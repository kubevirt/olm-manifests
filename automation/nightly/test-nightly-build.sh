#!/usr/bin/env bash

set -ex

# Get golang
docker login --username "$(cat "${DOCKER_USER}")" --password-stdin < "${DOCKER_PASSWORD}"
wget -q https://dl.google.com/go/go1.15.2.linux-amd64.tar.gz
tar -C /usr/local -xf go*.tar.gz
export PATH=/usr/local/go/bin:$PATH

# get latest commits
latest_kubevirt=$(curl -sL https://storage.googleapis.com/kubevirt-prow/devel/nightly/release/kubevirt/kubevirt/latest)
latest_kubevirt_image=$(curl -sL "https://storage.googleapis.com/kubevirt-prow/devel/nightly/release/kubevirt/kubevirt/${latest_kubevirt}/kubevirt-operator.yaml" | grep 'OPERATOR_IMAGE' -A1 | tail -n 1 | sed 's/.*value: //g')
IFS=: read -r kv_image kv_tag <<< "${latest_kubevirt_image}"
latest_kubevirt_commit=$(curl -sL "https://storage.googleapis.com/kubevirt-prow/devel/nightly/release/kubevirt/kubevirt/${latest_kubevirt}/commit")

# Update HCO dependencies
go mod edit -require "kubevirt.io/kubevirt@${latest_kubevirt_commit}"
go mod vendor
rm -rf kubevirt

# Get latest kubevirt
git clone https://github.com/kubevirt/kubevirt.git
(cd kubevirt; git checkout "${latest_kubevirt_commit}")
go mod edit -replace kubevirt.io/client-go=./kubevirt/staging/src/kubevirt.io/client-go
go mod vendor
go mod tidy

# set envs
build_date="$(date +%Y%m%d)"
export IMAGE_REGISTRY=docker.io
export DOCKER_PREFIX=kubevirtnightlybuilds
export OPERATOR_IMAGE=${DOCKER_PREFIX}/hyperconverged-cluster-operator
export WEBHOOK_IMAGE=${DOCKER_PREFIX}/hyperconverged-cluster-webhook
export IMAGE_TAG
IMAGE_TAG="${build_date}_$(git show -s --format=%h)"

# Build HCO & HCO Webhook
make container-build-operator container-push-operator container-build-webhook container-push-webhook

#update env file before re-building menifests
sed -i "s#docker.io/kubevirt/virt-#${kv_image/-*/-}#" deploy/images.csv
sed -i "s#^KUBEVIRT_VERSION=.*#KUBEVIRT_VERSION=\"${kv_tag}\"#" hack/config
sed -i "s#^CSV_VERSION=.*#CSV_VERSION=\"${IMAGE_TAG}\"#" hack/config

(
  cd ./tools/digester
  go build .
)
./automation/digester/update_images.sh
./hack/build-manifests.sh

REGISTRY_NAMESPACE=${DOCKER_PREFIX} CONTAINER_TAG=${IMAGE_TAG} make bundleRegistry
hco_bucket="kubevirt-prow/devel/nightly/release/kubevirt/hyperconverged-cluster-operator"
echo "${build_date}" > build-date
echo "${IMAGE_REGISTRY}/${DOCKER_PREFIX}/hco-container-registry:${IMAGE_TAG}" > hco-bundle
gsutil cp ./hco-bundle "gs://${hco_bucket}/${build_date}/hco-bundle-image"
gsutil cp ./build-date gs://${hco_bucket}/latest
