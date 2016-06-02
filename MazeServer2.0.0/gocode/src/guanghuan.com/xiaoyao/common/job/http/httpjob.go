package xyhttpjob

import (
	"bytes"
	xyjob "guanghuan.com/xiaoyao/common/job"
	xylog "guanghuan.com/xiaoyao/common/log"
	xyhttpservice "guanghuan.com/xiaoyao/common/service/http"
	"io/ioutil"
	"net/http"
)

const (
	HTTP_FAIL            = -1
	HTTP_SUCCESS         = 0
	HTTP_OP_FAIL         = 1
	HTTP_RESP_DATA_ERROR = 2
)

type HttpJob struct {
	xyjob.Job
	op      xyhttpservice.HttpOp
	Url     string
	ReqData []byte
}

func NewHttpJob(op xyhttpservice.HttpOp, url string, req_data []byte) (job *HttpJob) {
	job = &HttpJob{
		Job:     *xyjob.NewJob(),
		op:      op,
		Url:     url,
		ReqData: req_data,
	}
	return job
}

type HttpResult struct {
	xyjob.Result
	StatusCode int
	Url        string
	RespData   []byte
}

func NewHttpResult(job *HttpJob, resp_data []byte, status_code int, fail_reason int32, error_msg string) (result *HttpResult) {
	result = &HttpResult{
		Result:     *xyjob.NewResult(job.JobId(), fail_reason, error_msg),
		Url:        job.Url,
		RespData:   resp_data,
		StatusCode: status_code,
	}
	return
}

func (result *HttpResult) IsSuccess() bool {
	return (result.FailReason() != HTTP_SUCCESS)
}

func (result *HttpResult) IsStatusOK() bool {
	return (result.StatusCode == http.StatusOK)
}
func (job *HttpJob) Post() (result xyjob.IResult, err error) {
	var (
		resp        *http.Response
		err_msg     string
		resp_data   []byte
		fail_reason int32 = HTTP_SUCCESS
		status_code int
	)

	resp, err = http.Post(job.Url, "", bytes.NewReader(job.ReqData))

	if err != nil {
		fail_reason = HTTP_OP_FAIL
		err_msg = err.Error()
		xylog.ErrorNoId("Post Error: %s", err.Error())
		goto ErrorHandle
	}

	resp_data, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fail_reason = HTTP_RESP_DATA_ERROR
		err_msg = err.Error()
		xylog.ErrorNoId("Response Error: %s", err.Error())
		goto ErrorHandle
	}
	status_code = resp.StatusCode

ErrorHandle:
	result = NewHttpResult(job, resp_data, status_code, fail_reason, err_msg)

	return
}

func (job *HttpJob) Get() (result xyjob.IResult, err error) {
	result = NewHttpResult(job, nil, http.StatusOK, HTTP_SUCCESS, "")

	return
}

func (job *HttpJob) Execute() (result xyjob.IResult, err error) {
	switch job.op {
	case xyhttpservice.HttpPost:
		result, err = job.Post()
	case xyhttpservice.HttpGet:
		result, err = job.Get()
	default:
		// not support
		xylog.Warning("Op not support: %d", job.op)
	}
	return
}
