package main

import (
	"fmt"
	_ "github.com/davyxu/cellnet/codec/json"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/conf"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/app/impl"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/app/server"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/service"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/util"
	"github.com/urfave/cli"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/service/sphere"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/service/stats"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/service/user"
	slog "github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/service/log"
)

var (
	DefaultDifficultyFilePath string = fmt.Sprintf("%s%s%s", util.GetCurrPath(), string(os.PathSeparator), "work/conf/default_difficulty.yaml")
)

func main() {
	app := cli.NewApp()
	app.Name = "stratum-server"
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
	app.Usage = "lake pool stratum server"

	var serverConfig conf.ServerConfig
	var stratumConfig conf.StratumConfig
	var sphereGRPCConfig, statsGRPCConfig, logGRPCConfig conf.GRPCConfig
	ser := server.NewServer()
	router := server.NewRouter(0)
	sMgr := service.NewManager()
	config := conf.NewConfig()
	config.ServerConfig = &serverConfig
	config.StratumConfig = &stratumConfig
	impl.RegisterRoute(router)
	ser.SetConfig(config).SetServiceMgr(sMgr)

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
			EnvVar:      "STRATUM_SERVER_HOST",
			Destination: &serverConfig.Host,
		},
		cli.IntFlag{
			Name:        "server_port, sp",
			Value:       18801,
			Usage:       "server port",
			EnvVar:      "STRATUM_SERVER_PORT",
			Destination: &serverConfig.Port,
		},
		cli.StringFlag{
			Name:        "server_mode, sm",
			Value:       conf.ModeDev,
			Usage:       "server mode",
			EnvVar:      "STRATUM_SERVER_MODE",
			Destination: &serverConfig.Mode,
		},
		cli.Uint64Flag{
			Name:        "min_difficulty, mind",
			Value:       1,
			Usage:       "min difficulty",
			EnvVar:      "STRATUM_MIN_DIFFICULTY",
			Destination: &stratumConfig.MinDifficulty,
		},
		cli.Uint64Flag{
			Name:        "max_difficulty, maxd",
			Value:       4611686018427387904,
			Usage:       "max difficulty",
			EnvVar:      "STRATUM_MAX_DIFFICULTY",
			Destination: &stratumConfig.MaxDifficulty,
		},
		cli.Float64Flag{
			Name:        "difficulty_factor, df",
			Value:       10,
			Usage:       "stratum difficulty factor",
			EnvVar:      "STRATUM_DIFFICULTY_FACTOR",
			Destination: &stratumConfig.DifficultyFactor,
		},
		cli.Uint64Flag{
			Name:        "default_difficulty, dd",
			Value:       16384,
			Usage:       "default difficulty",
			EnvVar:      "STRATUM_DEFAULT_DIFFICULTY",
			Destination: &stratumConfig.DefaultDifficulty,
		},
		cli.StringFlag{
			Name:        "default_difficulty_file_path, ddfp",
			Value:       DefaultDifficultyFilePath,
			Usage:       "default difficulty file path",
			EnvVar:      "STRATUM_DEFAULT_DIFFICULTY_FILE_PATH",
			Destination: &stratumConfig.DefaultDifficultyFilePath,
		},
		cli.StringFlag{
			Name:        "subscribe_notify_placeholder, snp",
			Value:       "ae6812eb4cd7735a302a8a9dd95cf71f",
			Usage:       "subscribe notify placeholder",
			EnvVar:      "STRATUM_SUBSCRIBE_NOTIFY_PLACEHOLDER",
			Destination: &stratumConfig.NotifyPlaceholder,
		},
		cli.StringFlag{
			Name:        "subscribe_difficulty_placeholder, sdp",
			Value:       "1K3kCJFD2PYF99t2eBSQ3stmdA1jrp7Nyf",
			Usage:       "stratum subscribe difficulty placeholder",
			EnvVar:      "STRATUM_SUBSCRIBE_DIFFICULTY_PLACEHOLDER",
			Destination: &stratumConfig.DifficultyPlaceholder,
		},
		cli.IntFlag{
			Name:        "subscribe_extra_nonce1_length, sen1l",
			Value:       8,
			Usage:       "subscribe extra nonce1 length",
			EnvVar:      "STRATUM_SUBSCRIBE_EXTRA_NONCE1_LENGTH",
			Destination: &stratumConfig.ExtraNonce1Length,
		},
		cli.IntFlag{
			Name:        "subscribe_extra_nonce2_length, sen2l",
			Value:       4,
			Usage:       "subscribe extra nocne2 length",
			EnvVar:      "STRATUM_SUBSCRIBE_EXTRA_NONCE2_LENGTH",
			Destination: &stratumConfig.ExtraNonce2Length,
		},
		cli.BoolFlag{
			Name:        "authorize_skip_verify, sakv",
			EnvVar:      "STRATUM_AUTHORIZE_SKIP_VERIFY",
			Usage:       "authorize skip verify",
			Destination: &stratumConfig.SkipVerify,
		},
		cli.StringFlag{
			Name:        "authorize_name_separator, ans",
			EnvVar:      "STRATUM_AUTHORIZE_NAME_SEPARATOR",
			Value:       ".",
			Usage:       "authorize_name_separator",
			Destination: &stratumConfig.NameSeparator,
		},
		cli.StringFlag{
			Name:        "notify_coin_type, nct",
			EnvVar:      "STRATUM_NOTIFY_COIN_TYPE",
			Value:       "BTC",
			Usage:       "stratum notify coin type",
			Destination: &stratumConfig.CoinType,
		},
		cli.StringFlag{
			Name:        "notify_coin_tag, nctag",
			EnvVar:      "STRATUM_NOTIFY_COIN_TAG",
			Value:       "/lakepool/",
			Usage:       "stratum notify pool flag",
			Destination: &stratumConfig.PoolTag,
		},
		cli.StringFlag{
			Name:        "notify_payout_address, npa",
			EnvVar:      "NOTIFY_PAYOUT_ADDRESS",
			Value:       "1K3kCJFD2PYF99t2eBSQ3stmdA1jrp7Nyf",
			Usage:       "notify payout address",
			Destination: &stratumConfig.PayoutAddress,
		},
		cli.DurationFlag{
			Name:        "notify_loop_interval, nli",
			EnvVar:      "NOTIFY_LOOP_INTERVAL",
			Value:       25,
			Usage:       "notify loop interval (second)",
			Destination: &stratumConfig.NotifyLoopInterval,
		},
		cli.BoolFlag{
			Name:        "is_used_test_net, iutn",
			EnvVar:      "STRATUM_IS_USED_TEST_NET",
			Usage:       "is used test net",
			Destination: &stratumConfig.UsedTestNet,
		},
		cli.DurationFlag{
			Name:        "difficulty_check_loop_interval, dcli",
			EnvVar:      "STRATUM_DIFFICULTY_CHECK_LOOP_INTERVAL",
			Value:       60,
			Usage:       "difficulty check loop interval (second)",
			Destination: &stratumConfig.DifficultyCheckLoopInterval,
		},
		cli.DurationFlag{
			Name:        "job_hash_check_loop_interval, jhcli",
			EnvVar:      "STRATUM_JOB_HASH_CHECK_LOOP_INTERVAL",
			Value:       1000,
			Usage:       "job hash check loop interval (Millisecond)",
			Destination: &stratumConfig.JobHashCheckLoopInterval,
		},
		cli.StringFlag{
			Name:        "sphere_grpc_server_host",
			EnvVar:      "STRATUM_SPHERE_GRPC_SERVER_HOST",
			Value:       "localhost",
			Usage:       "sphere_grpc_server_host",
			Destination: &sphereGRPCConfig.Host,
		},
		cli.IntFlag{
			Name:        "sphere_grpc_server_port",
			EnvVar:      "STRATUM_SPHERE_GRPC_SERVER_PORT",
			Value:       80,
			Usage:       "sphere grpc server port",
			Destination: &sphereGRPCConfig.Port,
		},
		cli.StringFlag{
			Name:        "stats_grpc_server_host",
			EnvVar:      "STRATUM_STATS_GRPC_SERVER_HOST",
			Value:       "localhost",
			Usage:       "stats_grpc_server_host",
			Destination: &statsGRPCConfig.Host,
		},
		cli.IntFlag{
			Name:        "stats_grpc_server_port",
			EnvVar:      "STRATUM_STATS_GRPC_SERVER_PORT",
			Value:       8080,
			Usage:       "stats_grpc_server_port",
			Destination: &statsGRPCConfig.Port,
		},
		cli.StringFlag{
			Name:        "log_grpc_server_host",
			EnvVar:      "STRATUM_LOG_GRPC_SERVER_HOST",
			Value:       "localhost",
			Usage:       "log_grpc_server_host",
			Destination: &logGRPCConfig.Host,
		},
		cli.IntFlag{
			Name:        "log_grpc_server_port",
			EnvVar:      "STRATUM_LOG_GRPC_SERVER_PORT",
			Value:       8081,
			Usage:       "log_grpc_server_port",
			Destination: &logGRPCConfig.Port,
		},
	}

	bootstrap := func(c *cli.Context) error {
		if err := config.Validate(); err != nil {
			return err
		}
		//bind sphere service
		sphereService := sphere.NewGRPCService()
		sphereService.SetConfig(&sphereGRPCConfig)
		sMgr.SetSphereService(sphereService)

		//bind stratum service
		statsService := stats.NewGRPCService()
		statsService.SetConfig(&statsGRPCConfig)
		sMgr.SetStatsService(statsService)

		//bind log service
		logService := slog.NewGRPCService()
		logService.SetConfig(&logGRPCConfig)
		sMgr.SetLogService(logService)

		if stratumConfig.SkipVerify {
			userService := user.NewSkipUserService()
			sMgr.SetUserService(userService)
		}

		//init log mode
		if serverConfig.Mode == conf.ModeProd {
			log.SetFormatter(&log.TextFormatter{})
			log.SetLevel(log.InfoLevel)
		} else {
			log.SetFormatter(&log.JSONFormatter{})
			log.SetLevel(log.DebugLevel)
		}

		log.SetOutput(os.Stdout)
		return nil
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:     "run",
			Aliases:  []string{"Run", "RUN"},
			Usage:    "Stratum Server Start",
			Category: "core",
			Before:   bootstrap,
			Action: func(c *cli.Context) {
				if err := ser.Init(); err != nil {
					log.Fatal("Server initialization failed: ", err)
				}
				ser.Run(router)
				return
			},
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
