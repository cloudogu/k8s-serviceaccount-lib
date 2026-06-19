# Set these to the desired values
PROJECT_NAME=k8s-serviceaccount-lib
ARTIFACT_ID=k8s-serviceaccount-crd
APPEND_CRD_SUFFIX=false
VERSION=2.0.0

MAKEFILES_VERSION=10.8.0

GOTAG?=1.26.0
LINT_VERSION=v2.10.1
CLIENT_GEN_VERSION?=v0.35.3
CONTROLLER_GEN_VERSION?=v0.20.1

GO_BUILD_FLAGS?=-mod=vendor -a ./...
.DEFAULT_GOAL:=default

PRE_COMPILE = generate-deepcopy
IMAGE_IMPORT_TARGET=image-import
CHECK_VAR_TARGETS=check-all-vars-without-image

include build/make/variables.mk
INTEGRATION_TEST_NAME_PATTERN=.*_inttest$$

include build/make/self-update.mk
include build/make/dependencies-gomod.mk
include build/make/build.mk
include build/make/test-common.mk
include build/make/test-integration.mk
include build/make/test-unit.mk
include build/make/static-analysis.mk
include build/make/clean.mk
include build/make/release.mk
include build/make/mocks.mk

include build/make/digital-signature.mk
include build/make/k8s-controller.mk

CLIENT_GEN=$(UTILITY_BIN_PATH)/client-gen

default: compile

##@ Codegen
.PHONY: client-gen
client-gen: ${CLIENT_GEN} ## Download client-gen locally if necessary.

${CLIENT_GEN}: $(UTILITY_BIN_PATH)
	$(call go-get-tool,$(CLIENT_GEN),k8s.io/code-generator/cmd/client-gen@$(CLIENT_GEN_VERSION))

.PHONY: generate-client
generate-client: ${CLIENT_GEN} ## Generate client code from API definitions.
	@echo "Generating client..."
	@$(CLIENT_GEN) -v 5 \
		--output-dir "./" \
		--output-pkg "github.com/cloudogu/${PROJECT_NAME}" \
		--clientset-name "client" \
		--input "api/v1" \
		--input-base "${CURDIR}"

.PHONY: generate-crd-api
generate-crd-api: generate-deepcopy generate-client manifests
	@echo "generated deepcopy, api-client and manifests"

# Override make target to use k8s-serviceaccount-lib as label
.PHONY: crd-add-labels
crd-add-labels: $(BINARY_YQ)
	@echo "Adding labels to CRD..."
	@for file in ${HELM_CRD_SOURCE_DIR}/templates/*.yaml ; do \
		$(BINARY_YQ) -i e ".metadata.labels.app = \"ces\"" $${file} ;\
		$(BINARY_YQ) -i e ".metadata.labels.\"app.kubernetes.io/name\" = \"${PROJECT_NAME}\"" $${file} ;\
	done