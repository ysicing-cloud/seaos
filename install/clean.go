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
	"strings"
	"sync"

	"github.com/ergoapi/log"
	"github.com/ysicing-cloud/sealos/ipvs"
	sshcmd "github.com/ysicing-cloud/sealos/pkg/sshcmd/cmd"
)

type SealosClean struct {
	SealosInstaller
	cleanAll bool
}

// BuildClean clean the build resources.
func BuildClean(deleteNodes, deleteMasters []string) {
	i := &SealosClean{cleanAll: false, SealosInstaller: SealosInstaller{Log: log.GetInstance()}}
	masters := MasterIPs
	nodes := NodeIPs
	//1. 删除masters
	if len(deleteMasters) != 0 {
		if !CleanForce { // false
			prompt := fmt.Sprintf("clean command will clean masters [%s], continue clean (y/n)?", strings.Join(deleteMasters, ","))
			result := Confirm(prompt)
			if !result {
				i.Log.Debug("clean masters command is skip")
				goto node
			}
		}
		//只清除masters
		i.Masters = deleteMasters
	}

	//2. 删除nodes
node:
	if len(deleteNodes) != 0 {
		if !CleanForce { // flase
			prompt := fmt.Sprintf("clean command will clean nodes [%s], continue clean (y/n)?", strings.Join(deleteNodes, ","))
			result := Confirm(prompt)
			if !result {
				i.Log.Debug("clean nodes command is skip")
				goto all
			}
		}
		//只清除nodes
		i.Nodes = deleteNodes
	}
	//3. 删除所有节点
all:
	if len(deleteNodes) == 0 && len(deleteMasters) == 0 && CleanAll {
		if !CleanForce { // flase
			result := Confirm(`clean command will clean all masters and nodes, continue clean (y/n)?`)
			if !result {
				i.Log.Debug("clean all node command is skip")
				goto end
			}
		}
		// 所有master节点
		i.Masters = masters
		// 所有node节点
		i.Nodes = nodes
		i.cleanAll = true
	}
end:
	if len(i.Masters) == 0 && len(i.Nodes) == 0 {
		i.Log.Warn("clean nodes and masters is empty,please check your args and config.yaml.")
		os.Exit(-1)
	}
	i.CheckValid()
	i.Clean()
	if i.cleanAll {
		i.Log.Info("if clean all and clean sealos config")
		home, _ := os.UserHomeDir()
		cfgPath := home + defaultConfigPath
		sshcmd.Cmd("/bin/sh", "-c", "rm -rf "+cfgPath)
	}
}

// Clean clean cluster.
func (s *SealosClean) Clean() {
	var wg sync.WaitGroup
	//s 是要删除的数据
	//全局是当前的数据
	if len(s.Nodes) > 0 {
		//1. 再删除nodes
		for _, node := range s.Nodes {
			wg.Add(1)
			go func(node string) {
				defer wg.Done()
				s.cleanNode(node)
			}(node)
		}
		wg.Wait()
	}
	if len(s.Masters) > 0 {
		//2. 先删除master
		lock := sync.Mutex{}
		for _, master := range s.Masters {
			wg.Add(1)
			go func(master string) {
				lock.Lock()
				defer lock.Unlock()
				defer wg.Done()
				s.cleanMaster(master)
			}(master)
		}
		wg.Wait()
	}
}

func (s *SealosClean) cleanNode(node string) {
	cleanRoute(node)
	clean(node)
	//remove node
	NodeIPs = SliceRemoveStr(NodeIPs, node)
	if !s.cleanAll {
		s.Log.Debug("clean node in master")
		if len(MasterIPs) > 0 {
			hostname := isHostName(MasterIPs[0], node)
			cmd := "kubectl delete node %s"
			_ = SSHConfig.CmdAsync(MasterIPs[0], fmt.Sprintf(cmd, strings.TrimSpace(hostname)))
		}
	}
}

func (s *SealosClean) cleanMaster(master string) {
	clean(master)
	//remove master
	MasterIPs = SliceRemoveStr(MasterIPs, master)
	if !s.cleanAll {
		s.Log.Debug("clean node in master")
		if len(MasterIPs) > 0 {
			hostname := isHostName(MasterIPs[0], master)
			cmd := "kubectl delete node %s"
			_ = SSHConfig.CmdAsync(MasterIPs[0], fmt.Sprintf(cmd, strings.TrimSpace(hostname)))
		}
		//清空所有的nodes的数据
		yaml := ipvs.LvsStaticPodYaml(VIP, MasterIPs, LvscareImage)
		var wg sync.WaitGroup
		for _, node := range NodeIPs {
			wg.Add(1)
			go func(node string) {
				defer wg.Done()
				_ = SSHConfig.CmdAsync(node, fmt.Sprintf("mkdir -p /var/lib/rancher/k3s/agent/pod-manifests && echo '%s' > /var/lib/rancher/k3s/agent/pod-manifests/kube-sealyun-lvscare.yaml", yaml))
			}(node)
		}
		wg.Wait()
	}
}

func clean(host string) {
	cmd := "modprobe -r ipip  && lsmod"
	_ = SSHConfig.CmdAsync(host, cmd)
	cmd = "rm -rf ~/.kube"
	_ = SSHConfig.CmdAsync(host, cmd)
	cmd = "systemctl stop k3s && systemctl disable k3s && rm -rf /etc/systemd/system/k3s.service"
	_ = SSHConfig.CmdAsync(host, cmd)
	cmd = "rm -rf /var/lib/rancher /etc/rancher"
	_ = SSHConfig.CmdAsync(host, cmd)
	cmd = fmt.Sprintf("sed -i \"/%s/d\" /etc/hosts ", APIServer)
	_ = SSHConfig.CmdAsync(host, cmd)
}

func cleanRoute(node string) {
	// clean route
	cmdRoute := fmt.Sprintf("sealos route --host %s", IPFormat(node))
	status := SSHConfig.CmdToString(node, cmdRoute, "")
	if status != "ok" {
		// 删除为 vip创建的路由。
		delRouteCmd := fmt.Sprintf("sealos route del --host %s --gateway %s", VIP, IPFormat(node))
		SSHConfig.CmdToString(node, delRouteCmd, "")
	}
}
