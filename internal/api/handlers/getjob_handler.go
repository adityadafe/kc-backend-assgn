package api

import (
	"log"
	"net/http"

	"github.com/adityadafe/kc-backend-assgn/internal/models"
	"github.com/adityadafe/kc-backend-assgn/internal/storage"
	"github.com/adityadafe/kc-backend-assgn/internal/utils"
)

type GetJobInfoHandler struct {
	l  *log.Logger
	db storage.Storage
}

func NewGetJobInfoHandler(l *log.Logger, db storage.Storage) *GetJobInfoHandler {
	return &GetJobInfoHandler{l, db}
}

func (g *GetJobInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	jobID := r.URL.Query().Get("jobid")
	g.l.Println("recvd get job id:", jobID)

	jobStatus, err := g.db.GetJobStatus(jobID)

	//TODO: make these into functions
	if err != nil {
		jobFailed := new(models.GetJobResponseBodyFailed)
		jobFailed.JobId = jobID
		jobFailed.Status = storage.JobFailed
		jobFailed.Error.Error = "store_does_not_exist"
		utils.WriteJson(w, http.StatusBadRequest, jobFailed)
		return
	}

	if jobStatus == storage.JobFailed {
		jobFailed := new(models.GetJobResponseBodyFailed)
		jobFailed.JobId = jobID
		jobFailed.Status = storage.JobFailed
		utils.WriteJson(w, http.StatusBadRequest, jobFailed)
		return
	}

	jobOngoingOrCompleted := new(models.GetJobResponseBodyForCompletedOrOngoing)
	jobOngoingOrCompleted.JobId = jobID
	utils.WriteJson(w, http.StatusOK, jobOngoingOrCompleted)
}
