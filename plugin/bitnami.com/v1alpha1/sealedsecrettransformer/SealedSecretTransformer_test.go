// Learn more or give us feedback
// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package main_test

import (
	"fmt"
	"testing"

	kusttest_test "sigs.k8s.io/kustomize/api/testutils/kusttest"
)

func TestTransformer(t *testing.T) {
	tc := kusttest_test.NewPluginTestEnv(t).Set()
	defer tc.Reset()

	tc.BuildGoPlugin(
		"bitnami.com/v1alpha1", "", "SealedSecretTransformer")

	th := kusttest_test.NewKustTestHarnessAllowPlugins(t, "/app")

	rm := th.LoadAndRunTransformer(`
apiVersion: bitnami.com/v1alpha1
kind: SealedSecretTransformer
metadata:
  name: hash
disableNameSuffixHash: false
`, `
apiVersion: v1
kind: Secret
metadata:
  name: mySecret
spec:
  encryptedData:
    dockerConfig: Aiueo
---
apiVersion: bitnami.com/v1alpha1
kind: SealedSecret
metadata:
  name: mySealedSecret
spec:
  encryptedData:
    dockerConfig: Aiueo
---
apiVersion: v1
group: apps
kind: Deployment
metadata:
  name: deploy1
spec:
  template:
    metadata:
      name: app-nnginx
    spec:
      containers:
      - name: ngnix
        image: nginx:1.7.9
      volumes:
      - name: secret-volume
        secret:
          secretName: mySealedSecret
      - name: secret-volume-2
        secret:
          secretName: mySecret
`)
	fmt.Println(rm)

	actual, err := rm.AsYaml()
	if err != nil {
		t.Fatal("ee")
	}

	fmt.Printf("%s", actual)
	// 	th.AssertActualEqualsExpected(rm, `
	// apiVersion: v1
	// data:
	//   DB_PASSWORD: aWxvdmV5b3U=
	//   FRUIT: YXBwbGU=
	//   ROUTER_PASSWORD: YWRtaW4=
	//   VEGETABLE: Y2Fycm90
	//   obscure: CkxvcmVtIGlwc3VtIGRvbG9yIHNpdCBhbWV0LApjb25zZWN0ZXR1ciBhZGlwaXNjaW5nIGVsaXQuCg==
	// kind: Secret
	// metadata:
	//   name: mySecret
	//   namespace: whatever
	// type: Opaque
	// `)
}
