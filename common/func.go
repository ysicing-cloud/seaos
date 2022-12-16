package common

import "github.com/ergoapi/util/zos"

func GetDefaultLogDir() string {
	home := zos.GetHomeDir()
	return home + "/" + DefaultLogDir
}

func GetDefaultConfig() string {
	home := zos.GetHomeDir()
	return home + "/" + DefaultCfgDir + "/cluster.yaml"
}
