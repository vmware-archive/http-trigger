#!/usr/bin/env bash

# Copyright (c) 2016-2017 Bitnami
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

KUBELESS_MANIFESTS=${KUBELESS_MANIFESTS:?}
CONTROLLER_TAG=${CONTROLLER_TAG:?}

set -e

CUR_DIR="$( cd "$(dirname "$0")" >/dev/null 2>&1 && pwd -P )"
WORKSPACE_DIR="${CUR_DIR}/../"

git clone https://github.com/kubeless/kubeless.git ${GOPATH}/src/github.com/kubeless/kubeless
ln -s ${WORKSPACE_DIR}/ksonnet-lib ${GOPATH}/src/github.com/kubeless/kubeless/ksonnet-lib
cd ${GOPATH}/src/github.com/kubeless/kubeless
make binary
make all-yaml

# Replace the controller version that is included in the main Kubeless manifest
mkdir -p ${WORKSPACE_DIR}/build-manifests/
IFS=' ' read -r -a manifests <<< "$KUBELESS_MANIFESTS"
for f in "${manifests[@]}"; do
    sed -E -i.bak 's/image: .*\/http-trigger-controller:.*/image: kubeless\/http-trigger-controller:'"${CONTROLLER_TAG}"'/g' ${f}.yaml
    cp ${f}.yaml ${WORKSPACE_DIR}/build-manifests/
done
