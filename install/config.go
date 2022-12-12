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
	"fmt"
	"os"

	"sigs.k8s.io/yaml"

	"github.com/ergoapi/log"
	"github.com/ysicing-cloud/sealos/net"
)

const (
	defaultConfigPath      = "/.sealos"
	defaultConfigFile      = "/config.yaml"
	defaultAPIServerDomain = "apiserver.cluster.local"
)

// SealConfig for ~/.sealos/config.yaml
type SealConfig struct {
	Masters []string
	Nodes   []string
	//config from kubeadm.cfg. ex. cluster.local
	DNSDomain         string
	APIServerCertSANs []string

	//SSHConfig
	User       string
	Passwd     string
	PrivateKey string
	PkPassword string
	//ApiServer ex. apiserver.cluster.local
	APIServerDomain string
	Network         string
	VIP             string
	PkgURL          string
	Version         string
	Repo            string
	PodCIDR         string
	SvcCIDR         string
	//lvscare images
	LvscareName string
	LvscareTag  string

	Log log.Logger
}

// Dump is
func (c *SealConfig) Dump(path string) {
	home, _ := os.UserHomeDir()
	if path == "" {
		path = home + defaultConfigPath + defaultConfigFile
	}
	MasterIPs = ParseIPs(MasterIPs)
	c.Masters = MasterIPs
	NodeIPs = ParseIPs(NodeIPs)
	c.Nodes = ParseIPs(NodeIPs)
	c.User = SSHConfig.User
	c.Passwd = SSHConfig.Password
	c.PrivateKey = SSHConfig.PkFile
	c.PkPassword = SSHConfig.PkPassword
	c.APIServerDomain = APIServer
	c.Network = Network
	c.VIP = VIP
	c.PkgURL = PkgURL
	c.Version = Version
	c.Repo = Repo
	c.SvcCIDR = SvcCIDR
	c.PodCIDR = PodCIDR

	c.DNSDomain = DNSDomain
	c.APIServerCertSANs = APIServerCertSANs
	//lvscare
	c.LvscareName = LvscareImage.Image
	c.LvscareTag = LvscareImage.Tag
	y, err := yaml.Marshal(c)
	if err != nil {
		c.Log.Errorf("dump config file failed: %s", err)
	}

	err = os.MkdirAll(home+defaultConfigPath, os.ModePerm)
	if err != nil {
		c.Log.Warnf("create default sealos config dir failed, please create it by your self mkdir -p /root/.sealos && touch /root/.sealos/config.yaml")
	}

	if err = os.WriteFile(path, y, 0600); err != nil {
		c.Log.Warnf("write to file %s failed: %s", path, err)
	}
}

func Dump(path string, content interface{}) error {
	slog := log.GetInstance()
	y, err := yaml.Marshal(content)
	if err != nil {
		slog.Errorf("dump config file failed: %s", err)
		return err
	}
	home, _ := os.UserHomeDir()
	err = os.MkdirAll(home+defaultConfigPath, os.ModePerm)
	if err != nil {
		slog.Errorf("create dump dir failed %s", err)
		return err
	}

	_ = os.WriteFile(path, y, 0600)
	return nil
}

// Load is
func (c *SealConfig) Load(path string) (err error) {
	if path == "" {
		home, _ := os.UserHomeDir()
		path = home + defaultConfigPath + defaultConfigFile
	}

	y, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config file %s failed %w", path, err)
	}

	err = yaml.Unmarshal(y, c)
	if err != nil {
		return fmt.Errorf("unmarshal config file failed: %w", err)
	}

	MasterIPs = c.Masters
	NodeIPs = c.Nodes
	SSHConfig.User = c.User
	SSHConfig.Password = c.Passwd
	SSHConfig.PkFile = c.PrivateKey
	SSHConfig.PkPassword = c.PkPassword
	APIServer = c.APIServerDomain
	Network = c.Network
	VIP = c.VIP
	PkgURL = c.PkgURL
	Version = c.Version
	Repo = c.Repo
	PodCIDR = c.PodCIDR
	SvcCIDR = c.SvcCIDR
	DNSDomain = c.DNSDomain
	APIServerCertSANs = c.APIServerCertSANs
	//lvscare
	LvscareImage.Image = c.LvscareName
	LvscareImage.Tag = c.LvscareTag
	return
}

func Load(path string, content interface{}) error {
	slog := log.GetInstance()
	y, err := os.ReadFile(path)
	if err != nil {
		slog.Errorf("read config file %s failed %s", path, err)
		os.Exit(0)
	}

	err = yaml.Unmarshal(y, content)
	if err != nil {
		slog.Errorf("unmarshal config file failed: %s", err)
	}
	return nil
}

func (c *SealConfig) ShowDefaultConfig() {
	home, _ := os.UserHomeDir()
	c.Masters = []string{"192.168.0.2", "192.168.0.2", "192.168.0.2"}
	c.Nodes = []string{"192.168.0.3", "192.168.0.4"}
	c.User = "root"
	c.Passwd = "123456"
	c.PrivateKey = home + "/.ssh/id_rsa"
	c.APIServerDomain = defaultAPIServerDomain
	c.Network = net.CALICO
	c.VIP = "10.103.97.2"
	c.PkgURL = home + "/kube1.17.13.tar.gz"
	c.Version = "v1.17.13"
	c.Repo = "k8s.gcr.io"
	c.PodCIDR = "100.64.0.0/10"
	c.SvcCIDR = "10.96.0.0/12"
	c.APIServerCertSANs = []string{"apiserver.cluster.local", "127.0.0.1"}
	c.LvscareName = "fanux/lvscare"
	c.LvscareTag = "latest"

	y, err := yaml.Marshal(c)
	if err != nil {
		c.Log.Errorf("marshal config file failed: %s", err)
	}

	c.Log.Infof("\n\n%s\n\n", string(y))
	c.Log.Infof("Please save above config in ~/.sealos/config.yaml and edit values on your own")
}
