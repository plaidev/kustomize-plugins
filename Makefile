build:
	cd ./plugin/sealdsecret; go build -buildmode plugin -o ../../bin/sealdsecret ./sealdsecret.go