package main

import (
	"github.com/himanhimao/lakepool/backend/log_server/internal/pkg/conf"
	impl "github.com/himanhimao/lakepool/backend/log_server/internal/app"
	pb "github.com/himanhimao/lakepool_proto/backend/proto_log"
	"github.com/urfave/cli"
	"os"
	"time"
	"google.golang.org/grpc"
	"math"
	"github.com/influxdata/influxdb1-client/v2"
	log "github.com/sirupsen/logrus"
	"fmt"
	"net"
	"os/signal"
	"syscall"
	"google.golang.org/grpc/reflection"
)

func main() {
	app := cli.NewApp()
	app.Name = "log-server"
	app.Version = "0.1.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "himanhimao",
			Email: "himanhimmao@github.com",
		},
	}
	app.Copyright = "(c) 2019 himanhimao"
	app.HelpName = "contrive"
	app.Usage = "lake pool log server"

	var serverConfig conf.Config
	var logConfig conf.LogConfig
	var dbConfig client.HTTPConfig
	var influxClient client.Client
	var blockBPConfig client.BatchPointsConfig
	var shareBPConfig client.BatchPointsConfig

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "server_name, sn",
			Value:       app.Name,
			Hidden:      true,
			Destination: &serverConfig.Name,
		},
		cli.StringFlag{
			Name:        "server_host, sh",
			Value:       "0.0.0.0",
			Usage:       "server host",
			EnvVar:      "LOG_SERVER_HOST",
			Destination: &serverConfig.Host,
		},
		cli.IntFlag{
			Name:        "port, p",
			Value:       8082,
			Usage:       "server port",
			EnvVar:      "LOG_SERVER_PORT",
			Destination: &serverConfig.Port,
		},
		cli.StringFlag{
			Name:        "server_mode, sm",
			Value:       conf.ModeDev,
			Usage:       "server mode",
			EnvVar:      "LOG_MODE",
			Destination: &serverConfig.Mode,
		},
		cli.StringFlag{
			Name:        "db_addr, da",
			Value:       "http://127.0.0.1:8086",
			Usage:       "influx db http addr",
			EnvVar:      "LOG_DB_ADDR",
			Destination: &dbConfig.Addr,
		},
		cli.StringFlag{
			Name:        "db_username, du",
			Value:       "",
			Usage:       "influx db http username",
			EnvVar:      "LOG_DB_USERNAME",
			Destination: &dbConfig.Username,
		},
		cli.StringFlag{
			Name:        "db_password, dp",
			Value:       "",
			Usage:       "influx db http password",
			EnvVar:      "LOG_DB_PASSWORD",
			Destination: &dbConfig.Password,
		},
		cli.StringFlag{
			Name:        "block_database, bd",
			Value:       "mining_block",
			Usage:       "influx block batch points database name",
			EnvVar:      "LOG_BlOCK_DATABASE",
			Destination: &blockBPConfig.Database,
		},
		cli.StringFlag{
			Name:        "block_precision, bdp",
			Value:       "ns",
			Usage:       "influx db storage block database precision",
			EnvVar:      "LOG_BLOCK_PRECISION",
			Destination: &blockBPConfig.Precision,
			Hidden:      true,
		},
		cli.StringFlag{
			Name:        "share_database, sd",
			Value:       "mining_share",
			Usage:       "influx db storage share database name",
			EnvVar:      "LOG_SHARE_DATABASE",
			Destination: &shareBPConfig.Database,
		},
		cli.StringFlag{
			Name:        "share_precision, sdp",
			Value:       "m",
			Usage:       "influx db storage block database precision",
			EnvVar:      "LOG_SHARE_PRECISION",
			Destination: &shareBPConfig.Precision,
			Hidden:      true,
		},
		cli.StringFlag{
			Name:        "measurement_share_prefix, msp",
			Value:       "log_share",
			Usage:       "measurement share",
			EnvVar:      "LOG_MEASUREMENT_SHARE_PREFIX",
			Destination: &logConfig.MeasurementSharePrefix,
		},
		cli.StringFlag{
			Name:        "measurement_block_prefix, mbp",
			Value:       "log_block",
			Usage:       "measurement block",
			EnvVar:      "LOG_MEASUREMENT_BLOCK_PREFIX",
			Destination: &logConfig.MeasurementBlockPrefix,
		},
	}

	bootstrap := func(c *cli.Context) error {
		// init log mode
		if serverConfig.Mode == conf.ModeProd {
			log.SetFormatter(&log.TextFormatter{})
			log.SetLevel(log.InfoLevel)
		} else {
			log.SetFormatter(&log.JSONFormatter{})
			log.SetLevel(log.DebugLevel)
		}

		var err error
		if influxClient, err = client.NewHTTPClient(dbConfig); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("create db client error.")
		}
		return nil
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:     "run",
			Aliases:  []string{"Run", "RUN"},
			Usage:    "Log Server Start",
			Category: "core",
			Before:   bootstrap,
			Action: func(c *cli.Context) {
				port := fmt.Sprintf(":%d", serverConfig.Port)
				lis, err := net.Listen("tcp", port)
				log.Infoln("log server start, port:", serverConfig.Port)

				if err != nil {
					log.WithFields(log.Fields{"error": err}).Fatalf("failed to listen.")
				}

				serverOptions := []grpc.ServerOption{
					grpc.MaxRecvMsgSize(math.MaxInt32),
					grpc.MaxSendMsgSize(math.MaxInt32),
				}
				s := grpc.NewServer(serverOptions...)
				logServer := impl.NewLogServer(influxClient, &logConfig, blockBPConfig, shareBPConfig)
				pb.RegisterLogServer(s, logServer)
				reflection.Register(s)

				//Signal capture
				signalChannel := make(chan os.Signal, 1)
				signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)

				go func(c chan os.Signal) {
					sig := <-c
					log.Infoln(sig, "server shutting down......")
					s.Stop()
				}(signalChannel)

				if err := s.Serve(lis); err != nil {
					log.WithFields(log.Fields{
						"error": err,
					}).Fatalf("failed to serve.")
				}
				return
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
