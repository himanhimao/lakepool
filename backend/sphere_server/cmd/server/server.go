package main

import (
	"github.com/himanhimao/lakepool/backend/sphere_server/internal/pkg/conf"
	pb "github.com/himanhimao/lakepool_proto/backend/proto_sphere"
	"github.com/urfave/cli"
	"os"
	"time"
	"google.golang.org/grpc"
	log "github.com/sirupsen/logrus"
	"math"
	"github.com/himanhimao/lakepool/backend/sphere_server/internal/app"
	"google.golang.org/grpc/reflection"
	"github.com/gomodule/redigo/redis"
	"github.com/himanhimao/lakepool/backend/sphere_server/internal/pkg/service"
	"github.com/himanhimao/lakepool/backend/sphere_server/internal/pkg/service/cache"
	"github.com/himanhimao/lakepool/backend/sphere_server/internal/pkg/service/coin/btc"
	"os/signal"
	"syscall"
	"net"
)

func main() {
	app := cli.NewApp()
	app.Name = "sphere-server"
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
	app.Usage = "lake pool sphere server"

	var serverConfig conf.ServerConfig
	var btcConfig conf.BTCClientConfig
	var redisConfig conf.RedisConfig
	var btcRpcClient *btc.RpcClient
	var btcCoinConfig conf.CoinConfig
	mgr := service.NewManager()
	sphereConfig := conf.NewSphereConfig()

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
			EnvVar:      "SPHERE_SERVER_HOST",
			Destination: &serverConfig.Host,
		},
		cli.IntFlag{
			Name:        "server_port, p",
			Value:       80,
			Usage:       "server port",
			EnvVar:      "SPHERE_SERVER_PORT",
			Destination: &serverConfig.Port,
		},
		cli.StringFlag{
			Name:        "server_mode, sm",
			Value:       conf.ModeDev,
			Usage:       "server mode",
			EnvVar:      "SPHERE_MODE",
			Destination: &serverConfig.Mode,
		},
		cli.StringFlag{
			Name:        "btc_host, btc_h",
			Value:       "localhost",
			Usage:       "btc server host",
			EnvVar:      "SPHERE_BTC_HOST",
			Destination: &btcConfig.Host,
		},
		cli.IntFlag{
			Name:        "btc_port, btc_p",
			Value:       8332,
			Usage:       "btc server port",
			EnvVar:      "SPHERE_BTC_PORT",
			Destination: &btcConfig.Port,
		},
		cli.StringFlag{
			Name:        "btc_username, btc_u",
			Value:       "rpc_username",
			Usage:       "btc server username(Basic Auth)",
			EnvVar:      "SPHERE_BTC_USERNAME",
			Destination: &btcConfig.User,
		},
		cli.StringFlag{
			Name:        "btc_password, btc_pass",
			Value:       "3LTBKLD9mp38sf2",
			Usage:       "btc server password(Basic Auth)",
			EnvVar:      "SPHERE_BTC_PASSWORD",
			Destination: &btcConfig.Password,
		},
		cli.BoolFlag{
			Name:        "btc_usessl, btc_us",
			Usage:       "btc server usessl",
			EnvVar:      "SPHERE_BTC_USE_SSL",
			Destination: &btcConfig.UseSSL,
		},
		cli.StringFlag{
			Name:        "redis_host, rh",
			Usage:       "redis host",
			Value:       "127.0.0.1",
			EnvVar:      "SPHERE_REDIS_HOST",
			Destination: &redisConfig.Host,
		},
		cli.IntFlag{
			Name:        "redis_port, rp",
			Usage:       "redis port",
			Value:       6379,
			EnvVar:      "SPHERE_REDIS_PORT",
			Destination: &redisConfig.Port,
		},
		cli.StringFlag{
			Name:        "redis_password, rpw",
			Usage:       "redis password",
			Value:        "",
			EnvVar:      "SPHERE_REDIS_PASSWORD",
			Destination:  &redisConfig.Password,
		},
		cli.IntFlag{
			Name:        "redis_db_num, rdn",
			Usage:       "redis db num",
			Value:       15,
			EnvVar:      "SPHERE_REDIS_DB_NUM",
			Destination: &redisConfig.DBNum,
		},
		cli.IntFlag{
			Name:        "redis_pool_max_idle",
			Usage:       "redis pool max idle",
			Value:       0,
			EnvVar:      "SPHERE_REDIS_POOL_MAX_IDLE",
			Hidden:      true,
			Destination: &redisConfig.MaxIdle,
		},
		cli.IntFlag{
			Name:        "redis_pool_max_active",
			Usage:       "redis pool max active",
			Value:       5,
			EnvVar:      "SPHERE_REDIS_POOL_MAX_ACTIVE",
			Hidden:      true,
			Destination: &redisConfig.MaxActive,
		},
		cli.IntFlag{
			Name:        "redis_pool_idle_timeout",
			Usage:       "redis pool idle timeout",
			Value:       240,
			EnvVar:      "SPHERE_REDIS_POOL_IDLE_TIMEOUT",
			Hidden:      true,
			Destination: &redisConfig.IdleTimeout,
		},
		cli.DurationFlag{
			Name:        "btc_subscribe_pull_gbt_interval, spgi",
			Usage:       "subscribe pull gbt interval(ms)",
			Value:       time.Millisecond * 500,
			EnvVar:      "SPHERE_BTC_SUBSCRIBE_PULL_GET_GBT_INTERVAL",
			Destination: &btcCoinConfig.PullGBTInterval,
		},
		cli.DurationFlag{
			Name:        "btc_subscribe_notify_interval, sni",
			Usage:       "subscribe notify interval(s)",
			Value:       time.Second * 10,
			EnvVar:      "SPHERE_BTC_SUBSCRIBE_NOTIFY_INTERVAL",
			Destination: &btcCoinConfig.NotifyInterval,
		},
		cli.DurationFlag{
			Name:        "btc_job_cache_expire_ts",
			Usage:       "job cache expire ts(s)",
			Value:       time.Second * 7200,
			EnvVar:      "SPHERE_BTC_JOB_CACHE_EXPIRE_TS",
			Destination: &btcCoinConfig.JobCacheExpireTs,
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

		log.SetOutput(os.Stdout)
		var redisPool *redis.Pool

		btcRpcClient = btc.NewClient(btcConfig.Host, btcConfig.Port, btcConfig.User,
			btcConfig.Password, btcConfig.UseSSL)

		redisPool = &redis.Pool{
			MaxIdle:     redisConfig.MaxIdle,
			MaxActive:   redisConfig.MaxActive,
			IdleTimeout: time.Duration(redisConfig.IdleTimeout) * time.Second,
			Wait:        false,
			Dial: func() (redis.Conn, error) {
				redisConn, err := redis.Dial("tcp", redisConfig.FormatAddress(),
					redis.DialDatabase(redisConfig.DBNum), redis.DialPassword(redisConfig.Password))
				if err != nil {
					return nil, err
				}
				return redisConn, err
			},
		}

		redisCache := cache.NewRedisCache().SetClient(redisPool)
		mgr.SetCacheService(redisCache)

		btcCoin := btc.NewBTCCoin().SetRPCClient(btcRpcClient)
		mgr.SetCoinService(service.CoinTypeBTC, btcCoin)
		sphereConfig.Configs[service.CoinTypeBTC] = btcCoinConfig
		return nil
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name:     "run",
			Aliases:  []string{"Run", "RUN"},
			Usage:    "Sphere Server Start",
			Category: "core",
			Before:   bootstrap,
			Action: func(c *cli.Context) {
				log.Infoln("sphere  server start,  port:", serverConfig.Port)
				lis, err := net.Listen("tcp", serverConfig.FormatHostPort())
				if err != nil {
					log.WithFields(log.Fields{
						"error": err,
					}).Fatal("failed to listen")
				}

				serverOptions := []grpc.ServerOption{
					grpc.MaxRecvMsgSize(math.MaxInt32),
					grpc.MaxSendMsgSize(math.MaxInt32),
				}

				s := grpc.NewServer(serverOptions...)
				pb.RegisterSphereServer(s, &impl.SphereServer{Conf: sphereConfig, Mgr: mgr})
				reflection.Register(s)

				signalChannel := make(chan os.Signal, 1)
				signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)

				go func(server *grpc.Server, c chan os.Signal) {
					sig := <-c
					log.Infoln(sig, "server shutting down......")
					s.Stop()
				}(s, signalChannel)

				if err := s.Serve(lis); err != nil {
					log.WithFields(log.Fields{
						"error": err,
					}).Fatalf("failed to serve")
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
