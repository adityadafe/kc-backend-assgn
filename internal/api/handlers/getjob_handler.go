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

// Handler for "/status"
// GetJobHandler godoc
//
// @Summary Get job status
// @Description Retrieve current status of a job
// @Produce  json
// @Param   jobid query string true "Job ID (example: job-123)"
//
// @Success 200 {object} models.GetJobResponseBodyForCompletedOrOngoing
// @Failure 400 {object} models.GetJobResponseBodyFailed
// @Router /status [get]
func (g *GetJobInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	jobID := r.URL.Query().Get("jobid")
	g.l.Println("Request for job id:", jobID)

	jobInfo, err := g.db.GetJobStatus(jobID)
	if err != nil {
		utils.WriteJson(w, http.StatusBadRequest, `{}`)
		return
	}

	if jobInfo.Status == utils.JobFailed {
		failedErrors := make([]models.FailedJobError, 0, len(jobInfo.Errors))
		for _, e := range jobInfo.Errors {
			failedErrors = append(failedErrors, models.FailedJobError{
				StoreId: e.StoreID,
			})
		}

		utils.WriteJson(w, http.StatusOK, models.GetJobResponseBodyFailed{
			Status: jobInfo.Status,
			JobId:  jobID,
			Error:  failedErrors,
		})
		return
	}

	utils.WriteJson(w, http.StatusOK, models.GetJobResponseBodyForCompletedOrOngoing{
		Status: jobInfo.Status,
		JobId:  jobID,
	})
}
