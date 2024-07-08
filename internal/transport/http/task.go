package http

import (
	"context"
	"encoding/json"
	kitHttp "github.com/go-kit/kit/transport/http"
	"github.com/jmoiron/sqlx"
	"github.com/layzy-wolf/timeTrackerTest/internal/endpoints"
	"github.com/layzy-wolf/timeTrackerTest/internal/service"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type TaskServer struct {
	Get    *kitHttp.Server
	Begin  *kitHttp.Server
	Finish *kitHttp.Server
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoints.TaskGetRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Debugln(err)
		return nil, err
	}
	return request, nil
}

func decodeBeginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoints.TaskBeginRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Debugln(err)
		return nil, err
	}
	return request, nil
}

func decodeFinishRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoints.TaskFinishRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Debugln(err)
		return nil, err
	}
	return request, nil
}

func MakeTaskHandler(db *sqlx.DB) *TaskServer {
	srv := service.NewTask(db)
	en := endpoints.NewTaskEndpoints(srv)
	return &TaskServer{
		Get:    kitHttp.NewServer(en.GetEndpoint, decodeGetRequest, encodeResponse),
		Begin:  kitHttp.NewServer(en.BeginEndpoint, decodeBeginRequest, encodeResponse),
		Finish: kitHttp.NewServer(en.FinishEndpoint, decodeFinishRequest, encodeResponse),
	}
}
