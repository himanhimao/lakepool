package server

import (
	"github.com/himanhimao/lakepool/backend/stratum_server/internal/pkg/service"
	"sync"
	"time"
	log "github.com/sirupsen/logrus"
	"context"
)

const (
	InitStratumJobListLength = 10
	CleanInterval            = time.Minute * 1
)

type JobList struct {
	index int
	len   int
	list  []*service.StratumJob
}

type JobRepo struct {
	cleanHeight  int32
	latestHeight int32
	jobs       map[int32]*JobList
	mutex        sync.Mutex
}

func NewJobRepo() *JobRepo {
	jobs := make(map[int32]*JobList, 0)
	return &JobRepo{jobs: jobs}
}

func newJobList() *JobList {
	list := make([]*service.StratumJob, InitStratumJobListLength)
	return &JobList{0, len(list), list}
}

func (l *JobList) Append(job *service.StratumJob) {
	if l.index <= l.len-1 {
		l.list[l.index] = job
	} else {
		l.list = append(l.list, job)
		l.len = len(l.list)
	}
	l.index++
}

func (l *JobList) Get(index int) *service.StratumJob {
	if index > l.len {
		return nil
	}
	return l.list[index]
}

func (r *JobRepo) GetLatestHeight() int32 {
	return r.latestHeight
}

func (r *JobRepo) SetJob(height int32, job *service.StratumJob) {
	if height > 0 {
		r.mutex.Lock()
		var jobList *JobList
		if jobList = r.getJobList(height); jobList == nil {
			jobList = newJobList()
			r.jobs[height] = jobList
		}
		jobList.Append(job)

		if r.cleanHeight == 0 {
			r.cleanHeight = height
		}
		r.latestHeight = height
		r.mutex.Unlock()
	}
}

func (r *JobRepo) getJobList(height int32) *JobList {
	var jobList *JobList
	var ok bool
	if jobList, ok = r.jobs[height]; !ok {
		return nil
	}
	return jobList
}

func (r *JobRepo) GetJob(height int32, index int) *service.StratumJob {
	var jobList *JobList
	if jobList = r.getJobList(height); jobList == nil {
		return nil
	}
	return jobList.Get(index)
}

func (r *JobRepo) Clean(ctx context.Context) {
	cleanLoopTicker := time.NewTicker(CleanInterval)

	for {
		select {
		case <-ctx.Done():
			log.Debugln("job cleanup routine cancel")
			return
		case t := <-cleanLoopTicker.C:
			log.Debugln("clean job", t)
			if r.cleanHeight > 0 && r.cleanHeight < r.latestHeight {
				r.mutex.Lock()
				jobList := r.getJobList(r.cleanHeight)
				if jobList != nil {
					r.SetJob(r.cleanHeight, nil)
					delete(r.jobs, r.cleanHeight)
				}
				r.cleanHeight++
				r.mutex.Unlock()
			}
		}
	}

}
