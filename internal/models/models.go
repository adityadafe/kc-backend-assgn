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

type Result struct {
	StoreID   string `json:"store_id"`
	ImageURL  string `json:"image_url"`
	Perimeter int    `json:"perimeter"`
	Error     string `json:"error"`
}

type GetJobResponseBodyFailed struct {
	Status string           `json:"status"`
	JobId  string           `json:"job_id"`
	Error  []FailedJobError `json:"error"`
}

type FailedJobError struct {
	StoreId string `json:"store_id"`
	Error   string `json:"error"`
}

type GetJobResponseBodyForCompletedOrOngoing struct {
	Status string `json:"status"`
	JobId  string `json:"job_id"`
}
