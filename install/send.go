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
	"path"
	"sync"
)

// SendPackage is
func (s *SealosInstaller) SendPackage() {
	pkg := path.Base(PkgURL)
	afterHook := fmt.Sprintf("cd /tmp && tar zxvf %s  && cd /tmp/package/script && bash init.sh && bash prehook.sh", pkg)
	PkgURL = SendPackage(PkgURL, s.Hosts, "/tmp", nil, &afterHook)
}

// SendSealos is send the exec sealos to /usr/bin/sealos
func (s *SealosInstaller) SendSealos() {
	// send sealos first to avoid old version
	sealos := FetchSealosAbsPath()
	beforeHook := "ps -ef |grep -v 'grep'|grep sealos >/dev/null || rm -rf /usr/bin/sealos"
	SendPackage(sealos, s.Hosts, "/usr/bin", &beforeHook, nil)
}

func (s *SealosInstaller) sendFile(hosts []string, srcfile, dstfile string) {
	var wg sync.WaitGroup
	for _, node := range hosts {
		wg.Add(1)
		go func(node string) {
			defer wg.Done()
			SSHConfig.CopyLocalToRemote(node, srcfile, dstfile)
		}(node)
	}
	wg.Wait()
}
