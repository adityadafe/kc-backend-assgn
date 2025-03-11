package storage

import (
	"fmt"

	"github.com/adityadafe/kc-backend-assgn/internal/utils"
)

type JobId string

type StoreError struct {
	StoreID string `json:"store_id"`
	Error   string `json:"error"`
}

type JobInfo struct {
	Status string       `json:"status"`
	JobID  string       `json:"job_id"`
	Errors []StoreError `json:"error"`
}

type Store struct {
	storageMp map[JobId]*JobInfo
}

type Storage interface {
	AddNewJob(string)
	UpdateJob(string, string, string, string)
	GetJobStatus(string) (*JobInfo, error)
}

func CreateNewStore() Store {
	return Store{storageMp: make(map[JobId]*JobInfo)}
}

func (s *Store) AddNewJob(id string) {
	s.storageMp[JobId(id)] = &JobInfo{
		JobID:  id,
		Status: utils.JobOngoing,
		Errors: make([]StoreError, 0),
	}
}

func (s *Store) UpdateJob(jobId string, storeId string, jobStatus string, errMsg string) {

	key := JobId(jobId)
	jobInfo, exists := s.storageMp[key]
	if !exists {
		return
	}

	if jobStatus == utils.JobFailed {
		found := false
		for i, e := range jobInfo.Errors {
			if e.StoreID == storeId {
				jobInfo.Errors[i].Error = errMsg
				found = true
				break
			}
		}

		if !found {
			jobInfo.Errors = append(jobInfo.Errors, StoreError{
				StoreID: storeId,
				Error:   errMsg,
			})
		}

		jobInfo.Status = utils.JobFailed
	} else if jobInfo.Status != utils.JobFailed {
		jobInfo.Status = jobStatus
	}
}

func (s *Store) GetJobStatus(jobId string) (*JobInfo, error) {
	jobInfo, ok := s.storageMp[JobId(jobId)]
	if !ok {
		return nil, fmt.Errorf("err: %s", "job_id_does_not_exist")
	}
	return jobInfo, nil
}
