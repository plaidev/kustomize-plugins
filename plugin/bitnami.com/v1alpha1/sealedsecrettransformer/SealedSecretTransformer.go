package main

import (
	"fmt"

	"encoding/json"

	"sigs.k8s.io/kustomize/api/hasher"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/types"
	apiv1 "sigs.k8s.io/kustomize/pseudo/k8s/api/core/v1"
	metav1 "sigs.k8s.io/kustomize/pseudo/k8s/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/kustomize/pseudo/k8s/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

// Copied from https://github.com/bitnami-labs/sealed-secrets/blob/v0.9.5/pkg/apis/sealed-secrets/v1alpha1/types.go
// And some are modified to sync k8s/api version

// SecretTemplateSpec describes the structure a Secret should have
// when created from a template
type SecretTemplateSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Type              apiv1.SecretType `json:"type,omitempty" protobuf:"bytes,3,opt,name=type,casttype=SecretType"`
}

// SealedSecretSpec is the specification of a SealedSecret
type SealedSecretSpec struct {
	Template      SecretTemplateSpec `json:"template,omitempty"`
	Data          []byte             `json:"data,omitempty"`
	EncryptedData map[string]string  `json:"encryptedData"`
}

// SealedSecretConditionType describes the type of SealedSecret condition
type SealedSecretConditionType string

// SealedSecretCondition describes the state of a sealed secret at a certain point.
type SealedSecretCondition struct {
	Type               SealedSecretConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=DeploymentConditionType"`
	Status             apiv1.ConditionStatus     `json:"status" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/api/core/v1.ConditionStatus"`
	LastUpdateTime     metav1.Time               `json:"lastUpdateTime,omitempty" protobuf:"bytes,6,opt,name=lastUpdateTime"`
	LastTransitionTime metav1.Time               `json:"lastTransitionTime,omitempty" protobuf:"bytes,7,opt,name=lastTransitionTime"`
	Reason             string                    `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`
	Message            string                    `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
}

// SealedSecretStatus is the most recently observed status of the SealedSecret.
type SealedSecretStatus struct {
	ObservedGeneration int64                   `json:"observedGeneration,omitempty" protobuf:"varint,3,opt,name=observedGeneration"`
	Conditions         []SealedSecretCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,6,rep,name=conditions"`
}

// SealedSecret is the K8s representation of a "sealed Secret" - a
// regular k8s Secret that has been sealed (encrypted) using the
// controller's key.
type SealedSecret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SealedSecretSpec   `json:"spec"`
	Status SealedSecretStatus `json:"status"`
}

// SealedSecretList represents a list of SealedSecrets
type SealedSecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []SealedSecret `json:"items"`
}

type plugin struct {
	h                *resmap.PluginHelpers
	types.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
}

// KustomizePlugin is plugin
var KustomizePlugin plugin

func (p *plugin) Config(h *resmap.PluginHelpers, config []byte) (err error) {
	err = yaml.Unmarshal(config, p)
	if p.Name == "" {
		p.Name = p.Name
	}
	if p.Namespace == "" {
		p.Namespace = p.Namespace
	}
	p.h = h
	return
}

// Transform appends hash to sealed secrets
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

// encodeSealedSecret encodes a SealedSecret.
// EncryptedData, Kind, Name, and Type are taken into account.
func encodeSealedSecret(sec *SealedSecret) (string, error) {
	// json.Marshal sorts the keys in a stable order in the encoding
	data, err := json.Marshal(map[string]interface{}{"kind": "SealedSecret", "data": sec.Spec.EncryptedData, "name": sec.Name})
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
