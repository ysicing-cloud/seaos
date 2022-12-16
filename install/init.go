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

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/file"
	"github.com/ysicing-cloud/sealos/internal/pkg/k3s"
)

// BuildInit is
func BuildInit() {
	slog := log.GetInstance()
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
		Token:     Token,
		Log:       slog,
	}
	i.CheckValid()
	i.SendSealos()
	i.SendPackage()
	i.InstallMaster0()
	if len(masters) > 1 {
		i.JoinMasters(i.Masters[1:])
	}
	if len(nodes) > 0 {
		i.JoinNodes()
	}
}

func getDefaultSANs() string {
	return fmt.Sprintf("apiserver.cluster.local,%v", VIP)
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

func (s *SealosInstaller) genMasterService(master0 bool, path string) error {
	s.Log.Info("gen master k3s service")
	tpl := k3s.NewService("master", k3s.MetaData{Master0: master0, TLSSAN: getDefaultSANs(), Docker: true, Server: "https://apiserver.cluster.local:6443", Token: s.Token}).Template()
	return file.WriteToFile(path, []byte(tpl))
}

func (s *SealosInstaller) genWorkerService(path string) error {
	s.Log.Info("gen worker k3s service")
	tpl := k3s.NewService("worker", k3s.MetaData{Docker: true, Server: "https://apiserver.cluster.local:6443", Token: s.Token}).Template()
	return file.WriteToFile(path, []byte(tpl))
}

// InstallMaster0 is
func (s *SealosInstaller) InstallMaster0() {
	s.Log.Info("init first master")
	s.genMasterService(true, "/tmp/k3s.master0.service")
	s.sendFile([]string{s.Masters[0]}, "/tmp/k3s.master0.service", "/etc/systemd/system/k3s.service")

	// remote server run sealos init . it can not reach apiserver.cluster.local , should add masterip apiserver.cluster.local to /etc/hosts
	// err := s.appendAPIServer()
	// if err != nil {
	// 	s.Log.Warnf("append  %s %s to /etc/hosts err: %s", IPFormat(s.Masters[0]), APIServer, err)
	// }
	// //master0 do sth
	// cmd := fmt.Sprintf("grep -qF '%s %s' /etc/hosts || echo %s %s >> /etc/hosts", IPFormat(s.Masters[0]), APIServer, IPFormat(s.Masters[0]), APIServer)
	// _ = SSHConfig.CmdAsync(s.Masters[0], cmd)
	cmd := s.Command()
	SSHConfig.Cmd(s.Masters[0], cmd)
	cmd = fmt.Sprintf("mkdir -p /root/.kube && cp /etc/rancher/k3s/k3s.yaml /root/.kube/config && chmod 600 /root/.kube/config")
	SSHConfig.Cmd(s.Masters[0], cmd)
}
