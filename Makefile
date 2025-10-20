PROJECT_NAME := Pulumi Docker Build Resource Provider

PACK             := docker-build
PACKDIR          := sdk
PROJECT          := github.com/pulumi/pulumi-docker-build
NODE_MODULE_NAME := @pulumi/docker-build
NUGET_PKG_NAME   := Pulumi.DockerBuild

PROVIDER         := pulumi-resource-${PACK}
PROVIDER_PATH    := provider
VERSION_PATH     := ${PROVIDER_PATH}.Version
SCHEMA_PATH      := ${PROVIDER_PATH}/cmd/pulumi-resource-${PACK}/schema.json

GOPATH			 := $(shell go env GOPATH)

WORKING_DIR      := $(shell pwd)
EXAMPLES_DIR     := ${WORKING_DIR}/examples/yaml
TESTPARALLELISM  := 4

PULUMI           := pulumi
GOGLANGCILINT    := golangci-lint
GOTEST           := go test

# Override during CI using `make [TARGET] PROVIDER_VERSION=""` or by setting a PROVIDER_VERSION environment variable
# Local & branch builds will just used this fixed default version unless specified
PROVIDER_VERSION ?= 0.1.0-alpha.0+dev
# Use this normalised version everywhere rather than the raw input to ensure consistency.
VERSION_GENERIC = $(shell pulumictl convert-version --language generic --version "$(PROVIDER_VERSION)")

export PULUMI_IGNORE_AMBIENT_PLUGINS = true
export PULUMI_DISABLE_AUTOMATIC_PLUGIN_ACQUISITION = true

.PHONY: ensure
ensure:: tidy lint test_provider examples

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: provider
provider: bin/${PROVIDER} bin/pulumi-gen-${PACK} # Required by CI

.PHONY: local_generate
local_generate: sdk # Required by CI

provider_debug::
	(cd provider && go build -o $(WORKING_DIR)/bin/${PROVIDER} -gcflags="all=-N -l" -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION_GENERIC}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

test_provider:: # Required by CI
	${GOTEST} -short -v -coverprofile="coverage.txt" -coverpkg=./provider/... -timeout 2h -parallel ${TESTPARALLELISM} ./provider/...

test_examples: install_nodejs_sdk install_dotnet_sdk
	${GOTEST} -short -v -cover -tags=all -timeout 2h -parallel ${TESTPARALLELISM} ./examples/...

test_all:: test_provider test_examples

.PHONY:
gen_examples:

examples: $(shell mkdir -p examples)
examples: sdk examples/yaml examples/go examples/nodejs examples/python examples/dotnet examples/java

examples/yaml:
	rm -rf ${WORKING_DIR}/examples/yaml/app
	cp -r ${WORKING_DIR}/examples/app ${WORKING_DIR}/examples/yaml/app

examples/go: bin/${PROVIDER} ${WORKING_DIR}/examples/yaml/Pulumi.yaml
	$(call example,go)
	@git checkout examples/go/go.mod

examples/nodejs: bin/${PROVIDER} ${WORKING_DIR}/examples/yaml/Pulumi.yaml
	$(call example,nodejs)
	@git checkout examples/nodejs/package.json

examples/python: bin/${PROVIDER} ${WORKING_DIR}/examples/yaml/Pulumi.yaml
	$(call example,python)
	@git checkout examples/python/requirements.txt

examples/dotnet: bin/${PROVIDER} ${WORKING_DIR}/examples/yaml/Pulumi.yaml
	$(call example,dotnet)
	@git checkout examples/dotnet/provider-docker-build.csproj

examples/java: bin/${PROVIDER} ${WORKING_DIR}/examples/yaml/Pulumi.yaml
	$(call example,java)
	@git checkout examples/java/pom.xml

define pulumi_login
    export PULUMI_CONFIG_PASSPHRASE=asdfqwerty1234; \
    pulumi login --local;
endef

define example
	rm -rf ${WORKING_DIR}/examples/$(1)
	pulumi convert \
		--cwd ${WORKING_DIR}/examples/yaml \
		--logtostderr \
		--generate-only \
		--non-interactive \
		--language $(1) \
		--out ${WORKING_DIR}/examples/$(1)
	cp -r ${WORKING_DIR}/examples/app ${WORKING_DIR}/examples/$(1)/app
	cp ${WORKING_DIR}/examples/yaml/.dockerignore ${WORKING_DIR}/examples/$(1)/.dockerignore
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
build:: provider sdk/dotnet sdk/go sdk/nodejs sdk/python sdk/java ${SCHEMA_PATH}

# Required for the codegen action that runs in pulumi/pulumi
only_build:: build

.PHONY: lint
lint:
	${GOGLANGCILINT} run --fix -c .golangci.yml

install:: install_nodejs_sdk install_dotnet_sdk
	cp $(WORKING_DIR)/bin/${PROVIDER} ${GOPATH}/bin


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
generate_schema: ${SCHEMA_PATH} # Required by CI

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
	pulumi package get-schema ./bin/${PROVIDER} | jq 'del(.version)' > $(SCHEMA_PATH)

bin/${PROVIDER}: $(shell find ./provider -name '*.go') go.mod
	(cd provider && go build -o ../bin/${PROVIDER} -ldflags "-X ${PROJECT}/${VERSION_PATH}=${VERSION_GENERIC}" $(PROJECT)/${PROVIDER_PATH}/cmd/$(PROVIDER))

bin/pulumi-gen-${PACK}: # Required by CI
	touch bin/pulumi-gen-${PACK}

go.mod: $(shell find . -name '*.go')
go.sum: go.mod

sdk: $(shell mkdir -p sdk)
sdk: sdk/python sdk/nodejs sdk/java sdk/python sdk/go sdk/dotnet

# Folders can't be used for up-to-date checks as they will be marked as up-to-date even if the step fails - leading to a broken state.
.PHONY: sdk/*

sdk/python: TMPDIR := $(shell mktemp -d)
sdk/python: bin/${PROVIDER}
	rm -rf sdk/python
	$(PULUMI) package gen-sdk ./bin/$(PROVIDER) --language python -o ${TMPDIR}
	cp README.md ${TMPDIR}/python/
	cd ${TMPDIR}/python/ && \
		rm -rf ./bin/ ../python.bin/ && cp -R . ../python.bin && mv ../python.bin ./bin && \
		python3 -m venv venv && \
		./venv/bin/python -m pip install build && \
		cd ./bin && \
		../venv/bin/python -m build .
	mv -f ${TMPDIR}/python ${WORKING_DIR}/sdk/.

sdk/nodejs: TMPDIR := $(shell mktemp -d)
sdk/nodejs: bin/${PROVIDER}
	rm -rf sdk/nodejs
	$(PULUMI) package gen-sdk ./bin/$(PROVIDER) --language nodejs -o ${TMPDIR}
	cp README.md LICENSE ${TMPDIR}/nodejs
	cd ${TMPDIR}/nodejs/ && \
		yarn install && \
		yarn run tsc && \
		cp README.md LICENSE package.json yarn.lock bin/
	mv -f ${TMPDIR}/nodejs ${WORKING_DIR}/sdk/.

sdk/go: TMPDIR := $(shell mktemp -d)
sdk/go: PATH := "$(WORKING_DIR)/bin:$(PATH)"
sdk/go: bin/${PROVIDER}
	rm -rf sdk/go
	PATH=$(PATH) $(PULUMI) package gen-sdk ./bin/$(PROVIDER) --language go -o ${TMPDIR}
	cp go.mod ${TMPDIR}/go/dockerbuild/go.mod
	cd ${TMPDIR}/go/dockerbuild && \
		go mod edit -module=github.com/pulumi/pulumi-${PACK}/${PACKDIR}/go/dockerbuild && \
		go mod tidy
	mv -f ${TMPDIR}/go ${WORKING_DIR}/sdk/go

sdk/dotnet: TMPDIR := $(shell mktemp -d)
sdk/dotnet: bin/${PROVIDER}
	rm -rf sdk/dotnet
	$(PULUMI) package gen-sdk ./bin/${PROVIDER} --language dotnet -o ${TMPDIR}
	cd ${TMPDIR}/dotnet/ && \
		echo "$(VERSION_GENERIC)" > version.txt && \
		dotnet build
	mv -f ${TMPDIR}/dotnet ${WORKING_DIR}/sdk/.

sdk/java: PACKAGE_VERSION := $(shell pulumictl convert-version --language generic -v "$(VERSION_GENERIC)")
sdk/java: TMPDIR := $(shell mktemp -d)
sdk/java: bin/${PROVIDER}
	rm -rf sdk/java
	$(PULUMI) package gen-sdk --language java ./bin/${PROVIDER} -o ${TMPDIR}
	cd ${TMPDIR}/java/ && gradle --console=plain build
	mv -f ${TMPDIR}/java ${WORKING_DIR}/sdk/.

docs: $(shell find docs/yaml -type f) $(shell find ./provider/internal/embed -name '*.md') ${SCHEMA_PATH}
	go generate docs/generate.go
	@touch docs

# Set these variables to enable signing of the windows binary
AZURE_SIGNING_CLIENT_ID ?=
AZURE_SIGNING_CLIENT_SECRET ?=
AZURE_SIGNING_TENANT_ID ?=
AZURE_SIGNING_KEY_VAULT_URI ?=
SKIP_SIGNING ?=

bin/jsign-6.0.jar:
	wget https://github.com/ebourg/jsign/releases/download/6.0/jsign-6.0.jar --output-document=bin/jsign-6.0.jar

sign-goreleaser-exe-amd64: GORELEASER_ARCH := amd64_v1
sign-goreleaser-exe-arm64: GORELEASER_ARCH := arm64

# Set the shell to bash to allow for the use of bash syntax.
sign-goreleaser-exe-%: SHELL:=/bin/bash
sign-goreleaser-exe-%: bin/jsign-6.0.jar
	@# Only sign windows binary if fully configured.
	@# Test variables set by joining with | between and looking for || showing at least one variable is empty.
	@# Move the binary to a temporary location and sign it there to avoid the target being up-to-date if signing fails.
	@set -e; \
	if [[ "${SKIP_SIGNING}" != "true" ]]; then \
		if [[ "|${AZURE_SIGNING_CLIENT_ID}|${AZURE_SIGNING_CLIENT_SECRET}|${AZURE_SIGNING_TENANT_ID}|${AZURE_SIGNING_KEY_VAULT_URI}|" == *"||"* ]]; then \
			echo "Can't sign windows binaries as required configuration not set: AZURE_SIGNING_CLIENT_ID, AZURE_SIGNING_CLIENT_SECRET, AZURE_SIGNING_TENANT_ID, AZURE_SIGNING_KEY_VAULT_URI"; \
			echo "To rebuild with signing delete the unsigned windows exe file and rebuild with the fixed configuration"; \
			if [[ "${CI}" == "true" ]]; then exit 1; fi; \
		else \
			file=dist/build-provider-sign-windows_windows_${GORELEASER_ARCH}/pulumi-resource-docker-build.exe; \
			mv $${file} $${file}.unsigned; \
			az login --service-principal \
				--username "${AZURE_SIGNING_CLIENT_ID}" \
				--password "${AZURE_SIGNING_CLIENT_SECRET}" \
				--tenant "${AZURE_SIGNING_TENANT_ID}" \
				--output none; \
			ACCESS_TOKEN=$$(az account get-access-token --resource "https://vault.azure.net" | jq -r .accessToken); \
			java -jar bin/jsign-6.0.jar \
				--storetype AZUREKEYVAULT \
				--keystore "PulumiCodeSigning" \
				--url "${AZURE_SIGNING_KEY_VAULT_URI}" \
				--storepass "$${ACCESS_TOKEN}" \
				$${file}.unsigned; \
			mv $${file}.unsigned $${file}; \
			az logout; \
		fi; \
	fi

# To make an immediately observable change to .ci-mgmt.yaml:
#
# - Edit .ci-mgmt.yaml
# - Run make ci-mgmt to apply the change locally.
#
ci-mgmt: .ci-mgmt.yaml
	go run github.com/pulumi/ci-mgmt/provider-ci@master generate
.PHONY: ci-mgmt
	fi
