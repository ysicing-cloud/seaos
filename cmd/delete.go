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

	"github.com/ysicing-cloud/sealos/internal/pkg/util/factory"

	"github.com/spf13/cobra"

	"github.com/ysicing-cloud/sealos/install"
	"github.com/ysicing-cloud/sealos/pkg/appmanager"
)

var (
	deleteExamlpe = `
	# delete from default config in installed tar
	sealos delete --pkg-url mysql.tar

	# delete from your define config when sealos install use
	sealos delete --pkg-url mysql.tar -c /root/config -w /data/

	# delete force no interactive
	sealos delete --pkg-url mysql.tar -f
`
)

func DeleteCmd(f factory.Factory) *cobra.Command {
	slog := f.GetLog()
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "delete kubernetes apps installled by sealos..",
		Long: `delete kubernetes apps, like dashboard prometheus installled by sealos.. 
`,
		Example: deleteExamlpe,
		Run: func(cmd *cobra.Command, args []string) {
			cfg := appmanager.GetDeleteFlags(AppURL)
			_ = appmanager.DeleteApp(cfg, cfgFile)
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			slog.Fatal("the delete app feature not support")
			if install.ExitDeleteCase(AppURL) {
				_ = cmd.Help()
				os.Exit(install.ErrorExitOSCase)
			}
		},
	}

	deleteCmd.Flags().StringVar(&AppURL, "pkg-url", "", "APP offline pluginsfile localtion ex. /root/prometheus.tar.gz")
	deleteCmd.Flags().StringVarP(&install.PackageConfig, "pkg-config", "c", "", "packageConfig for delete installed package config")
	deleteCmd.Flags().StringVarP(&install.WorkDir, "workdir", "w", "/root", "workdir for install package home, keep the same with installed")
	deleteCmd.Flags().BoolVarP(&install.CleanForce, "force", "f", false, "if this is true, will no prompt")
	return deleteCmd
}
