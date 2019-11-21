// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

//go:generate pluginator
package main

import (
	"fmt"

	"encoding/json"

	ss "github.com/bitnami-labs/sealed-secrets/pkg/apis/sealed-secrets/v1alpha1"
	"sigs.k8s.io/kustomize/api/hasher"
	"sigs.k8s.io/kustomize/api/ifc"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/pseudo/k8s/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

// SealedSecretArgs is args
type SealedSecretArgs struct {
	// GeneratorArgs for the secret.
	types.GeneratorArgs `json:",inline,omitempty" yaml:",inline,omitempty"`

	// Type of the secret.
	//
	// This is the same field as the secret type field in v1/Secret:
	// It can be "Opaque" (default), or "kubernetes.io/tls".
	//
	// If type is "kubernetes.io/tls", then "literals" or "files" must have exactly two
	// keys: "tls.key" and "tls.crt"
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
}

type plugin struct {
	h                *resmap.PluginHelpers
	hasher           ifc.KunstructuredHasher
	types.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	types.GeneratorOptions
	SealedSecretArgs
}

//noinspection GoUnusedGlobalVariable
var KustomizePlugin plugin

func (p *plugin) Config(h *resmap.PluginHelpers, config []byte) (err error) {
	p.GeneratorOptions = types.GeneratorOptions{}
	p.SealedSecretArgs = SealedSecretArgs{}
	err = yaml.Unmarshal(config, p)
	if p.SealedSecretArgs.Name == "" {
		p.SealedSecretArgs.Name = p.Name
	}
	if p.SealedSecretArgs.Namespace == "" {
		p.SealedSecretArgs.Namespace = p.Namespace
	}
	p.hasher = h.ResmapFactory().RF().Hasher()
	p.h = h
	return
}

// Transform appends hash to generated resources.
func (p *plugin) Transform(m resmap.ResMap) error {
	for _, res := range m.Resources() {
		u := unstructured.Unstructured{
			Object: res.Map(),
		}
		kind := u.GetKind()
		if kind == "SealedSecret" {
			sec, err := unstructuredToSealedSecret(u)
			if err != nil {
				return err
			}
			h, err := secretHash(sec)
			if err != nil {
				return err
			}
			res.SetName(fmt.Sprintf("%s-%s", res.GetName(), h))
		}
	}
	return nil
}

func secretHash(sec *ss.SealedSecret) (string, error) {
	encoded, err := encodeSealedSecret(sec)
	if err != nil {
		return "", err
	}
	h, err := hasher.Encode(hasher.Hash(encoded))
	if err != nil {
		return "", err
	}
	return h, nil
}

// encodeSecret encodes a Secret.
// Data, Kind, Name, and Type are taken into account.
func encodeSealedSecret(sec *ss.SealedSecret) (string, error) {
	// json.Marshal sorts the keys in a stable order in the encoding
	data, err := json.Marshal(map[string]interface{}{"kind": "SealedSecret", "spec": sec.Spec, "name": sec.Name})
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func unstructuredToSealedSecret(u unstructured.Unstructured) (*ss.SealedSecret, error) {
	marshaled, err := json.Marshal(u.Object)
	if err != nil {
		return nil, err
	}
	var out ss.SealedSecret
	err = json.Unmarshal(marshaled, &out)
	return &out, err
}
