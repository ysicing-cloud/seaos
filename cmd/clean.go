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

	"github.com/ergoapi/log"
	"github.com/spf13/cobra"
	"github.com/ysicing-cloud/sealos/install"
	"github.com/ysicing-cloud/sealos/internal/pkg/util/factory"
)

var exampleCleanCmd = `
	# clean  master
	sealos clean --master 192.168.0.2 \
	--master 192.168.0.3

	# clean  node  use --force/-f will be not prompt 
	sealos clean --node 192.168.0.4 \
	--node 192.168.0.5 --force

	# clean master and node
	sealos clean --master 192.168.0.2-192.168.0.3 \
 	--node 192.168.0.4-192.168.0.5

	# clean your kubernets HA cluster and use --force/-f will be not prompt (danger)
	sealos clean --all -f
`

func CleanCmd(f factory.Factory) *cobra.Command {
	cleanCmd := &cobra.Command{
		Use:     "clean",
		Short:   "Simplest way to clean your kubernets HA cluster",
		Long:    `sealos clean`,
		Example: exampleCleanCmd,
		Run:     CleanCmdFunc,
	}
	// Here you will define your flags and configuration settings.
	cleanCmd.Flags().StringSliceVar(&install.NodeIPs, "node", []string{}, "clean node ips.kubernetes multi-nodes ex. 192.168.0.5-192.168.0.5")
	cleanCmd.Flags().StringSliceVar(&install.MasterIPs, "master", []string{}, "clean master ips.kubernetes multi-nodes ex. 192.168.0.5-192.168.0.5")
	cleanCmd.PersistentFlags().BoolVarP(&install.CleanForce, "force", "f", false, "if this is true, will no prompt")
	cleanCmd.PersistentFlags().BoolVar(&install.CleanAll, "all", false, "if this is true, delete all ")
	return cleanCmd
}

func CleanCmdFunc(cmd *cobra.Command, args []string) {
	slog := log.GetInstance()
	deleteNodes := install.ParseIPs(install.NodeIPs)
	deleteMasters := install.ParseIPs(install.MasterIPs)
	c := &install.SealConfig{}
	if err := c.Load(cfgFile); err != nil {
		slog.Error(err)
		os.Exit(-1)
	}

	// 使用 sealos clean --node   不小心写了 masterip.
	if ok, node := deleteOrJoinNodeIsExistInCfgNodes(deleteNodes, c.Masters); ok {
		slog.Errorf(`clean master Use "sealos clean --master %s" to clean it, exit...`, node)
		os.Exit(-1)
	}
	// 使用 sealos clean --master 不小心写了 nodeip.
	if ok, node := deleteOrJoinNodeIsExistInCfgNodes(deleteMasters, c.Nodes); ok {
		slog.Errorf(`clean nodes Use "sealos clean --node %s" to clean it, exit...`, node)
		os.Exit(-1)
	}

	install.BuildClean(deleteNodes, deleteMasters)
	c.Dump(cfgFile)
}

// IsExistNodes
func deleteOrJoinNodeIsExistInCfgNodes(deleteOrJoinNodes []string, nodes []string) (bool, string) {
	for _, node := range nodes {
		for _, deleteOrJoinNode := range deleteOrJoinNodes {
			// 如果ips 相同. 则说明cfg配置文件已经存在该node.
			if node == deleteOrJoinNode {
				return true, node
			}
		}
	}
	return false, ""
}
