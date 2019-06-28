package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/yeqown/infrastructure/types"
	"github.com/yeqown/micro-server-demo/global"

	_ "github.com/go-sql-driver/mysql"
	"github.com/urfave/cli"
	"github.com/yeqown/infrastructure/pkg/cfgutil"
)

func init() {
	var (
		curEnv types.Envrion
		cfg    = new(global.Config)
	)

	env := strings.TrimSpace(os.Getenv("ENV"))
	curEnv = types.ParseEnvrion(env)
	global.SetEnv(curEnv)

	rcloser, err := cfgutil.Open(
		fmt.Sprintf("configs/%s.json", curEnv.String()),
	)
	if err != nil {
		log.Fatalf("could not open config file: %v", err)
	}

	// load JSON config file depends on `curEnv`
	if err := cfgutil.LoadJSON(rcloser, cfg); err != nil {
		panic(err)
	}
	defer rcloser.Close()

	global.SetConfig(cfg)
}

// RUN GRPC SERVER and REST-HTTP SERVER
func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Author = "com.github.yeqown"
	app.Name = "micro-server-demo-cli"
	app.Usage = "micro-server-demo"

	mountCommands(app)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func mountCommands(app *cli.App) {
	app.Commands = []cli.Command{
		getAutoGenerateDBCommand(),
		getStartServerCommand(),
	}
}
