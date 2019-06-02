package impl

import (
	"github.com/davyxu/cellnet"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/app/server"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/context"
	"time"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/cellnet/proto"
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/service"
	log "github.com/sirupsen/logrus"
)

func JobNotify(s *server.Server, ev cellnet.Event) {
	sid := ev.Session().ID()
	contextSet := ev.Session().Peer().(cellnet.ContextSet)
	stratumContext, ok := contextSet.GetContext(sid)
	stratumContext.(*context.StratumContext).NotifyMutex.Lock()


	if !ok {
		log.WithFields(log.Fields{"sid": sid}).Errorln("job notify check not found context")
		stratumContext.(*context.StratumContext).NotifyMutex.Unlock()
		ev.Session().Close()
		return
	}

	workerName := stratumContext.(*context.StratumContext).WorkerName
	notifyLoopInterval := s.GetConfig().StratumConfig.NotifyLoopInterval
	notifyLoopIntervalTs := int(notifyLoopInterval)
	latestNotifyTs := stratumContext.(*context.StratumContext).LatestNotifyTs
	diffTs := int(time.Now().Unix() - latestNotifyTs)

	// 如果时间不满足, sleep
	if diffTs-notifyLoopIntervalTs < 0 {
		sleepTs := notifyLoopIntervalTs - diffTs
		time.Sleep(time.Duration(sleepTs) * time.Second)
	}

	latestNotifyJobHeight := stratumContext.(*context.StratumContext).LatestNotifyJobHeight
	latestNotifyJobIndex := stratumContext.(*context.StratumContext).LatestNotifyJobIndex
	latestDifficulty := stratumContext.(*context.StratumContext).LatestDifficulty

	newNotifyTs := time.Now().Unix()
	newNotifyJobIndex := latestNotifyJobIndex + 1
	newJob := s.GetJobRepo().GetJob(latestNotifyJobHeight, newNotifyJobIndex)
	if newJob == nil {
		//如果未找到任务, return
		log.WithFields(log.Fields{
			"sid":         sid,
			"worker_name": workerName,
		}).Errorln("jobNotify height not found new job")
		stratumContext.(*context.StratumContext).NotifyMutex.Unlock()
		return
	}
	jobId := service.GenerateJobId(latestNotifyJobHeight, latestNotifyJobIndex, latestDifficulty)
	newJob = newJob.Fill(jobId, latestNotifyTs, false)
	stratumContext.(*context.StratumContext).LatestNotifyJobHeight = latestNotifyJobHeight
	stratumContext.(*context.StratumContext).LatestNotifyTs =  newNotifyTs
	stratumContext.(*context.StratumContext).LatestNotifyJobIndex = newNotifyJobIndex
	stratumContext.(*context.StratumContext).NotifyCount ++

	notifyRESP := proto.NewJSONRpcNotifyRESP(newJob.ToJSONInterface())
	ev.Session().Send(notifyRESP)

	log.WithFields(log.Fields{
		"sid":       sid,
		"height":    newJob.Meta.Height,
		"job_index": newNotifyJobIndex,
	}).Infoln("Sent new job")
	stratumContext.(*context.StratumContext).NotifyMutex.Unlock()
}

func JobDifficultyCheck(s *server.Server, ev cellnet.Event) {
	sid := ev.Session().ID()
	contextSet := ev.Session().Peer().(cellnet.ContextSet)
	stratumContext, ok := contextSet.GetContext(sid)
	if ! ok {
		log.WithFields(log.Fields{"sid": sid}).Errorln(sid, "difficulty check not found context")
		ev.Session().Close()
		return
	}

	slideWindow := stratumContext.(*context.StratumContext).SlideWindow
	workerName := stratumContext.(*context.StratumContext).WorkerName
	minDifficulty := s.GetConfig().StratumConfig.MinDifficulty
	maxDifficulty := s.GetConfig().StratumConfig.MaxDifficulty
	difficulty := stratumContext.(*context.StratumContext).LatestDifficulty
	difficultyCheckInterval := s.GetConfig().StratumConfig.DifficultyCheckLoopInterval

	factor := 1 * s.GetConfig().StratumConfig.DifficultyFactor
	if factor <= 0 {
		factor = float64(1)
	}
	exceptCount := float64(float64(difficultyCheckInterval) / factor)
	averageCount := slideWindow.Average(difficultyCheckInterval * time.Second)
	log.Debugln(sid, workerName, "exceptCount", exceptCount, "averageCount", averageCount)

	var curDifficulty uint64
	if exceptCount > averageCount {
		if difficulty == minDifficulty {
			log.WithFields(log.Fields{
				"sid":         sid,
				"worker_name": workerName,
			}).Debugln("Difficulty is equal to the minimum difficulty")
			return
		}
		if averageCount == 0 {
			averageCount = 1
		}
		rate := exceptCount / averageCount
		curDifficulty = uint64(float64(difficulty) / rate)

		if curDifficulty < minDifficulty {
			curDifficulty = minDifficulty
		}
	} else {
		if difficulty == maxDifficulty {
			log.WithFields(log.Fields{
				"sid":         sid,
				"worker_name": workerName,
			}).Debugln("Difficulty is equal to the maximum difficulty")
			return
		}
		rate := averageCount / exceptCount
		curDifficulty = uint64(float64(difficulty) * rate)
		if curDifficulty > maxDifficulty {
			curDifficulty = maxDifficulty
		}
	}

	if curDifficulty == difficulty {
		log.WithFields(log.Fields{
			"sid":         sid,
			"worker_name": workerName,
		}).Debugln("Difficulty does not need to change")
		return
	}

	stratumContext.(*context.StratumContext).LatestDifficulty = curDifficulty
	difficultyResp := proto.NewJSONRpcSetDifficultyRESP(curDifficulty)
	ev.Session().Send(difficultyResp)
}

func JobHeightCheck(s *server.Server, ev cellnet.Event) {
	sid := ev.Session().ID()

	contextSet := ev.Session().Peer().(cellnet.ContextSet)
	stratumContext, ok := contextSet.GetContext(sid)
	stratumContext.(*context.StratumContext).NotifyMutex.Lock()

	if !ok {
		log.WithFields(log.Fields{
			"sid": sid,
		}).Errorln("job height check not found context")
		stratumContext.(*context.StratumContext).NotifyMutex.Unlock()
		ev.Session().Close()
		return
	}

	workerName := stratumContext.(*context.StratumContext).WorkerName
	latestDifficulty := stratumContext.(*context.StratumContext).LatestDifficulty
	currentNotifyJobHeight := stratumContext.(*context.StratumContext).LatestNotifyJobHeight
	latestJobHeight := s.GetJobRepo().GetLatestHeight()
	if currentNotifyJobHeight != latestJobHeight {
		latestJobIndex := 0
		latestNotifyTs := time.Now().Unix()
		newJob := s.GetJobRepo().GetJob(latestJobHeight, latestJobIndex)

		if newJob == nil {
			log.WithFields(log.Fields{
				"sid":         sid,
				"worker_name": workerName,
			}).Errorln("job height check not found new job")
			stratumContext.(*context.StratumContext).NotifyMutex.Unlock()
			ev.Session().Close()
			return
		}

		jobId := service.GenerateJobId(latestJobHeight, latestJobIndex, latestDifficulty)
		newJob = newJob.Fill(jobId, latestNotifyTs, true)
		stratumContext.(*context.StratumContext).LatestNotifyJobHeight = latestJobHeight
		stratumContext.(*context.StratumContext).LatestNotifyTs = latestNotifyTs
		stratumContext.(*context.StratumContext).LatestNotifyJobIndex = latestJobIndex
		stratumContext.(*context.StratumContext).NotifyCount++

		notifyRESP := proto.NewJSONRpcNotifyRESP(newJob.ToJSONInterface())
		ev.Session().Send(notifyRESP)

		log.WithFields(log.Fields{
			"sid":         sid,
			"worker_name": workerName,
			"height":      newJob.Meta.Height,
		}).Infoln("Sent a new height job")
	}

	stratumContext.(*context.StratumContext).NotifyMutex.Unlock()
}
