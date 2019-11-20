module github.com/plaidev/kustomize-plugin/plugin/builtin/sealedsecrettransformer

go 1.13

require (
	github.com/bitnami-labs/sealed-secrets v0.9.5
	sigs.k8s.io/kustomize/api v0.2.0
	sigs.k8s.io/kustomize/pseudo/k8s v0.1.0
	sigs.k8s.io/yaml v1.1.0
)
