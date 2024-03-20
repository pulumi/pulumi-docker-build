PROJECT_NAME := Pulumi Docker Build Resource Provider

PACK             := dockerbuild
PACKDIR          := sdk
PROJECT          := github.com/pulumi/pulumi-dockerbuild
NODE_MODULE_NAME := @pulumi/dockerbuild
NUGET_PKG_NAME   := Pulumi.Dockerbuild

PROVIDER         := pulumi-resource-${PACK}
VERSION          ?= $(shell pulumictl get version)
PROVIDER_PATH    := provider
VERSION_PATH     := ${PROVIDER_PATH}.Version
SCHEMA_PATH      := ${PROVIDER_PATH}/cmd/pulumi-resource-${PACK}/schema.json

GOPATH			 := $(shell go env GOPATH)

WORKING_DIR      := $(shell pwd)
EXAMPLES_DIR     := ${WORKING_DIR}/examples/yaml
TESTPARALLELISM  := 4

PULUMI           := bin/pulumi

.PHONY: ensure
ensure:: tidy lint test_provider examples

.PHONY: tidy
tidy: go.sum

.PHONY: provider
provider: bin/${PROVIDER} bin/pulumi-gen-${PACK} # Required by CI

provider_debug::
	(cd provider && go build -o $(WORKING_DIR)/bin/${PROVIDER} -gcflags="all=-N -l" -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

test_provider:: # Required by CI
	go test -short -v -cover -timeout 2h -parallel ${TESTPARALLELISM} ./provider/...

test_examples: install_nodejs_sdk install_dotnet_sdk
	go test -short -v -cover -tags=all -timeout 2h -parallel ${TESTPARALLELISM} ./examples/...

test_all:: test_provider test_examples

.PHONY:
gen_examples:

examples: $(shell mkdir -p examples)
examples: sdk examples/go examples/nodejs examples/python examples/dotnet examples/java

examples/go: ${PULUMI} bin/${PROVIDER} ${WORKING_DIR}/examples/yaml/Pulumi.yaml
	$(call example,go)
	@git checkout examples/go/go.mod

examples/nodejs: ${PULUMI} bin/${PROVIDER} ${WORKING_DIR}/examples/yaml/Pulumi.yaml
	$(call example,nodejs)
	@git checkout examples/nodejs/package.json

examples/python: ${PULUMI} bin/${PROVIDER} ${WORKING_DIR}/examples/yaml/Pulumi.yaml
	$(call example,python)
	@git checkout examples/python/requirements.txt

examples/dotnet: ${PULUMI} bin/${PROVIDER} ${WORKING_DIR}/examples/yaml/Pulumi.yaml
	$(call example,dotnet)
	@git checkout examples/dotnet/provider-dockerbuild.csproj

examples/java: ${PULUMI} bin/${PROVIDER} ${WORKING_DIR}/examples/yaml/Pulumi.yaml
	$(call example,java)
	@git checkout examples/java/pom.xml

${PULUMI}: go.sum
	GOBIN=${WORKING_DIR}/bin go install github.com/pulumi/pulumi/pkg/v3/cmd/pulumi

define pulumi_login
    export PULUMI_CONFIG_PASSPHRASE=asdfqwerty1234; \
    pulumi login --local;
endef

define example
	echo "GOT $(1)"
	rm -rf ${WORKING_DIR}/examples/$(1)
	$(PULUMI) convert \
		--cwd ${WORKING_DIR}/examples/yaml \
		--logtostderr \
		--generate-only \
		--non-interactive \
		--language $(1) \
		--out ${WORKING_DIR}/examples/$(1)
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

.PHONY: lint
lint:
	golangci-lint run --fix -c .golangci.yml --timeout 10m

install:: install_nodejs_sdk install_dotnet_sdk
	cp $(WORKING_DIR)/bin/${PROVIDER} ${GOPATH}/bin

GO_TEST 	 := go test -v -count=1 -cover -timeout 2h -parallel ${TESTPARALLELISM}


install_dotnet_sdk:: # Required by CI
	rm -rf $(WORKING_DIR)/nuget/$(NUGET_PKG_NAME).*.nupkg
	mkdir -p $(WORKING_DIR)/nuget
	find . -name '*.nupkg' -print -exec cp -p {} ${WORKING_DIR}/nuget \;

install_python_sdk:: # Required by CI

install_go_sdk:: # Required by CI

install_nodejs_sdk:: # Required by CI
	-yarn unlink --cwd $(WORKING_DIR)/sdk/nodejs/bin
	yarn link --cwd $(WORKING_DIR)/sdk/nodejs/bin

.PHONY: codegen
codegen: # Required by CI

.PHONY: generate_schema
generate_schema: # Required by CI

.PHONY: build_go install_go_sdk
generate_go: sdk/go # Required by CI
build_go: # Required by CI

.PHONY: build_java install_java_sdk
generate_java: sdk/java # Required by CI
build_java: # Required by CI

.PHONY: build_python install_python_sdk
generate_python: sdk/python # Required by CI
build_python: # Required by CI

.PHONY: build_nodejs install_nodejs_sdk
generate_nodejs: sdk/nodejs # Required by CI
build_nodejs: # Required by CI

.PHONY: build_dotnet install_dotnet_sdk
generate_dotnet: sdk/dotnet # Required by CI
build_dotnet: # Required by CI

${SCHEMA_PATH}: bin/${PROVIDER}
	pulumi package get-schema bin/${PROVIDER} > $(SCHEMA_PATH)

bin/${PROVIDER}: $(shell find ./provider -name '*.go') go.mod
	(cd provider && go build -o ../bin/${PROVIDER} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

bin/pulumi-gen-${PACK}: # Required by CI
	touch bin/pulumi-gen-${PACK}

$(shell find . -name '*.go'):

go.mod: $(shell find . -name '*.go')
go.sum: go.mod
	go mod tidy

sdk: $(shell mkdir -p sdk)
sdk: sdk/python sdk/nodejs sdk/java sdk/python sdk/go sdk/dotnet

sdk/python: PYPI_VERSION := $(shell pulumictl get version --language python)
sdk/python: TMPDIR := $(shell mktemp -d)
sdk/python: $(PULUMI) bin/${PROVIDER}
	rm -rf sdk/python
	$(PULUMI) package gen-sdk bin/$(PROVIDER) --language python -o ${TMPDIR}
	cp README.md ${TMPDIR}/python/
	cd ${TMPDIR}/python/ && \
		rm -rf ./bin/ ../python.bin/ && cp -R . ../python.bin && mv ../python.bin ./bin && \
		sed -i.bak -e 's/^  version = .*/  version = "$(PYPI_VERSION)"/g' ./bin/pyproject.toml && \
		rm ./bin/pyproject.toml.bak && \
		python3 -m venv venv && \
		./venv/bin/python -m pip install build && \
		cd ./bin && \
		../venv/bin/python -m build .
	mv -f ${TMPDIR}/python ${WORKING_DIR}/sdk/.

sdk/nodejs: NODE_VERSION := $(shell pulumictl get version --language javascript)
sdk/nodejs: TMPDIR := $(shell mktemp -d)
sdk/nodejs: $(PULUMI) bin/${PROVIDER}
	rm -rf sdk/nodejs
	$(PULUMI) package gen-sdk bin/$(PROVIDER) --language nodejs -o ${TMPDIR}
	cp README.md LICENSE ${TMPDIR}/nodejs
	cd ${TMPDIR}/nodejs/ && \
		yarn install && \
		yarn run tsc && \
		cp README.md LICENSE package.json yarn.lock bin/ && \
		sed -i.bak 's/$${VERSION}/$(NODE_VERSION)/g' bin/package.json && \
		rm ./bin/package.json.bak
	mv -f ${TMPDIR}/nodejs ${WORKING_DIR}/sdk/.

sdk/go: TMPDIR := $(shell mktemp -d)
sdk/go: $(PULUMI) bin/${PROVIDER}
	rm -rf sdk/go
	$(PULUMI) package gen-sdk bin/$(PROVIDER) --language go -o ${TMPDIR}
	cp go.mod ${TMPDIR}/go/${PACK}/go.mod
	cd ${TMPDIR}/go/${PACK} && \
		go mod edit -module=github.com/pulumi/pulumi-${PACK}/${PACKDIR}/go/${PACK} && \
		go mod tidy
	mv -f ${TMPDIR}/go ${WORKING_DIR}/sdk/.

sdk/dotnet: DOTNET_VERSION  := $(shell pulumictl get version --language dotnet)
sdk/dotnet: TMPDIR := $(shell mktemp -d)
sdk/dotnet: $(PULUMI) bin/${PROVIDER}
	rm -rf sdk/dotnet
	$(PULUMI) package gen-sdk bin/${PROVIDER} --language dotnet -o ${TMPDIR}
	cd ${TMPDIR}/dotnet/ && \
		echo "$(DOTNET_VERSION)" > version.txt && \
		dotnet build /p:Version=${DOTNET_VERSION}
	mv -f ${TMPDIR}/dotnet ${WORKING_DIR}/sdk/.

sdk/java: PACKAGE_VERSION := $(shell pulumictl get version --language generic)
sdk/java: TMPDIR := $(shell mktemp -d)
sdk/java: $(PULUMI) bin/${PROVIDER}
	rm -rf sdk/java
	$(PULUMI) package gen-sdk --language java bin/${PROVIDER} -o ${TMPDIR}
	cd ${TMPDIR}/java/ && gradle --console=plain build
	mv -f ${TMPDIR}/java ${WORKING_DIR}/sdk/.
