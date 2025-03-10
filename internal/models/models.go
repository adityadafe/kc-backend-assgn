package models

type StoreVisit struct {
	StoreId   string   `json:"store_id"`
	ImageUrls []string `json:"image_url"`
	VisitTime string   `json:"visit_time"`
}

type JobPayload struct {
	Count  int          `json:"count"`
	Visits []StoreVisit `json:"visits"`
}

type SubmitJobResponseBody struct {
	JobId string `json:"job_id"`
}

type GetJobResponseBodyForCompletedOrOngoing struct {
	Status string `json:"status"`
	JobId  string `json:"job_id"`
}

type FailedJobError struct {
	StoreId string `json:"store_id"`
	Error   string `json:"error"`
}

type GetJobResponseBodyFailed struct {
	Status string         `json:"status"`
	JobId  string         `json:"job_id"`
	Error  FailedJobError `json:"error"`
}
