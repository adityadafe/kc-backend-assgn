package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adityadafe/kc-backend-assgn/internal/models"
	"github.com/adityadafe/kc-backend-assgn/internal/process"
	"github.com/adityadafe/kc-backend-assgn/internal/storage"
	"github.com/adityadafe/kc-backend-assgn/internal/utils"
	"github.com/google/uuid"
)

type SubmitJobHandler struct {
	l  *log.Logger
	db storage.Storage
}

func NewSubmitJobHandler(l *log.Logger, db storage.Storage) *SubmitJobHandler {
	return &SubmitJobHandler{l, db}
}

// Handler for "/submit"
// SubmitJobHandler godoc
//
// @Summary Submit job for image processing
// @Description Add a new job to processing queue
// @Accept   json
// @Produce  json
// @Param   job body models.JobPayload true "Job details"
//
// @Success 201 {object} models.SubmitJobResponseBody
// @Failure 400 {object} models.SubmitJobFailedResponseBody
// @Router /submit [post]
func (s *SubmitJobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.l.Println("Recvd a job")
	newJobPayload := new(models.JobPayload)

	err := json.NewDecoder(r.Body).Decode(newJobPayload)

	if err != nil {
		s.l.Println("failed to decode")
		utils.WriteJson(w, http.StatusInternalServerError, models.SubmitJobFailedResponseBody{Error: "Internal server Error"})
		return
	}

	if newJobPayload.Count != len(newJobPayload.Visits) {
		s.l.Println("count != len(visits)")
		utils.WriteJson(w, http.StatusBadRequest, models.SubmitJobFailedResponseBody{Error: "Bad Request"})
		return
	}

	for _, eachJob := range newJobPayload.Visits {
		if len(eachJob.ImageUrls) < 1 || eachJob.StoreId == "" || eachJob.VisitTime == "" {
			utils.WriteJson(w, http.StatusBadRequest, models.SubmitJobFailedResponseBody{Error: "Bad Request"})
			return
		}
	}

	newSubmitJobResponseBody := new(models.SubmitJobResponseBody)
	newSubmitJobResponseBody.JobId = uuid.NewString()

	s.db.AddNewJob(newSubmitJobResponseBody.JobId)

	go process.ProcessJob(newSubmitJobResponseBody.JobId, *newJobPayload, s.db, s.l)

	utils.WriteJson(w, http.StatusCreated, newSubmitJobResponseBody)
}
