# kustomize-plugins

## docker

- https://hub.docker.com/repository/docker/plaid/kustomize-plugins

## sealed secret transformer

### prerequisites

- go 1.13.4
- kustomize 3.4.0
- sealed secret 0.9.5

### installation

```
$ git clone https://github.com/plaidev/kustomize-plugins.git
$ make setup
$ XDG_CONFIG_HOME=<PLUGIN_PATH> make build
# $XDG_CONFIG_HOME/kustomize/plugin/bitnami.com/v1alpha1/sealedsecrettransformer/SealedSecretTransformer.so will be made
```

### test

- unit test: `make unit-test`
- test: `make test`
