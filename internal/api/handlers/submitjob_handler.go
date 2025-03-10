package api

import (
	"fmt"
	"log"
	"net/http"
)

type SubmitJobHandler struct {
	l *log.Logger
}

func NewSubmitJobHandler(l *log.Logger) *SubmitJobHandler {
	return &SubmitJobHandler{l}
}

func (s *SubmitJobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.l.Println("recvd submit job handler")

	fmt.Fprintf(w, "some data")
}
