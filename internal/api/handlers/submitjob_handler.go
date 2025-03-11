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

func (s *SubmitJobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.l.Println("Recvd a job")
	newJobPayload := new(models.JobPayload)

	err := json.NewDecoder(r.Body).Decode(newJobPayload)

	if err != nil {
		s.l.Println("failed to decode")
		utils.WriteJson(w, http.StatusInternalServerError, `{error:"Interal Server error"}`)
		return
	}

	if newJobPayload.Count != len(newJobPayload.Visits) {
		s.l.Println("count != len(visits)")
		utils.WriteJson(w, http.StatusBadRequest, `{error:"Bad request"}`)
		return
	}

	for _, eachJob := range newJobPayload.Visits {
		if len(eachJob.ImageUrls) < 1 || eachJob.StoreId == "" || eachJob.VisitTime == "" {
			utils.WriteJson(w, http.StatusBadRequest, `{error:"Bad request"}`)
			return
		}
	}

	newSubmitJobResponseBody := new(models.SubmitJobResponseBody)
	newSubmitJobResponseBody.JobId = uuid.NewString()

	s.db.AddNewJob(newSubmitJobResponseBody.JobId)

	go process.ProcessJob(newSubmitJobResponseBody.JobId, *newJobPayload, s.db, s.l)

	utils.WriteJson(w, http.StatusCreated, newSubmitJobResponseBody)
}
