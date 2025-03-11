package storage

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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
	CheckStore(string) error
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

func (s *Store) CheckStore(storeID string) error {
	const csvPath = "./store_manager.csv"

	absPath, err := filepath.Abs(csvPath)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path for %q: %w", csvPath, err)
	}

	if _, err := os.Stat(absPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("store data file not found at %q", absPath)
		}
		return fmt.Errorf("error accessing file %q: %w", absPath, err)
	}

	file, err := os.Open(absPath)
	if err != nil {
		return fmt.Errorf("failed to open store data file %q: %w", absPath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("invalid CSV header in %q: %w", absPath, err)
	}

	for recordNumber := 1; ; recordNumber++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading record #%d in %q: %w",
				recordNumber, absPath, err)
		}

		if len(record) < 3 {
			continue
		}

		if strings.TrimSpace(record[2]) == storeID {
			return nil
		}
	}

	return fmt.Errorf("store ID %q not found in %q", storeID, absPath)
}
