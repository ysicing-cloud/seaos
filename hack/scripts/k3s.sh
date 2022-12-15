#!/bin/bash

export INSTALL_K3S_MIRROR=cn

curl -sfL https://rancher-mirror.oss-cn-beijing.aliyuncs.com/k3s/k3s-install.sh | K3S_TOKEN=SECRET sh -s - server --cluster-init
curl -sfL https://rancher-mirror.oss-cn-beijing.aliyuncs.com/k3s/k3s-install.sh | K3S_TOKEN=SECRET sh -s - server --server https://<ip or hostname of server1>:6443
