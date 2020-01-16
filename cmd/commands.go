package main

import (
	"log"

	"github.com/yeqown/micro-server-demo/global"
	"github.com/yeqown/micro-server-demo/model"
	"github.com/yeqown/micro-server-demo/router"

	"github.com/yeqown/infrastructure/framework/gormic"

	"github.com/jinzhu/gorm"
	"github.com/urfave/cli"
)

func connectDB() *gorm.DB {
	dbConn, err := gormic.ConnectSqlite3(global.GetConfig().Sqlite3)
	if err != nil {
		panic(err)
	}
	return dbConn
}

func getAutoGenerateDBCommand() cli.Command {
	dbConn := connectDB()

	return cli.Command{
		Name:  "generate",
		Usage: "auto generate DB tables",
		Action: func(c *cli.Context) error {
			dbConn.AutoMigrate(
				&model.FooModel{},
				// more...
			)
			return nil
		},
	}
}

func getStartServerCommand() cli.Command {
	return cli.Command{
		Name:  "start",
		Usage: "starting REST and GRPC servers",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "rpcPort",
				Value: 8080,
				Usage: "rpc server port",
			},
			cli.IntFlag{
				Name:  "httpPort",
				Value: 8081,
				Usage: "http server port",
			},
		},
		Action: func(c *cli.Context) error {
			httpPort := c.Int("httpPort")
			go func() {
				_httpSrv := router.NewHTTP(httpPort)
				if err := _httpSrv.Run(); err != nil {
					log.Fatal(err)
				}
			}()

			rpcPort := c.Int("rpcPort")
			_grpcSrv := router.NewgRPC(rpcPort)
			log.Printf("running gRPC server on: %d\n", rpcPort)
			return _grpcSrv.Run()
		},
	}
}
