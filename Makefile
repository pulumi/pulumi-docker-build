PROJECT_NAME := Pulumi Docker Build Resource Provider

PACK             := dockerbuild
PACKDIR          := sdk
PROJECT          := github.com/pulumi/pulumi-dockerbuild
NODE_MODULE_NAME := @pulumi/dockerbuild
NUGET_PKG_NAME   := Pulumi.DockerBuild

PROVIDER        := pulumi-resource-${PACK}
VERSION         ?= $(shell pulumictl get version)
PROVIDER_PATH   := provider
VERSION_PATH    := ${PROVIDER_PATH}.Version
SCHEMA_PATH     := ${PROVIDER_PATH}/cmd/pulumi-resource-${PACK}/schema.json

GOPATH			:= $(shell go env GOPATH)

WORKING_DIR     := $(shell pwd)
EXAMPLES_DIR    := ${WORKING_DIR}/examples/yaml
TESTPARALLELISM := 4

PYPI_VERSION    := $(shell pulumictl get version --language python)
NODE_VERSION    := $(shell pulumictl get version --language javascript)
DOTNET_VERSION  := $(shell pulumictl get version --language dotnet)

ensure:: tidy lint test_provider sdk

.PHONY:
tidy: provider/go.sum tests/go.sum

provider::
	(cd provider && go build -o $(WORKING_DIR)/bin/${PROVIDER} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))
	touch bin/pulumi-gen-${PACK}

provider_debug::
	(cd provider && go build -o $(WORKING_DIR)/bin/${PROVIDER} -gcflags="all=-N -l" -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

test_provider::
	cd tests && go test -short -v -count=1 -cover -timeout 2h -parallel ${TESTPARALLELISM} ./...

dotnet_sdk:: DOTNET_VERSION := $(shell pulumictl get version --language dotnet)
dotnet_sdk::
	rm -rf sdk/dotnet
	pulumi package gen-sdk $(WORKING_DIR)/bin/$(PROVIDER) --language dotnet
	cd ${PACKDIR}/dotnet/&& \
		echo "${DOTNET_VERSION}" >version.txt && \
		dotnet build /p:Version=${DOTNET_VERSION}

go_sdk:: $(WORKING_DIR)/bin/$(PROVIDER)
	rm -rf sdk/go
	pulumi package gen-sdk $(WORKING_DIR)/bin/$(PROVIDER) --language go

nodejs_sdk:: VERSION := $(shell pulumictl get version --language javascript)
nodejs_sdk::
	rm -rf sdk/nodejs
	pulumi package gen-sdk $(WORKING_DIR)/bin/$(PROVIDER) --language nodejs
	cd ${PACKDIR}/nodejs/ && \
		yarn install && \
		yarn run tsc && \
		cp ../../README.md ../../LICENSE package.json yarn.lock bin/ && \
		sed -i.bak 's/$${VERSION}/$(VERSION)/g' bin/package.json && \
		rm ./bin/package.json.bak

python_sdk:: PYPI_VERSION := $(shell pulumictl get version --language python)
python_sdk::
	rm -rf sdk/python
	pulumi package gen-sdk $(WORKING_DIR)/bin/$(PROVIDER) --language python
	cp README.md ${PACKDIR}/python/
	cd ${PACKDIR}/python/ && \
		python3 setup.py clean --all 2>/dev/null && \
		rm -rf ./bin/ ../python.bin/ && cp -R . ../python.bin && mv ../python.bin ./bin && \
		sed -i.bak -e 's/^VERSION = .*/VERSION = "$(PYPI_VERSION)"/g' -e 's/^PLUGIN_VERSION = .*/PLUGIN_VERSION = "$(VERSION)"/g' ./bin/setup.py && \
		rm ./bin/setup.py.bak && \
		cd ./bin && python3 setup.py build sdist

.PHONY:
gen_examples: examples/go examples/nodejs examples/python examples/dotnet examples/java

examples: gen_examples

examples/%: bin/${PROVIDER} ${WORKING_DIR}/examples/yaml/Pulumi.yaml
	rm -rf ${WORKING_DIR}/examples/$*
	pulumi convert \
		--cwd ${WORKING_DIR}/examples/yaml \
		--logtostderr \
		--generate-only \
		--non-interactive \
		--language $* \
		--out ${WORKING_DIR}/examples/$*

define pulumi_login
    export PULUMI_CONFIG_PASSPHRASE=asdfqwerty1234; \
    pulumi login --local;
endef

up::
	$(call pulumi_login) \
	cd ${EXAMPLES_DIR} && \
	pulumi stack init dev && \
	pulumi stack select dev && \
	pulumi config set name dev && \
	pulumi up -y

down::
	$(call pulumi_login) \
	cd ${EXAMPLES_DIR} && \
	pulumi stack select dev && \
	pulumi destroy -y && \
	pulumi stack rm dev -y

devcontainer::
	git submodule update --init --recursive .devcontainer
	git submodule update --remote --merge .devcontainer
	cp -f .devcontainer/devcontainer.json .devcontainer.json

.PHONY: build

build:: provider dotnet_sdk go_sdk nodejs_sdk python_sdk

# Required for the codegen action that runs in pulumi/pulumi
only_build:: build

lint::
	for DIR in "provider" "tests" ; do \
		pushd $$DIR && golangci-lint run --fix -c ../.golangci.yml --timeout 10m && popd ; \
	done

install:: install_nodejs_sdk install_dotnet_sdk
	cp $(WORKING_DIR)/bin/${PROVIDER} ${GOPATH}/bin

GO_TEST 	 := go test -v -count=1 -cover -timeout 2h -parallel ${TESTPARALLELISM}

test_all:: test_provider
	cd tests/sdk/nodejs && $(GO_TEST) ./...
	cd tests/sdk/python && $(GO_TEST) ./...
	cd tests/sdk/dotnet && $(GO_TEST) ./...
	cd tests/sdk/go && $(GO_TEST) ./...

install_dotnet_sdk::
	rm -rf $(WORKING_DIR)/nuget/$(NUGET_PKG_NAME).*.nupkg
	mkdir -p $(WORKING_DIR)/nuget
	find . -name '*.nupkg' -print -exec cp -p {} ${WORKING_DIR}/nuget \;

install_python_sdk::
	#target intentionally blank

install_go_sdk::
	#target intentionally blank

install_nodejs_sdk::
	-yarn unlink --cwd $(WORKING_DIR)/sdk/nodejs/bin
	yarn link --cwd $(WORKING_DIR)/sdk/nodejs/bin

.PHONY: codegen
codegen:

.PHONY: generate_schema
generate_schema:

.PHONY: build_go install_go_sdk
generate_go: sdk/go
build_go:

.PHONY: build_java install_java_sdk
generate_java: sdk/java
build_java:

.PHONY: build_python install_python_sdk
generate_python: sdk/python
build_python:

.PHONY: build_nodejs install_nodejs_sdk
generate_nodejs: sdk/nodejs
build_nodejs:

.PHONY: build_dotnet install_dotnet_sdk
generate_dotnet: sdk/dotnet
build_dotnet:

${SCHEMA_PATH}: bin/${PROVIDER}
	pulumi package get-schema bin/${PROVIDER} > $(SCHEMA_PATH)

bin/${PROVIDER}: provider/**.go
	(cd provider && go build -o ../bin/${PROVIDER} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

provider/go.sum: provider/go.mod
	cd provider && go mod tidy

tests/go.sum: tests/go.mod
	cd tests && go mod tidy

$(shell mkdir -p sdk)
sdk: sdk/python sdk/nodejs sdk/java sdk/python sdk/go sdk/dotnet

TMPDIR := $(shell mktemp -d)
sdk/python: bin/${PROVIDER}
	pulumi package gen-sdk bin/$(PROVIDER) --language python -o ${TMPDIR}
	cp README.md ${TMPDIR}/python/
	cd ${TMPDIR}/python/ && \
		rm -rf ./bin/ ../python.bin/ && cp -R . ../python.bin && mv ../python.bin ./bin && \
		sed -i.bak -e 's/^  version = .*/  version = "$(PYPI_VERSION)"/g' ./bin/pyproject.toml && \
		rm ./bin/pyproject.toml.bak && \
		python3 -m venv venv && \
		./venv/bin/python -m pip install build && \
		cd ./bin && \
		../venv/bin/python -m build .
	mv ${TMPDIR}/python ${WORKING_DIR}/sdk/python

TMPDIR := $(shell mktemp -d)
sdk/nodejs: bin/${PROVIDER}
	pulumi package gen-sdk bin/$(PROVIDER) --language nodejs -o ${TMPDIR}
	cp README.md LICENSE ${TMPDIR}/nodejs
	cd ${TMPDIR}/nodejs/ && \
		yarn install && \
		yarn run tsc && \
		cp README.md LICENSE package.json yarn.lock bin/ && \
		sed -i.bak 's/$${VERSION}/$(NODE_VERSION)/g' bin/package.json && \
		rm ./bin/package.json.bak
	mv ${TMPDIR}/nodejs ${WORKING_DIR}/sdk/nodejs

TMPDIR := $(shell mktemp -d)
sdk/go: bin/${PROVIDER}
	pulumi package gen-sdk bin/$(PROVIDER) --language go -o ${TMPDIR}
	mv ${TMPDIR}/go ${WORKING_DIR}/sdk/go

TMPDIR := $(shell mktemp -d)
sdk/dotnet: bin/${PROVIDER}
	pulumi package gen-sdk bin/${PROVIDER} --language dotnet -o ${TMPDIR}
	cd ${TMPDIR}/dotnet/ && \
		echo "$(DOTNET_VERSION)" > version.txt && \
		dotnet build /p:Version=${DOTNET_VERSION}
	mv ${TMPDIR}/dotnet ${WORKING_DIR}/sdk/dotnet

TMPDIR := $(shell mktemp -d)
sdk/java: bin/${PROVIDER}
	pulumi package gen-sdk --language java bin/${PROVIDER} -o ${TMPDIR}
	cd ${TMPDIR}/java/ && gradle --console=plain build
	mv ${TMPDIR}/java ${WORKING_DIR}/sdk/java
