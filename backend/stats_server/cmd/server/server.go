package main

import (
	"github.com/himanhimao/lakepool/backend/stats_server/internal/pkg/conf"
	pb "github.com/himanhimao/lakepool_proto/backend/proto_stats"
	"github.com/urfave/cli"
	"os"
	"time"
	"google.golang.org/grpc"
	"math"
	impl "github.com/himanhimao/lakepool/backend/stats_server/internal/app"
	"github.com/influxdata/influxdb1-client/v2"
	"google.golang.org/grpc/reflection"
	log "github.com/sirupsen/logrus"
	"fmt"
	"net"
	"os/signal"
	"syscall"
	"github.com/himanhimao/lakepool/backend/stats_server/internal/pkg/worker"
)

func main() {
	app := cli.NewApp()
	app.Name = "stats-server"
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
	app.Usage = "lake pool stats server"

	var serverConfig conf.Config
	var storeWorkerConfig worker.Config
	var dbConfig client.HTTPConfig
	var dbClient client.Client
	var statsConfig conf.StatsConfig
	var batchPointsConfig client.BatchPointsConfig
	var newDelayer *worker.StoreWorker

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
			EnvVar:      "STATS_SERVER_HOST",
			Destination: &serverConfig.Host,
		},
		cli.IntFlag{
			Name:        "port, p",
			Value:       8080,
			Usage:       "server port",
			EnvVar:      "STATS_PORT",
			Destination: &serverConfig.Port,
		},
		cli.StringFlag{
			Name:        "server_mode, sm",
			Value:       conf.ModeDev,
			Usage:       "server mode",
			EnvVar:      "STATS_SERVER_MODE",
			Destination: &serverConfig.Mode,
		},
		cli.StringFlag{
			Name:        "db_addr, da",
			Value:       "http://127.0.0.1:8086",
			Usage:       "influx db http addr",
			EnvVar:      "STATS_DB_ADDR",
			Destination: &dbConfig.Addr,
		},
		cli.StringFlag{
			Name:        "db_username, du",
			Value:       "db username",
			Usage:       "influx db http username",
			EnvVar:      "STATS_DB_USERNAME",
			Destination: &dbConfig.Username,
		},
		cli.StringFlag{
			Name:        "db_password, dp",
			Value:       "db password",
			Usage:       "influx db http password",
			EnvVar:      "STATS_DB_PASSWORD",
			Destination: &dbConfig.Password,
		},
		cli.StringFlag{
			Name:        "db_database, bpd",
			Value:       "mining_stats",
			Usage:       "influx db database name",
			EnvVar:      "STATS_DB_DATABASE",
			Destination: &batchPointsConfig.Database,
		},
		cli.StringFlag{
			Name:        "db_database_measurement_precision, lp",
			Value:       "ns",
			Usage:       "influx db database measurement precision",
			EnvVar:      "STATS_DB_MEASUREMENT_PRECISION",
			Destination: &batchPointsConfig.Precision,
			Hidden:      true,
		},
		cli.StringFlag{
			Name:        "coin_type, ct",
			Value:       "BTC",
			Usage:       "coin type",
			EnvVar:      "STATS_COIN_TYPE",
			Destination: &statsConfig.CoinType,
		},
		cli.StringFlag{
			Name:        "measurement_share_prefix, msp",
			Value:       "stats_share",
			Usage:       "measurement share prefix",
			EnvVar:      "STATS_MEASUREMENT_SHARE_PREFIX",
			Destination: &statsConfig.MeasurementSharePrefix,
		},
		cli.DurationFlag{
			Name:        "store_cycle",
			Value:       1,
			Usage:       "delayer store interval(second.)",
			EnvVar:      "STATS_STORE_INTERVAL",
			Destination: &storeWorkerConfig.StoreInterval,
		},
		cli.IntFlag{
			Name:        "store queues num",
			Value:       3,
			Usage:       "delayer store queues num",
			EnvVar:      "STATS_STORE_QUEUES_NUM",
			Destination: &storeWorkerConfig.StoreQueuesNum,
		},
		cli.IntFlag{
			Name:        "pack_point_num,pps",
			Value:       10,
			Usage:       "pack point num",
			EnvVar:      "STATS_PACKET_POINT_NUM",
			Destination: &storeWorkerConfig.PackPointsSize,
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
		if dbClient, err = client.NewHTTPClient(dbConfig); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("create db client error.")
		}

		newDelayer = worker.NewDelayer(&storeWorkerConfig, batchPointsConfig, dbClient)
		if err = newDelayer.Init(); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("delayer init error.")
		}

		return nil
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:     "run",
			Aliases:  []string{"Run", "RUN"},
			Usage:    "Stats Server Start",
			Category: "core",
			Before:   bootstrap,
			Action: func(c *cli.Context) {
				//run delayer loop
				log.Infoln("delayer looping")
				if err := newDelayer.Loop(); err != nil {
					log.WithFields(log.Fields{"error": err}).Error("delayer loop error.")
				}

				log.Infoln("stats server start, port:", serverConfig.Port)
				port := fmt.Sprintf(":%d", serverConfig.Port)
				lis, err := net.Listen("tcp", port)

				if err != nil {
					log.WithFields(log.Fields{"error": err}).Fatalf("failed to listen.")
				}

				serverOptions := []grpc.ServerOption{
					grpc.MaxRecvMsgSize(math.MaxInt32),
					grpc.MaxSendMsgSize(math.MaxInt32),
				}
				s := grpc.NewServer(serverOptions...)
				statsServer := impl.NewStatsServer(newDelayer, &statsConfig)
				pb.RegisterStatsServer(s, statsServer)
				reflection.Register(s)

				//信号捕捉
				signalChannel := make(chan os.Signal, 1)
				signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)

				go func(c chan os.Signal) {
					sig := <-c
					log.Infoln(sig, "server shutting down......")
					newDelayer.Stop()
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
