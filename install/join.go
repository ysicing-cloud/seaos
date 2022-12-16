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
	"sync"

	"github.com/ergoapi/log"
	"github.com/ysicing-cloud/sealos/ipvs"
)

// BuildJoin is
func BuildJoin(joinMasters, joinNodes []string) {
	if len(joinMasters) > 0 {
		joinMastersFunc(joinMasters)
	}
	if len(joinNodes) > 0 {
		joinNodesFunc(joinNodes)
	}
}

func joinMastersFunc(joinMasters []string) {
	masters := MasterIPs
	nodes := NodeIPs
	i := &SealosInstaller{
		Hosts:     joinMasters,
		Masters:   masters,
		Nodes:     nodes,
		APIServer: APIServer,
		Log:       log.GetInstance(),
	}
	i.CheckValid()
	i.SendSealos()
	i.SendPackage()
	i.JoinMasters(joinMasters)
	//master join to MasterIPs
	MasterIPs = append(MasterIPs, joinMasters...)
	i.lvscare()
}

// 返回/etc/hosts记录
func getApiserverHost(ipAddr string) (host string) {
	return fmt.Sprintf("%s %s", ipAddr, APIServer)
}

// joinNodesFunc is join nodes func
func joinNodesFunc(joinNodes []string) {
	// 所有node节点
	nodes := joinNodes
	i := &SealosInstaller{
		Hosts:   nodes,
		Masters: MasterIPs,
		Nodes:   nodes,
		Log:     log.GetInstance(),
	}
	i.CheckValid()
	i.SendSealos()
	i.SendPackage()
	i.JoinNodes()
	//node join to NodeIPs
	NodeIPs = append(NodeIPs, joinNodes...)
}

// JoinMasters is
func (s *SealosInstaller) JoinMasters(masters []string) {
	var wg sync.WaitGroup
	//join master do sth
	cmd := s.Command()
	for _, master := range masters {
		wg.Add(1)
		go func(master string) {
			defer wg.Done()
			s.genMasterService(false, "/tmp/k3s.masterx.service")
			SSHConfig.CopyLocalToRemote(master, "/tmp/k3s.masterx.service", "/etc/systemd/system/k3s.service")
			cmdHosts := fmt.Sprintf("echo %s >> /etc/hosts", getApiserverHost(IPFormat(s.Masters[0])))
			_ = SSHConfig.CmdAsync(master, cmdHosts)
			// cmdMult := fmt.Sprintf("%s --apiserver-advertise-address %s", cmd, IpFormat(master))
			_ = SSHConfig.CmdAsync(master, cmd)
			// cmdHosts = fmt.Sprintf(`sed "s/%s/%s/g" -i /etc/hosts`, getApiserverHost(IPFormat(s.Masters[0])), getApiserverHost(IPFormat(master)))
			// _ = SSHConfig.CmdAsync(master, cmdHosts)
			copyk8sConf := `rm -rf .kube/config && mkdir -p /root/.kube && cp -a /etc/rancher/k3s/k3s.yaml /root/.kube/config && chmod 600 /root/.kube/config`
			_ = SSHConfig.CmdAsync(master, copyk8sConf)
		}(master)
	}
	wg.Wait()
}

// JoinNodes is
func (s *SealosInstaller) JoinNodes() {
	var masters string
	var wg sync.WaitGroup
	for _, master := range s.Masters {
		masters += fmt.Sprintf(" --rs %s:6443", IPFormat(master))
	}
	cmd := s.Command()
	ipvsCmd := fmt.Sprintf("sealos ipvs --vs %s:6443 %s --health-path /healthz --health-schem https --run-once", VIP, masters)
	for _, node := range s.Nodes {
		wg.Add(1)
		go func(node string) {
			defer wg.Done()
			s.genWorkerService("/tmp/k3s.worker.service")
			s.sendFile([]string{node}, "/tmp/k3s.worker.service", "/etc/systemd/system/k3s.service")

			cmdHosts := fmt.Sprintf("echo %s %s >> /etc/hosts", VIP, APIServer)
			_ = SSHConfig.CmdAsync(node, cmdHosts)

			// 如果不是默认路由， 则添加 vip 到 master的路由。
			cmdRoute := fmt.Sprintf("sealos route --host %s", IPFormat(node))
			status := SSHConfig.CmdToString(node, cmdRoute, "")
			if status != "ok" {
				// 以自己的ip作为路由网关
				addRouteCmd := fmt.Sprintf("sealos route add --host %s --gateway %s", VIP, IPFormat(node))
				SSHConfig.CmdToString(node, addRouteCmd, "")
			}

			_ = SSHConfig.CmdAsync(node, ipvsCmd) // create ipvs rules before we join node
			//create lvscare static pod
			yaml := ipvs.LvsStaticPodYaml(VIP, MasterIPs, LvscareImage)
			_ = SSHConfig.CmdAsync(node, cmd)
			_ = SSHConfig.Cmd(node, "mkdir -p /var/lib/rancher/k3s/agent/pod-manifests")
			SSHConfig.CopyConfigFile(node, "/var/lib/rancher/k3s/agent/pod-manifests/kube-sealyun-lvscare.yaml", []byte(yaml))

			cleaninstall := `rm -rf /tmp/package`
			_ = SSHConfig.CmdAsync(node, cleaninstall)
		}(node)
	}
	wg.Wait()
}

func (s *SealosInstaller) lvscare() {
	var wg sync.WaitGroup
	for _, node := range s.Nodes {
		wg.Add(1)
		go func(node string) {
			defer wg.Done()
			yaml := ipvs.LvsStaticPodYaml(VIP, MasterIPs, LvscareImage)
			_ = SSHConfig.Cmd(node, "rm -rf /var/lib/rancher/k3s/agent/pod-manifests/kube-sealyun-lvscare.yaml || :")
			SSHConfig.CopyConfigFile(node, "/var/lib/rancher/k3s/agent/pod-manifests/kube-sealyun-lvscare.yaml", []byte(yaml))
		}(node)
	}
	wg.Wait()
}
