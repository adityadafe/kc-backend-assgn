package api

import (
	"fmt"
	"log"
	"net/http"
)

type GetJobInfoHandler struct {
	l *log.Logger
}

func NewGetJobInfoHandler(l *log.Logger) *GetJobInfoHandler {
	return &GetJobInfoHandler{l}
}

func (g *GetJobInfoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	jobID := r.URL.Query().Get("jobid")

	g.l.Println("recvd get job id handler:", jobID)

	fmt.Fprintf(w, "some data from get job info")
}
