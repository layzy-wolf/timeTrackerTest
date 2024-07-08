package http

import (
	"context"
	"encoding/json"
	kitHttp "github.com/go-kit/kit/transport/http"
	"github.com/jmoiron/sqlx"
	"github.com/layzy-wolf/timeTrackerTest/internal/endpoints"
	"github.com/layzy-wolf/timeTrackerTest/internal/env"
	"github.com/layzy-wolf/timeTrackerTest/internal/service"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type PeopleServer struct {
	Create *kitHttp.Server
	Read   *kitHttp.Server
	Update *kitHttp.Server
	Delete *kitHttp.Server
}

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoints.PeopleCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Debugln(err)
		return nil, err
	}
	return request, nil
}

func decodeReadRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoints.PeopleReadRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Debugln(err)
		return nil, err
	}
	return request, nil
}

func decodeUpdateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoints.PeopleUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Debugln(err)
		return nil, err
	}
	return request, nil
}

func decodeDeleteRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request endpoints.PeopleDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Debugln(err)
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func MakePeopleHandler(conf *env.Config, db *sqlx.DB) *PeopleServer {
	srv := service.NewPeople(conf, db)
	en := endpoints.NewPeopleEndpoints(srv)
	return &PeopleServer{
		Create: kitHttp.NewServer(en.Create, decodeCreateRequest, encodeResponse),
		Read:   kitHttp.NewServer(en.Read, decodeReadRequest, encodeResponse),
		Update: kitHttp.NewServer(en.Update, decodeUpdateRequest, encodeResponse),
		Delete: kitHttp.NewServer(en.Delete, decodeDeleteRequest, encodeResponse),
	}
}
