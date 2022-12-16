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
	"strings"

	"github.com/ergoapi/log"
	"github.com/sirupsen/logrus"
)

type CleanCluster interface {
	Check
	Clean
}

type JoinNodeAndMaster interface {
	Check
	Send
	Join
}

type Init interface {
	Check
	Send
	PreInit
	Join
	Print
}

type Install interface {
	Check
	Send
	Apply
}

var (
	JoinToken       string
	TokenCaCertHash string
	CertificateKey  string
)

// SealosInstaller is
type SealosInstaller struct {
	Hosts     []string
	Masters   []string
	Nodes     []string
	APIServer string
	Token     string
	Log       log.Logger `json:"-"`
}

func (s *SealosInstaller) Command() (cmd string) {
	return "systemctl enable k3s && systemctl start k3s"
}

// decode output to join token  hash and key
func decodeOutput(output []byte) {
	s0 := string(output)
	logrus.Debugf("[globals]decodeOutput: %s", s0)
	slice := strings.Split(s0, "kubeadm join")
	slice1 := strings.Split(slice[1], "Please note")
	logrus.Infof("[globals]join command is: %s", slice1[0])
	decodeJoinCmd(slice1[0])
}

// 192.168.0.200:6443 --token 9vr73a.a8uxyaju799qwdjv --discovery-token-ca-cert-hash sha256:7c2e69131a36ae2a042a339b33381c6d0d43887e2de83720eff5359e26aec866 --experimental-control-plane --certificate-key f8902e114ef118304e561c3ecd4d0b543adc226b7a07f675f56564185ffe0c07
func decodeJoinCmd(cmd string) {
	logrus.Debugf("[globals]decodeJoinCmd: %s", cmd)
	stringSlice := strings.Split(cmd, " ")

	for i, r := range stringSlice {
		r = strings.ReplaceAll(r, "\t", "")
		r = strings.ReplaceAll(r, "\n", "")
		r = strings.ReplaceAll(r, "\\", "")
		r = strings.TrimSpace(r)
		logrus.Debugf("[####]%d :%s:", i, r)
		// switch r {
		// case "--token":
		// 	JoinToken = stringSlice[i+1]
		// case "--discovery-token-ca-cert-hash":
		// 	TokenCaCertHash = stringSlice[i+1]
		// case "--certificate-key":
		// 	CertificateKey = stringSlice[i+1][:64]
		// }
		if strings.Contains(r, "--token") {
			JoinToken = stringSlice[i+1]
		}

		if strings.Contains(r, "--discovery-token-ca-cert-hash") {
			TokenCaCertHash = stringSlice[i+1]
		}

		if strings.Contains(r, "--certificate-key") {
			CertificateKey = stringSlice[i+1][:64]
		}
	}
	logrus.Debugf("[####]JoinToken :%s", JoinToken)
	logrus.Debugf("[####]TokenCaCertHash :%s", TokenCaCertHash)
	logrus.Debugf("[####]CertificateKey :%s", CertificateKey)
}
