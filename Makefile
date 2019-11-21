MAKEFILE_DIR := $(shell pwd)
XDG_CONFIG_HOME=$(MAKEFILE_DIR)/test
PLUGIN_BIN=$(XDG_CONFIG_HOME)/kustomize/plugin/bitnami.com/v1alpha1/sealedsecrettransformer/SealedSecretTransformer.so
BIN_DIR=$(MAKEFILE_DIR)/bin

setup:
	./scripts/install_kustomize.sh

build:
	mkdir -p $(dir $(PLUGIN_BIN))
	mkdir -p $(XDG_CONFIG_HOME)/app/kustomizeconfig
	cp ./transformerconfigs/bitnami.com/v1alpha1/sealedsecret.yml $(XDG_CONFIG_HOME)/app/kustomizeconfig
	cd ./plugin/bitnami.com/v1alpha1/sealedsecrettransformer; go build -buildmode plugin -o $(PLUGIN_BIN) ./SealedSecretTransformer.go

unit-test:
	cd ./plugin/bitnami.com/v1alpha1/sealedsecrettransformer; go test;

test:
	XDG_CONFIG_HOME=$(XDG_CONFIG_HOME) $(BIN_DIR)/kustomize build $(XDG_CONFIG_HOME)/app --enable_alpha_plugins

.PHONY: setup build unit-test test