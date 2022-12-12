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
	"fmt"
	"os"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/zos"
	"github.com/sirupsen/logrus"
	"github.com/ysicing-cloud/sealos/cmd/flags"
	"github.com/ysicing-cloud/sealos/install"
	"github.com/ysicing-cloud/sealos/internal/pkg/util/factory"

	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	Info        bool
	globalFlags *flags.GlobalFlags
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// create a new factory
	f := factory.DefaultFactory()
	// build the root command
	rootCmd := BuildRoot(f)
	rootCmd.AddCommand(InitCmd(f))
	rootCmd.AddCommand(JoinCmd(f))
	rootCmd.AddCommand(CleanCmd(f))
	rootCmd.AddCommand(ConfigCmd(f))
	rootCmd.AddCommand(DeleteCmd(f))
	rootCmd.AddCommand(ExecCmd(f))
	rootCmd.AddCommand(InstallCmd(f))
	rootCmd.AddCommand(IPVSCmd(f))
	rootCmd.AddCommand(NewRouteCmd(f))
	if err := rootCmd.Execute(); err != nil {
		if globalFlags.Debug {
			f.GetLog().Fatalf("%v", err)
		} else {
			f.GetLog().Fatal(err)
		}
	}
}

// BuildRoot creates a new root command from the
func BuildRoot(f factory.Factory) *cobra.Command {
	rootCmd := NewRootCmd(f)
	persistentFlags := rootCmd.PersistentFlags()
	globalFlags = flags.SetGlobalFlags(persistentFlags)
	return rootCmd
}

// NewRootCmd returns a new root command
func NewRootCmd(f factory.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "sealos",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		PersistentPreRunE: func(cobraCmd *cobra.Command, args []string) error {
			if cobraCmd.Annotations != nil {
				return nil
			}
			qlog := f.GetLog()
			if globalFlags.Silent {
				qlog.SetLevel(logrus.FatalLevel)
			} else if globalFlags.Debug {
				qlog.SetLevel(logrus.DebugLevel)
			}
			// TODO apply extra flags
			return nil
		},
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home := zos.GetHomeDir()
	logHome := fmt.Sprintf("%s/.sealos/log", home)
	if !install.FileExist(logHome) {
		err := os.MkdirAll(logHome, os.ModePerm)
		if err != nil {
			fmt.Println("create default sealos config dir failed, please create it by your self mkdir -p /root/.sealos && touch /root/.sealos/config.yaml")
		}
	}
	log.StartFileLogging(logHome, "sealos.log")
}
