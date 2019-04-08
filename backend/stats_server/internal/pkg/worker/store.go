package worker

import (
	"github.com/influxdata/influxdb1-client/v2"
	"gopkg.in/oleiade/lane.v1"
	"sync"
	"errors"
	log "github.com/sirupsen/logrus"
	"time"
	"sync/atomic"
	"runtime"
)

var (
	ErrorUninitialized           = errors.New("Uninitialized Error.")
	ErrorInvalidConfigQueuesNum  = errors.New("Invalid Config QueuesNum.")
	ErrorInvalidConfigStoreCycle = errors.New("Invalid Config StoreCycle.")
	ErrorInvalidPackPointsSize   = errors.New("Invalid Config PackPointsSize.")
)

type Config struct {
	StoreQueuesNum int
	StoreInterval  time.Duration
	PackPointsSize int
}

type StoreQueues struct {
	quartet []*lane.Deque
}

type StoreWorker struct {
	config            *Config
	batchPointsConfig client.BatchPointsConfig
	dbClient          client.Client
	initialized       bool
	index             int32
	queues            *StoreQueues
	initOnce          sync.Once
	stopOnce          sync.Once
	loopOnce          sync.Once
	stopC             chan struct{}
}

func NewStoreQueues(num int) *StoreQueues {
	quartet := make([]*lane.Deque, num)
	i := 0
	for i < num {
		queue := lane.NewDeque()
		quartet[i] = queue
		i = i + 1
	}
	return &StoreQueues{quartet: quartet}
}

func NewDelayer(conf *Config, batchConf client.BatchPointsConfig, client client.Client) *StoreWorker {
	return &StoreWorker{config: conf, batchPointsConfig: batchConf, dbClient: client}
}

func (s *StoreQueues) Get(index int) (*lane.Deque) {
	if index > (len(s.quartet) - 1) {
		return nil
	}
	return s.quartet[index]
}

func (s *StoreQueues) Set(index int, value interface{}) {
	s.Get(index)
}

func (c *Config) validate() error {
	if c.StoreQueuesNum <= 0 {
		return ErrorInvalidConfigQueuesNum
	}

	if c.StoreInterval <= 0 {
		return ErrorInvalidConfigStoreCycle
	}

	if c.PackPointsSize <= 0 {
		return ErrorInvalidPackPointsSize
	}

	return nil
}

func (s *StoreWorker) Init() error {
	if err := s.config.validate(); err != nil {
		return err
	}
	s.initOnce.Do(func() {
		log.Infoln("database", s.batchPointsConfig.Database)
		storeQueues := NewStoreQueues(s.config.StoreQueuesNum)
		s.queues = storeQueues
		s.index = -1
		s.stopC = make(chan struct{})
		s.initialized = true
	})

	return nil
}

func (s *StoreWorker) AsyncStoreSharePoint(point *client.Point) error {
	nextIndex := s.genNextIndex()
	queue := s.queues.Get(int(nextIndex))
	queue.Append(point)
	atomic.StoreInt32(&s.index, nextIndex)
	return nil
}

func (s *StoreWorker) SyncStoreSharePoint(point *client.Point) error {
	return s.storeSinglePoint(point)
}

func (s *StoreWorker) Loop() error {
	if !s.initialized {
		return ErrorUninitialized
	}
	s.loopOnce.Do(func() {
		queueNumber := 0
		for queueNumber < s.config.StoreQueuesNum {
			go s.loopQueue(s.queues.Get(queueNumber), queueNumber)
			queueNumber = queueNumber + 1
		}
	})
	return nil
}

func (s *StoreWorker) Stop() {
	s.stopOnce.Do(func() {
		log.Infoln("delayer stop...")
		s.stopC <- struct{}{}
		s.dbClient.Close()
		s.safeStop()
	})
}

func (s *StoreWorker) loopQueue(queue *lane.Deque, number int) {
	log.WithFields(log.Fields{
		"queue_num":  number,
		"queue_size": queue.Size(),
	}).Infoln("loop queue")

	for {
		select {
		case <-s.stopC:
			return
		default:
			queueSize := queue.Size()
			if queueSize > 0 {
				var packPointsSize int
				if queue.Size() >= s.config.PackPointsSize {
					packPointsSize = s.config.PackPointsSize
				} else {
					packPointsSize = queue.Size()
				}
				var i, j int
				for j < queueSize {
					bp, _ := client.NewBatchPoints(s.batchPointsConfig)
					i = 0
					for !queue.Empty() {
						bp.AddPoint(queue.Shift().(*client.Point))
						j++
						if i == packPointsSize-1 {
							break
						}
						i++
					}

					if err := s.store(bp); err != nil {
						//TODO 记录到error文件 后续处理
						log.WithFields(log.Fields{"error": err}).Error("store batch points error ")
					} else {
						log.WithFields(log.Fields{
							"queue_number":       number,
							"pack_points_size":   packPointsSize,
							"store_points_count": j,
						}).Debugln("store batch points success.", s.batchPointsConfig.Database)
					}
				}
			} else {
				runtime.Gosched()
			}
			time.Sleep(time.Second * s.config.StoreInterval)
		}
	}
}

func (s *StoreWorker) genNextIndex() int32 {
	var nextIndex int32
	if s.index+1 > int32(s.config.StoreQueuesNum-1) {
		nextIndex = 0
	} else {
		nextIndex = s.index + 1
	}
	return nextIndex
}

func (s *StoreWorker) safeStop() {
	queueNumber := 0
	for queueNumber < s.config.StoreQueuesNum {
		queue := s.queues.Get(queueNumber)
		for !queue.Empty() {
			s.storeSinglePoint(queue.Shift().(*client.Point))
		}
		queueNumber = queueNumber + 1
	}
}

func (s *StoreWorker) store(batchPoints client.BatchPoints) error {
	if err := s.dbClient.Write(batchPoints); err != nil {
		return err
	}
	return nil
}

func (s *StoreWorker) storeSinglePoint(point *client.Point) error {
	bp, _ := client.NewBatchPoints(s.batchPointsConfig)
	bp.AddPoint(point)

	if err := s.dbClient.Write(bp); err != nil {
		return err
	}
	return nil
}
