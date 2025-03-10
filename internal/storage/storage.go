package storage

import "fmt"

type JobId string
type JobStatus string

const (
	JobOngoing   string = "ongoing"
	JobCompleted string = "completed"
	JobFailed    string = "failed"

	DbDoesNotExist string = "does_not_exist"
)

type Store struct {
	storageMp map[JobId]JobStatus
}

type Storage interface {
	AddNewJob(string)
	UpdateJob(string, string)
	GetJobStatus(string) (string, error)
}

func CreateNewStore() Store {
	return Store{storageMp: make(map[JobId]JobStatus)}
}

func (s *Store) AddNewJob(id string) {
	s.storageMp[JobId(id)] = JobStatus(JobOngoing)
}

func (s *Store) UpdateJob(id string, status string) {
	s.storageMp[JobId(id)] = JobStatus(status)
}

func (s *Store) GetJobStatus(id string) (string, error) {
	val, ok := (s.storageMp[JobId(id)])

	if !ok {
		return "", fmt.Errorf("err: %s", DbDoesNotExist)
	}

	return string(val), nil
}
