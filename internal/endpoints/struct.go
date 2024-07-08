package endpoints

import (
	"github.com/layzy-wolf/timeTrackerTest/internal/service"
	"time"
)

type PeopleCreateRequest struct {
	PassportNumber string `json:"passportNumber"`
}

type PeopleCreateResponse struct {
	People int   `json:"people"`
	Error  error `json:"error"`
}

type PeopleReadRequest struct {
	Filter map[string]string `json:"filter"`
	Page   int               `json:"page"`
}

type PeopleReadResponse struct {
	People []service.SPeople `json:"people"`
	Error  error             `json:"error"`
}

type PeopleUpdateRequest struct {
	People int               `json:"people"`
	Update map[string]string `json:"update"`
}

type PeopleDeleteRequest struct {
	People int `json:"people"`
}

type TaskGetRequest struct {
	People      int       `json:"people"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
}

type TaskGetResponse struct {
	Tasks []service.STask `json:"tasks"`
	Error error           `json:"error"`
}

type TaskBeginRequest struct {
	People int    `json:"people"`
	Task   string `json:"task"`
}

type TaskBeginResponse struct {
	Task  int   `json:"task"`
	Error error `json:"error"`
}

type TaskFinishRequest struct {
	Task int `json:"task"`
}

type BoolResponse struct {
	Success bool  `json:"success"`
	Error   error `json:"error"`
}
