package main

import (
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/ontio/ontology/common/log"
	"github.com/ontio/layer2deploy/cmd"
	"github.com/ontio/layer2deploy/layer2config"
)

func setupAPP() *cli.App {
	app := cli.NewApp()
	app.Usage = "layer2deploy CLI"
	app.Action = startLayer2Deploy
	app.Version = layer2config.Version
	app.Copyright = "Copyright in 2018 The Ontology Authors"
	app.Flags = []cli.Flag{
		cmd.LogLevelFlag,
		cmd.RestPortFlag,
		cmd.NetworkIdFlag,
		cmd.ConfigfileFlag,
	}
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

func main() {
	if err := setupAPP().Run(os.Args); err != nil {
		cmd.PrintErrorMsg(err.Error())
		os.Exit(1)
	}
}

func startLayer2Deploy(ctx *cli.Context) {
	initLog(ctx)
	if err := initConfig(ctx); err != nil {
		log.Errorf("[initConfig] error: %s", err)
		return
	}
	log.Infof("config: %v\n", layer2config.DefLayer2Config)
	waitToExit()
}

func initLog(ctx *cli.Context) {
	logLevel := ctx.GlobalInt(cmd.GetFlagName(cmd.LogLevelFlag))
	log.InitLog(logLevel, log.Stdout)
}

// GetPassword gets password from user input
func getDBPassword() ([]byte, error) {
	fmt.Printf("DB Password:")
	passwd, err := gopass.GetPasswd()
	if err != nil {
		return nil, err
	}
	return passwd, nil
}

func initConfig(ctx *cli.Context) error {
	//init config
	return cmd.SetOntologyConfig(ctx)
}

func waitToExit() {
	exit := make(chan bool, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range sc {
			log.Infof("saga server received exit signal: %s.", sig.String())
			close(exit)
			break
		}
	}()
	<-exit
}
