# kustomize-plugins

## docker

- https://hub.docker.com/repository/docker/plaid/kustomize-plugins

## sealed secret transformer
This plugins was developed in aims to restart pods when sealed secret are modified. But this plugin doesnt work by defualt. Because sealed secret which metadata.name modified with hash cannot be decrypted bacause of the implementation of sealed secrte encryption. To make it work, you need to set sealed secret scope to `namespace-wide`. Before using this plugin, please consider some options to restart pods when secred changed (e.g. https://github.com/stakater/Reloader). 

more details about this.
- https://kubernetes.slack.com/archives/CM0H415UG/p1574409839114400
- https://github.com/plaidev/kustomize-plugins

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
