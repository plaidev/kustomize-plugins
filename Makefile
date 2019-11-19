build:
	cd ./plugin/sealdsecret; go build -buildmode plugin -o ../../bin/sealdsecret ./sealdsecret.go
	# go build -buildmode plugin -o ./bin ./plugin/sealdsecret/sealdsecret.go