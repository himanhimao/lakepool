package main

import (
	"github.com/himanhimao/lakepool/backend/stats_server/internal/pkg/conf"
	pb "github.com/himanhimao/lakepool_proto/backend/proto_log"
	"github.com/urfave/cli"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
	"github.com/influxdata/influxdb1-client/v2"
	"os/signal"
	"syscall"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	impl "github.com/himanhimao/lakepool/backend/stats_server/internal/app"
)

func main() {
	app := cli.NewApp()
	app.Name = "stats-pusher"
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
	app.Usage = "lake pool stats pusher"

	var dbConfig client.HTTPConfig
	var dbClient client.Client
	var pusherConfig conf.PusherConfig
	var logServerGRPCConfig conf.GRPCConfig
	var logGRPCClient pb.LogClient
	var pushWorker *impl.PushWorker

	//dbConfig.InsecureSkipVerify = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "client_mode, sm",
			Value:       conf.ModeDev,
			Usage:       "client mode",
			EnvVar:      "STATS_PUSHER_CLIENT_MODE",
			Destination: &pusherConfig.Mode,
		},
		cli.StringFlag{
			Name:        "coin_type, cp",
			Value:       "BTC",
			Usage:       "coin type",
			EnvVar:      "STATS_COIN_TYPE",
			Destination: &pusherConfig.CoinType,
		},
		cli.StringFlag{
			Name:        "db_addr, da",
			Value:       "http://127.0.0.1:8086",
			Usage:       "influx db http addr",
			EnvVar:      "STATS_PUSHER_DB_ADDR",
			Destination: &dbConfig.Addr,
		},
		cli.StringFlag{
			Name:        "db_username, du",
			Value:       "",
			Usage:       "influx db http username",
			EnvVar:      "STATS_PUSHER_DB_USERNAME",
			Destination: &dbConfig.Username,
		},
		cli.StringFlag{
			Name:        "db_password, dp",
			Value:       "",
			Usage:       "influx db http password",
			EnvVar:      "STATS_PUSHER_DB_PASSWORD",
			Destination: &dbConfig.Password,
		},
		cli.StringFlag{
			Name:        "db_database, dd",
			Value:       "mining_stats",
			Usage:       "influx db database name",
			EnvVar:      "STATS_PUSHER_DB_DATABASE",
			Destination: &pusherConfig.Database,
		},
		cli.StringFlag{
			Name:        "db_measurement_prefix, dmp",
			Value:       "stats_share",
			Usage:       "influx db measurement prefix",
			EnvVar:      "STATS_PUSHER_DB_PREFIX",
			Destination: &pusherConfig.MeasurementPrefix,
			Hidden:      true,
		},
		cli.StringFlag{
			Name:        "db_measurement_suffix, dmn",
			Value:       "1min",
			Usage:       "influx db measurement name",
			EnvVar:      "STATS_PUSHER_DB_SUFFIX",
			Destination: &pusherConfig.MeasurementSuffix,
			Hidden:      true,
		},
		cli.StringFlag{
			Name:        "db_measurement_precision",
			Value:       "m",
			Usage:       "influx db measurement precision",
			EnvVar:      "STATS_PUSHER_DB_PRECISION",
			Destination: &pusherConfig.Precision,
			Hidden:      true,
		},
		cli.StringFlag{
			Name:        "db_measurement_retention_strategy, dbrs",
			Value:       "two_weeks",
			Usage:       "influx db measurement retention strategy",
			EnvVar:      "STATS_PUSHER_MEASUREMENT_RETENTION_STRATEGY",
			Destination: &pusherConfig.RetentionStrategy,
			Hidden:      true,
		},
		cli.DurationFlag{
			Name:        "read_interval, ri",
			Value:       5,
			Usage:       "read interval(second.)",
			EnvVar:      "STATS_PUSHER_READ_INTERVAL",
			Destination: &pusherConfig.ReadInterval,
		},
		cli.IntFlag{
			Name:        "per_size, ps",
			Value:       200,
			Usage:       "push per size",
			EnvVar:      "STATS_PUSHER_PER_SIZE",
			Destination: &pusherConfig.PerSize,
		},
		cli.StringFlag{
			Name:        "log_server_grpc_host, lsgh",
			Value:       "localhost",
			Usage:       "log server grpc host",
			EnvVar:      "STATS_PUSHER_LOG_SERVER_GRPC_HOST",
			Destination: &logServerGRPCConfig.Host,
		},
		cli.IntFlag{
			Name:        "log_server_grpc_port,  lsgp",
			Value:       8082,
			Usage:       "log server grpc port",
			EnvVar:      "STATS_PUSHER_LOG_SERVER_GRPC_PORT",
			Destination: &logServerGRPCConfig.Port,
		},
	}

	bootstrap := func(c *cli.Context) error {
		// init log mode
		if pusherConfig.Mode == conf.ModeProd {
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

		address := logServerGRPCConfig.FormatHostPort()
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
		if err != nil {
			return err
		}
		logGRPCClient = pb.NewLogClient(conn)
		pushWorker = impl.NewPushWorker(logGRPCClient, dbClient, &pusherConfig)
		return nil
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:     "run",
			Aliases:  []string{"Run", "RUN"},
			Usage:    "Stats pusher start",
			Category: "core",
			Before:   bootstrap,
			Action: func(c *cli.Context) {
				//信号捕捉
				signalChannel := make(chan os.Signal, 1)
				signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
				go func(c chan os.Signal) {
					sig := <-c
					pushWorker.Stop()
					log.Infoln(sig, "server shutting down......")
				}(signalChannel)
				pushWorker.Run()
				return
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
