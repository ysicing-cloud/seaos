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

func JoinCmd(f factory.Factory) *cobra.Command {
	slog := f.GetLog()
	joinCmd := &cobra.Command{
		Use:   "join",
		Short: "Simplest way to join your kubernets HA cluster",
		Long:  `sealos join --node 192.168.0.5`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(install.MasterIPs) == 0 && len(install.NodeIPs) == 0 {
				slog.Error("this command is join feature,master and node is empty at the same time.please check your args in command.")
				_ = cmd.Help()
				os.Exit(0)
			}
		},
		Run: JoinCmdFunc,
	}
	joinCmd.Flags().StringSliceVar(&install.MasterIPs, "master", []string{}, "kubernetes multi-master ex. 192.168.0.5-192.168.0.5")
	joinCmd.Flags().StringSliceVar(&install.NodeIPs, "node", []string{}, "kubernetes multi-nodes ex. 192.168.0.5-192.168.0.5")
	return joinCmd
}

func JoinCmdFunc(cmd *cobra.Command, args []string) {
	beforeNodes := install.ParseIPs(install.NodeIPs)
	beforeMasters := install.ParseIPs(install.MasterIPs)
	slog := log.GetInstance()
	c := &install.SealConfig{
		Log: slog,
	}
	if err := c.Load(cfgFile); err != nil {
		slog.Error(err)
		os.Exit(0)
	}

	cfgNodes := append(c.Masters, c.Nodes...)
	joinNodes := append(beforeNodes, beforeMasters...)

	if ok, node := deleteOrJoinNodeIsExistInCfgNodes(joinNodes, cfgNodes); ok {
		slog.Errorf(`[%s] has already exist in your cluster. please check.`, node)
		os.Exit(-1)
	}

	install.BuildJoin(beforeMasters, beforeNodes)
	c.Dump(cfgFile)
}
