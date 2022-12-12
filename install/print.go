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
	"encoding/json"
	"strings"
)

// Print is
func (s *SealosInstaller) Print(process ...string) {
	if len(process) == 0 {
		configJSON, _ := json.Marshal(s)
		s.Log.Infof("\n[globals]sealos config is: %s", string(configJSON))
	} else {
		var sb strings.Builder
		for _, v := range process {
			sb.Write([]byte("==>"))
			sb.Write([]byte(v))
		}
		s.Log.Debug(sb.String())
	}
}
func (s *SealosInstaller) PrintFinish() {
	s.Log.Info("sealos install success.")
}
