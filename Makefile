MAKEFILE_DIR := $(shell pwd)
XDG_CONFIG_HOME=$(MAKEFILE_DIR)/temp
PLUGIN_BIN=$(XDG_CONFIG_HOME)/kustomize/plugin/bitnami.com/v1alpha1/sealedsecrettransformer/SealedSecretTransformer.so
BIN_DIR=$(MAKEFILE_DIR)/bin

setup:
	./scripts/install_kustomize.sh
unit-test:
	cd ./plugin/v1/sealedsecrettransformer; go test;

test:
	rm -rf $(XDG_CONFIG_HOME) || true
	mkdir -p $(dir $(PLUGIN_BIN))
	cd ./plugin/v1/sealedsecrettransformer; go build -buildmode plugin -o $(PLUGIN_BIN) ./SealedSecretTransformer.go
	mkdir -p ./test/kustomizeconfig
	cp ./transformerconfigs/v1/sealedsecret.yml ./test/kustomizeconfig
	XDG_CONFIG_HOME=$(XDG_CONFIG_HOME) $(BIN_DIR)/kustomize build ./test --enable_alpha_plugins

build:
	make test
	cd ./plugin/sealdsecretgenerator; XDG_CONFIG_HOME=$(MAKEFILE_DIR); go build -buildmode plugin -o ../../bin/sealdsecretgenerator ./SealdSecretGenerator.go

.PHONY: unit-test test build