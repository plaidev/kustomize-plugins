#!/usr/bin/env bash
set -eux

CURRENT_DIR=$(cd $(dirname $0); pwd)
ROOT_DIR=$CURRENT_DIR/..
BIN_DIR=$ROOT_DIR/bin
VENDOR_DIR=$ROOT_DIR/vendor

rm -rf ${BIN_DIR} ${VENDOR_DIR} || true

mkdir -p ${BIN_DIR}
mkdir -p ${VENDOR_DIR}


# get kustomize 3.4.0
curl -L https://github.com/kubernetes-sigs/kustomize/archive/kustomize/v3.4.0.tar.gz -o ${VENDOR_DIR}/kustomize.tar.gz
tar zxvf ${VENDOR_DIR}/kustomize.tar.gz -C ${VENDOR_DIR}
rm ${VENDOR_DIR}/kustomize.tar.gz

# build
cd ${VENDOR_DIR}/kustomize-kustomize-v3.4.0/kustomize; go build -o ${BIN_DIR}/kustomize