# kustomize-plugins

## docker

- https://hub.docker.com/r/plaid/kustomize-plugins

## sealed secret transformer

This plugin was developed in aims to restart pods when sealed secret are modified. But this plugin won't work by default. Because sealed secret which metadata.name modified with `name-${hash}` cannot be decrypted bacause of the scope of sealed secret encryption. To make it work, you need to set sealed secret scope to `namespace-wide`. Before using this plugin, please consider the other options to restart pods when secret changed (e.g. https://github.com/stakater/Reloader).

more details about this.

- https://kubernetes.slack.com/archives/CM0H415UG/p1574409839114400
- https://github.com/bitnami-labs/sealed-secrets/issues/167

### prerequisites

- go 1.13.4
- kustomize 3.5.5
- sealed secret 0.12.3

### installation

```sh
$ git clone https://github.com/plaidev/kustomize-plugins.git
$ make setup
$ XDG_CONFIG_HOME=<PLUGIN_PATH> make build
# then sealed secret transformer plugin will be built in
# $XDG_CONFIG_HOME/kustomize/plugin/bitnami.com/v1alpha1/sealedsecrettransformer/SealedSecretTransformer.so
```

### test

- unit test: `make unit-test`
- test: `make test`
