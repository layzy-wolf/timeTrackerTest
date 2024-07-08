package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/layzy-wolf/timeTrackerTest/internal/service"
	log "github.com/sirupsen/logrus"
)

type PeopleEndpoints struct {
	Create endpoint.Endpoint
	Read   endpoint.Endpoint
	Update endpoint.Endpoint
	Delete endpoint.Endpoint
}

func NewPeopleEndpoints(srv *service.People) *PeopleEndpoints {
	return &PeopleEndpoints{
		Create: MakeCreateEndpoint(srv),
		Read:   MakeReadEndpoint(srv),
		Update: MakeUpdateEndpoint(srv),
		Delete: MakeDeleteEndpoint(srv),
	}
}

func MakeCreateEndpoint(srv *service.People) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PeopleCreateRequest)
		res, err := srv.Create(req.PassportNumber)

		if err != nil {
			log.Debugln(err)
		}

		return PeopleCreateResponse{
			People: res,
			Error:  err,
		}, err
	}
}

func MakeReadEndpoint(srv *service.People) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var people []service.SPeople
		req := request.(PeopleReadRequest)
		res, err := srv.Read(req.Filter, req.Page)

		if err != nil {
			log.Debugln(err)
		}

		for _, val := range res {
			people = append(people, val)
		}

		return PeopleReadResponse{
			People: people,
			Error:  err,
		}, err
	}
}

func MakeUpdateEndpoint(srv *service.People) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PeopleUpdateRequest)

		success, err := srv.Update(req.People, req.Update)

		if err != nil {
			log.Debugln(err)
		}

		return BoolResponse{
			Success: success,
			Error:   err,
		}, err
	}
}

func MakeDeleteEndpoint(srv *service.People) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PeopleDeleteRequest)

		success, err := srv.Delete(req.People)

		if err != nil {
			log.Debugln(err)
		}

		return BoolResponse{
			Success: success,
			Error:   err,
		}, err
	}
}
