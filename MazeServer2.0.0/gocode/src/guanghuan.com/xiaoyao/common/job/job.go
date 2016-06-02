package xyjob

import (
	"math/rand"
	"sync"
	"time"
)

type JobHandler func(job IJob) (IResult, error)

var (
	next_job_id  int32
	job_id_mutex sync.Mutex
)

func Init() {
	rand.Seed(time.Now().UnixNano())
	next_job_id = int32(rand.Intn(5) * 10000)
	if next_job_id == 0 {
		next_job_id = 1
	}
}

func NextJobId() (id int32) {
	job_id_mutex.Lock()
	defer job_id_mutex.Unlock()

	id = next_job_id
	next_job_id++

	return
}

// interfaces
type IJob interface {
	ExecuteWithHandler(handler JobHandler) (result IResult, err error)
	Execute() (IResult, error)
	JobId() int32
}

type IResult interface {
	IsSuccess() bool
	FailReason() int32
	Error() string
	JobId() int32
}

// predefined error code
const (
	FAIL    = -1
	SUCCESS = 0
)

// default job implementation
type Job struct {
	job_id int32
}

func NewJob() (job *Job) {
	job = &Job{
		job_id: NextJobId(),
	}
	return
}

func (job *Job) JobId() int32 {
	return job.job_id
}

func (job *Job) Execute() (result IResult, err error) {
	return
}

func (job *Job) ExecuteWithHandler(handler JobHandler) (result IResult, err error) {
	result, err = handler(job)
	return
}

// default result implementation
type Result struct {
	job_id      int32
	fail_reason int32
	error_msg   string
}

func NewResult(job_id int32, fail_reason int32, err_msg string) (result *Result) {
	result = &Result{
		job_id:      job_id,
		fail_reason: fail_reason,
		error_msg:   err_msg,
	}
	return
}

func (result *Result) IsSuccess() bool {
	return (result.fail_reason != 0)
}

func (result *Result) Error() string {
	return result.error_msg
}

func (result *Result) FailReason() int32 {
	return result.fail_reason
}

func (result *Result) JobId() int32 {
	return result.job_id
}
