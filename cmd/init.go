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

package cmd

import (
	"os"

	"github.com/ergoapi/util/expass"
	"github.com/ergoapi/util/zos"
	"github.com/spf13/cobra"

	"github.com/ysicing-cloud/sealos/install"
	"github.com/ysicing-cloud/sealos/internal/pkg/util/factory"
)

var contact = `
      ___           ___           ___           ___       ___           ___     
     /\  \         /\  \         /\  \         /\__\     /\  \         /\  \    
    /::\  \       /::\  \       /::\  \       /:/  /    /::\  \       /::\  \   
   /:/\ \  \     /:/\:\  \     /:/\:\  \     /:/  /    /:/\:\  \     /:/\ \  \  
  _\:\~\ \  \   /::\~\:\  \   /::\~\:\  \   /:/  /    /:/  \:\  \   _\:\~\ \  \ 
 /\ \:\ \ \__\ /:/\:\ \:\__\ /:/\:\ \:\__\ /:/__/    /:/__/ \:\__\ /\ \:\ \ \__\
 \:\ \:\ \/__/ \:\~\:\ \/__/ \/__\:\/:/  / \:\  \    \:\  \ /:/  / \:\ \:\ \/__/
  \:\ \:\__\    \:\ \:\__\        \::/  /   \:\  \    \:\  /:/  /   \:\ \:\__\  
   \:\/:/  /     \:\ \/__/        /:/  /     \:\  \    \:\/:/  /     \:\/:/  /  
    \::/  /       \:\__\         /:/  /       \:\__\    \::/  /       \::/  /   
     \/__/         \/__/         \/__/         \/__/     \/__/         \/__/  

                  官方文档：sealyun.com
                  项目地址：github.com/ysicing-cloud/sealos
                  QQ群   ：98488045
                  常见问题：sealyun.com/faq
`

var exampleInit = `
	# init with password with three master one node
	sealos init --passwd your-server-password  \
	--master 192.168.0.2 --master 192.168.0.3 --master 192.168.0.4 \
	--node 192.168.0.5 --user root \
	--version v1.18.0 --pkg-url=/root/kube1.18.0.tar.gz 
	
	# init with pk-file , when your server have different password
	sealos init --pk /root/.ssh/id_rsa \
	--master 192.168.0.2 --node 192.168.0.5 --user root \
	--version v1.18.0 --pkg-url=/root/kube1.18.0.tar.gz 

	# when use multi network. set a can-reach with --interface 
 	sealos init --interface 192.168.0.254 \
	--master 192.168.0.2 --master 192.168.0.3 --master 192.168.0.4 \
	--node 192.168.0.5 --user root --passwd your-server-password \
	--version v1.18.0 --pkg-url=/root/kube1.18.0.tar.gz 
	
	# when your interface is not "eth*|en*|em*" like.
	sealos init --interface your-interface-name \
	--master 192.168.0.2 --master 192.168.0.3 --master 192.168.0.4 \
	--node 192.168.0.5 --user root --passwd your-server-password \
	--version v1.18.0 --pkg-url=/root/kube1.18.0.tar.gz
`

func InitCmd(f factory.Factory) *cobra.Command {
	slog := f.GetLog()
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Simplest way to init your kubernets HA cluster",
		Long: `sealos init --master 192.168.0.2 --master 192.168.0.3 --master 192.168.0.4 \
	--node 192.168.0.5 --user root --passwd your-server-password \
	--version v1.18.0 --pkg-url=/root/kube1.18.0.tar.gz`,
		Example: exampleInit,
		Run: func(cmd *cobra.Command, args []string) {
			c := &install.SealConfig{}
			// 没有重大错误可以直接保存配置. 但是apiservercertsans为空. 但是不影响用户 clean
			// 如果用户指定了配置文件,并不使用--master, 这里就不dump, 需要使用load获取配置文件了.
			if cfgFile != "" && len(install.MasterIPs) == 0 {
				err := c.Load(cfgFile)
				if err != nil {
					slog.Errorf("load cfgFile %s err: %q", cfgFile, err)
					os.Exit(1)
				}
			} else {
				c.Dump(cfgFile)
			}
			install.BuildInit()
			// 安装完成后生成完整版
			c.Dump(cfgFile)
			slog.Info(contact)
		},
	}

	// Here you will define your flags and configuration settings.
	initCmd.Flags().StringVar(&install.SSHConfig.User, "user", "root", "servers user name for ssh")
	initCmd.Flags().StringVar(&install.SSHConfig.Password, "passwd", "", "password for ssh")
	initCmd.Flags().StringVar(&install.SSHConfig.PkFile, "pk", zos.GetHomeDir()+"/.ssh/id_rsa", "private key for ssh")
	initCmd.Flags().StringVar(&install.SSHConfig.PkPassword, "pk-passwd", "", "private key password for ssh")

	initCmd.Flags().StringVar(&install.APIServer, "apiserver", "apiserver.cluster.local", "apiserver domain name")
	initCmd.Flags().StringVar(&install.Token, "token", expass.PwGenAlphaNum(16), "random token")
	initCmd.Flags().StringVar(&install.VIP, "vip", "10.103.97.2", "virtual ip")
	initCmd.Flags().StringSliceVar(&install.MasterIPs, "master", []string{}, "k3s multi-masters ex. 192.168.0.2-192.168.0.4")
	initCmd.Flags().StringSliceVar(&install.NodeIPs, "node", []string{}, "k3s multi-nodes ex. 192.168.0.5-192.168.0.5")
	initCmd.Flags().StringSliceVar(&install.CertSANS, "cert-sans", []string{}, "k3s apiServerCertSANs ex. 47.0.0.22 sealyun.com ")

	initCmd.Flags().StringVar(&install.PkgURL, "pkg-url", "", "http://store.lameleg.com/kube1.14.1.tar.gz download offline package url, or file location ex. /root/kube1.14.1.tar.gz")
	initCmd.Flags().StringVar(&install.PodCIDR, "podcidr", "100.64.0.0/10", "Specify range of IP addresses for the pod network")
	initCmd.Flags().StringVar(&install.SvcCIDR, "svccidr", "10.96.0.0/12", "Use alternative range of IP address for service VIPs")
	initCmd.Flags().StringVar(&install.LvscareImage.Image, "lvscare-image", "fanux/lvscare", "lvscare image name")
	initCmd.Flags().StringVar(&install.LvscareImage.Tag, "lvscare-tag", "latest", "lvscare image tag name")
	return initCmd
}
