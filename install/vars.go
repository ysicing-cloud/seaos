// Copyright © 2021 sealos.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package install

import (
	"regexp"

	"github.com/sealyun/lvscare/care"

	"github.com/ysicing-cloud/sealos/ipvs"
	"github.com/ysicing-cloud/sealos/pkg/sshcmd/sshutil"
)

var (
	MasterIPs         []string
	NodeIPs           []string
	CertSANS          []string
	DNSDomain         string
	APIServerCertSANs []string
	SSHConfig         sshutil.SSH
	APIServer         string
	Token             string
	CriSocket         string
	CgroupDriver      string
	KubeadmAPI        string

	VIP     string
	PkgURL  string
	PodCIDR string
	SvcCIDR string

	Envs          []string // read env from -e
	PackageConfig string   // install/delete package config
	Values        string   // values for  install package values.yaml
	WorkDir       string   // workdir for install/delete package home

	Ipvs         care.LvsCare
	LvscareImage ipvs.LvscareImage

	Network string // network type, calico or flannel etc..

	WithoutCNI bool // if true don't install cni plugin

	Interface string //network interface name, like "eth.*|en.*"

	BGP bool // the ipip mode of the calico

	MTU string // mtu size

	YesRx = regexp.MustCompile("^(?i:y(?:es)?)$")

	CleanForce bool
	CleanAll   bool

	InDocker     bool
	SnapshotName string
	EtcdBackDir  string
	RestorePath  string

	OssEndpoint      string
	AccessKeyID      string
	AccessKeySecrets string
	BucketName       string
	ObjectPath       string
)

type metadata struct {
	K8sVersion string `json:"k8sVersion"`
	CniVersion string `json:"cniVersion"`
	CniName    string `json:"cniName"`
}
