package cmd

import (
	"github.com/ontio/layer2deploy/layer2config"
	"github.com/urfave/cli"
	"strings"
)

var (
	LogLevelFlag = cli.UintFlag{
		Name:  "loglevel",
		Usage: "Set the log level to `<level>` (0~6). 0:Trace 1:Debug 2:Info 3:Warn 4:Error 5:Fatal 6:MaxLevel",
		Value: uint(layer2config.DEFAULT_LOG_LEVEL),
	}
	RestPortFlag = cli.UintFlag{
		Name:  "restport",
		Usage: "restful server listening port `<number>`",
		Value: 0,
	}
	NetworkIdFlag = cli.UintFlag{
		Name:  "networkid",
		Usage: "Network id `<number>`. 1=ontology main net, 2=polaris test net, 3=testmode, and other for custom network",
		Value: layer2config.NETWORK_ID_SOLO_NET,
	}
	ConfigfileFlag = cli.StringFlag{
		Name:   "config",
		Usage:  "specify configfile",
		Value:  "config.json",
		EnvVar: "CONFIG_FILE",
	}
)

func GetFlagName(flag cli.Flag) string {
	name := flag.GetName()
	if name == "" {
		return ""
	}
	return strings.TrimSpace(strings.Split(name, ",")[0])
}
