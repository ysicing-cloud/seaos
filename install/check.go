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
)

// SetHosts set hosts. if can't access to hostName, set /etc/hosts
func SetHosts(hostIP, hostName string) {
	cmd := fmt.Sprintf("cat /etc/hosts |grep %s || echo '%s %s' >> /etc/hosts", hostName, IPFormat(hostIP), hostName)
	_ = SSHConfig.CmdAsync(hostIP, cmd)
}

// CheckValid is
func (s *SealosInstaller) CheckValid() {
	//hosts := append(Masters, Nodes...)
	// 所有master节点
	//masters := append(Masters, ParseIPs(MasterIPs)...)
	// 所有node节点
	//nodes := append(Nodes, ParseIPs(NodeIPs)...)
	//hosts := append(masters, nodes...)
	var hosts = append(s.Masters, s.Nodes...)
	if len(s.Hosts) == 0 && len(hosts) == 0 {
		s.Log.Error("hosts not allow empty")
		os.Exit(1)
	}
	if SSHConfig.User == "" {
		s.Log.Error("user not allow empty")
		os.Exit(1)
	}
	dict := make(map[string]bool)
	for _, h := range s.Hosts {
		hostname := SSHConfig.CmdToString(h, "hostname", "") //获取主机名
		if hostname == "" {
			s.Log.Errorf("[%s] ------------ check error", h)
			os.Exit(1)
		} else {
			SetHosts(h, hostname)
			if _, ok := dict[hostname]; !ok {
				dict[hostname] = true //不冲突, 主机名加入字典
			} else {
				s.Log.Error("duplicate hostnames is not allowed")
				os.Exit(1)
			}
			s.Log.Infof("[%s]  ------------ check ok", h)
		}
	}
}
