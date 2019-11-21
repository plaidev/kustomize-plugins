// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

//go:generate pluginator
package main

import (
	"fmt"

	"encoding/json"

	"sigs.k8s.io/kustomize/api/hasher"
	"sigs.k8s.io/kustomize/api/ifc"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/types"
	apiv1 "sigs.k8s.io/kustomize/pseudo/k8s/api/core/v1"
	metav1 "sigs.k8s.io/kustomize/pseudo/k8s/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/kustomize/pseudo/k8s/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

// SecretTemplateSpec describes the structure a Secret should have
// when created from a template
type SecretTemplateSpec struct {
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Used to facilitate programmatic handling of secret data.
	// +optional
	Type apiv1.SecretType `json:"type,omitempty" protobuf:"bytes,3,opt,name=type,casttype=SecretType"`
}

// SealedSecretSpec is the specification of a SealedSecret
type SealedSecretSpec struct {
	// Template defines the structure of the Secret that will be
	// created from this sealed secret.
	// +optional
	Template SecretTemplateSpec `json:"template,omitempty"`

	// Data is deprecated and will be removed eventually. Use per-value EncryptedData instead.
	Data          []byte            `json:"data,omitempty"`
	EncryptedData map[string]string `json:"encryptedData"`
}

// SealedSecretConditionType describes the type of SealedSecret condition
type SealedSecretConditionType string

// SealedSecretCondition describes the state of a sealed secret at a certain point.
type SealedSecretCondition struct {
	// Type of condition for a sealed secret.
	// Valid value: "Synced"
	Type SealedSecretConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=DeploymentConditionType"`
	// Status of the condition for a sealed secret.
	// Valid values for "Synced": "True", "False", or "Unknown".
	Status apiv1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/api/core/v1.ConditionStatus"`
	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty" protobuf:"bytes,6,opt,name=lastUpdateTime"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,7,opt,name=lastTransitionTime"`
	// The reason for the condition's last transition.
	Reason string `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`
	// A human readable message indicating details about the transition.
	Message string `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
}

// SealedSecretStatus is the most recently observed status of the SealedSecret.
type SealedSecretStatus struct {
	// ObservedGeneration reflects the generation most recently observed by the sealed-secrets controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,3,opt,name=observedGeneration"`

	// Represents the latest available observations of a sealed secret's current state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []SealedSecretCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,6,rep,name=conditions"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient

// SealedSecret is the K8s representation of a "sealed Secret" - a
// regular k8s Secret that has been sealed (encrypted) using the
// controller's key.
type SealedSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SealedSecretSpec   `json:"spec"`
	Status SealedSecretStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SealedSecretList represents a list of SealedSecrets
type SealedSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []SealedSecret `json:"items"`
}

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

func secretHash(sec *SealedSecret) (string, error) {
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
func encodeSealedSecret(sec *SealedSecret) (string, error) {
	// json.Marshal sorts the keys in a stable order in the encoding
	data, err := json.Marshal(map[string]interface{}{"kind": "SealedSecret", "spec": sec.Spec, "name": sec.Name})
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func unstructuredToSealedSecret(u unstructured.Unstructured) (*SealedSecret, error) {
	marshaled, err := json.Marshal(u.Object)
	if err != nil {
		return nil, err
	}
	var out SealedSecret
	err = json.Unmarshal(marshaled, &out)
	return &out, err
}
