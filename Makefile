MAKEFILE_DIR := $(shell pwd)

test:
	cd ./plugin/sealdsecretgenerator; XDG_CONFIG_HOME=$(MAKEFILE_DIR); go test;

build:
	make test
	cd ./plugin/sealdsecretgenerator; XDG_CONFIG_HOME=$(MAKEFILE_DIR); go build -buildmode plugin -o ../../bin/sealdsecretgenerator ./SealdSecretGenerator.go