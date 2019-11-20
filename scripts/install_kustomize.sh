#!/usr/bin/env bash
set -eux

CURRENT_DIR=$(cd $(dirname $0); pwd)
ROOT_DIR=$CURRENT_DIR/..
BIN_DIR=$ROOT_DIR/bin

mkdir -p ${BIN_DIR}

# install kustomize
curl -L https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv3.4.0/kustomize_v3.4.0_darwin_amd64.tar.gz -o ${BIN_DIR}/kustomize.tar.gz
tar zxvf ${BIN_DIR}/kustomize.tar.gz -C ${BIN_DIR}
rm ${BIN_DIR}/kustomize.tar.gz
chmod +x ${BIN_DIR}/kustomize
