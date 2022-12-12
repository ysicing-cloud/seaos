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
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
)

var ConfigType string

func Config() {
	switch ConfigType {
	case "kubeadm":
		printlnKubeadmConfig()
	case "join":
		printlnJoinKubeadmConfig()
	default:
		printlnKubeadmConfig()
	}
}

func joinKubeadmConfig() string {
	var sb strings.Builder
	sb.Write([]byte(JoinCPTemplateText))
	return sb.String()
}

func printlnJoinKubeadmConfig() {
	fmt.Println(joinKubeadmConfig())
}

func kubeadmConfig() string {
	var sb strings.Builder
	sb.Write([]byte(InitTemplateText))
	return sb.String()
}

func printlnKubeadmConfig() {
	fmt.Println(kubeadmConfig())
}

// Template is
func Template() []byte {
	return TemplateFromTemplateContent(kubeadmConfig())
}

// JoinTemplate is generate JoinCP nodes configuration by master ip.
func JoinTemplate(ip string, cgroup string) []byte {
	return JoinTemplateFromTemplateContent(joinKubeadmConfig(), ip, cgroup)
}

func JoinTemplateFromTemplateContent(templateContent, ip, cgroup string) []byte {
	tmpl, err := template.New("text").Parse(templateContent)
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("join template parse failed: %v", err)
		}
	}()
	if err != nil {
		panic(1)
	}
	var envMap = make(map[string]interface{})
	envMap["Master0"] = IPFormat(MasterIPs[0])
	envMap["Master"] = ip
	envMap["TokenDiscovery"] = JoinToken
	envMap["TokenDiscoveryCAHash"] = TokenCaCertHash
	envMap["VIP"] = VIP
	envMap["KubeadmApi"] = KubeadmAPI
	envMap["CriSocket"] = CriSocket
	envMap["CgroupDriver"] = cgroup
	var buffer bytes.Buffer
	_ = tmpl.Execute(&buffer, envMap)
	return buffer.Bytes()
}

func TemplateFromTemplateContent(templateContent string) []byte {
	tmpl, err := template.New("text").Parse(templateContent)
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("template parse failed: %v", err)
		}
	}()
	if err != nil {
		panic(1)
	}
	var masters []string
	getmasters := MasterIPs
	for _, h := range getmasters {
		masters = append(masters, IPFormat(h))
	}
	var envMap = make(map[string]interface{})
	envMap["CertSANS"] = CertSANS
	envMap["VIP"] = VIP
	envMap["Masters"] = masters
	envMap["ApiServer"] = APIServer
	envMap["PodCIDR"] = PodCIDR
	envMap["SvcCIDR"] = SvcCIDR
	envMap["Master0"] = IPFormat(MasterIPs[0])
	envMap["Network"] = Network
	envMap["CgroupDriver"] = CgroupDriver
	envMap["KubeadmApi"] = KubeadmAPI
	envMap["CriSocket"] = CriSocket
	var buffer bytes.Buffer
	_ = tmpl.Execute(&buffer, envMap)
	return buffer.Bytes()
}
