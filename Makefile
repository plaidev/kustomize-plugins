MAKEFILE_DIR := $(shell pwd)
SELAED_SECRET_OUT=$(XDG_CONFIG_HOME)/kustomize/plugin/bitnami.com/v1alpha1/sealedsecrettransformer/SealedSecretTransformer.so
TEST_XDG_CONFIG_HOME=$(MAKEFILE_DIR)/test
PLUGIN_TEST_BIN=$(TEST_XDG_CONFIG_HOME)/kustomize/plugin/bitnami.com/v1alpha1/sealedsecrettransformer/SealedSecretTransformer.so
BIN_DIR=$(MAKEFILE_DIR)/bin

setup:
	./scripts/install_kustomize.sh

build:
	mkdir -p $(dir $(SELAED_SECRET_OUT))
	cd ./plugin/bitnami.com/v1alpha1/sealedsecrettransformer; go build -buildmode plugin -o $(SELAED_SECRET_OUT) ./SealedSecretTransformer.go

unit-test:
	cd ./plugin/bitnami.com/v1alpha1/sealedsecrettransformer; go test;

test:
	mkdir -p $(dir $(PLUGIN_TEST_BIN))
	cd ./plugin/bitnami.com/v1alpha1/sealedsecrettransformer; go build -buildmode plugin -o $(PLUGIN_TEST_BIN) ./SealedSecretTransformer.go
	XDG_CONFIG_HOME=$(TEST_XDG_CONFIG_HOME) $(BIN_DIR)/kustomize build $(TEST_XDG_CONFIG_HOME)/app --enable_alpha_plugins

.PHONY: setup build unit-test test