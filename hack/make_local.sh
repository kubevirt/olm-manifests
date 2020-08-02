#!/usr/bin/env bash
set -ex

LOCAL_DIR=_local
FORMAT=${FORMAT:-txt}

set -o allexport
source hack/config
set +o allexport

mkdir -p "${LOCAL_DIR}"
./hack/make_local.py "${LOCAL_DIR}" "${FORMAT}"
sed "s/\(^.*\/operator.yaml$\)/### \1/" deploy/deploy.sh > _local/deploy.sh
sed -i "s|^kubectl create -f https://raw.githubusercontent.com/kubevirt/hyperconverged-cluster-operator/master/deploy|kubectl apply -f deploy|g" _local/deploy.sh
chmod +x _local/deploy.sh

_local/deploy.sh
kubectl apply -f _local/local.yaml
