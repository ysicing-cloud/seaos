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
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// BuildInit is
func BuildInit() {
	MasterIPs = ParseIPs(MasterIPs)
	NodeIPs = ParseIPs(NodeIPs)
	// 所有master节点
	masters := MasterIPs
	// 所有node节点
	nodes := NodeIPs
	hosts := append(masters, nodes...)
	i := &SealosInstaller{
		Hosts:     hosts,
		Masters:   masters,
		Nodes:     nodes,
		APIServer: APIServer,
	}
	i.CheckValid()
	i.Print()
	i.SendSealos()
	i.SendPackage()
	i.Print("SendPackage")

	i.InstallMaster0()
	i.Print("SendPackage", "KubeadmConfigInstall", "InstallMaster0")
	if len(masters) > 1 {
		i.JoinMasters(i.Masters[1:])
		i.Print("SendPackage", "KubeadmConfigInstall", "InstallMaster0", "JoinMasters")
	}
	if len(nodes) > 0 {
		i.JoinNodes()
		i.Print("SendPackage", "KubeadmConfigInstall", "InstallMaster0", "JoinMasters", "JoinNodes")
	}
	i.PrintFinish()
}

func getDefaultSANs() []string {
	var sans = []string{"127.0.0.1", "apiserver.cluster.local", VIP}
	// 指定的certSANS不为空, 则添加进去
	if len(CertSANS) != 0 {
		sans = append(sans, CertSANS...)
	}
	for _, master := range MasterIPs {
		sans = append(sans, IPFormat(master))
	}
	return sans
}

func (s *SealosInstaller) appendAPIServer() error {
	etcHostPath := "/etc/hosts"
	etcHostMap := fmt.Sprintf("%s %s", IPFormat(s.Masters[0]), APIServer)
	file, err := os.OpenFile(etcHostPath, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		s.Log.Errorf("open %s file error %s", etcHostPath, err)
		os.Exit(1)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		str, err := reader.ReadString('\n')
		if strings.Contains(str, APIServer) {
			s.Log.Infof("local %s is already exists %s", etcHostPath, APIServer)
			return nil
		}
		if err == io.EOF {
			break
		}
	}
	write := bufio.NewWriter(file)
	_, _ = write.WriteString(etcHostMap)
	return write.Flush()
}

// InstallMaster0 is
func (s *SealosInstaller) InstallMaster0() {
	s.SendKubeConfigs([]string{s.Masters[0]})
	s.sendNewCertAndKey([]string{s.Masters[0]})

	// remote server run sealos init . it can not reach apiserver.cluster.local , should add masterip apiserver.cluster.local to /etc/hosts
	err := s.appendAPIServer()
	if err != nil {
		s.Log.Warnf("append  %s %s to /etc/hosts err: %s", IPFormat(s.Masters[0]), APIServer, err)
	}
	//master0 do sth
	cmd := fmt.Sprintf("grep -qF '%s %s' /etc/hosts || echo %s %s >> /etc/hosts", IPFormat(s.Masters[0]), APIServer, IPFormat(s.Masters[0]), APIServer)
	_ = SSHConfig.CmdAsync(s.Masters[0], cmd)

	cmd = s.Command("", InitMaster)

	output := SSHConfig.Cmd(s.Masters[0], cmd)
	if output == nil {
		s.Log.Errorf("[%s] install kubernetes failed. please clean and uninstall.", s.Masters[0])
		os.Exit(1)
	}
	decodeOutput(output)

	cmd = `mkdir -p /root/.kube && cp /etc/kubernetes/admin.conf /root/.kube/config && chmod 600 /root/.kube/config`
	SSHConfig.Cmd(s.Masters[0], cmd)
}

// SendKubeConfigs
func (s *SealosInstaller) SendKubeConfigs(masters []string) {
	s.sendKubeConfigFile(masters, "kubelet.conf")
	s.sendKubeConfigFile(masters, "admin.conf")
	s.sendKubeConfigFile(masters, "controller-manager.conf")
	s.sendKubeConfigFile(masters, "scheduler.conf")
}

func (s *SealosInstaller) SendJoinMasterKubeConfigs(masters []string) {
	s.sendKubeConfigFile(masters, "admin.conf")
	s.sendKubeConfigFile(masters, "controller-manager.conf")
	s.sendKubeConfigFile(masters, "scheduler.conf")
}
