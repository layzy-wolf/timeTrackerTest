package endpoints

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/layzy-wolf/timeTrackerTest/internal/service"
	log "github.com/sirupsen/logrus"
)

type TaskEndpoints struct {
	GetEndpoint    endpoint.Endpoint
	BeginEndpoint  endpoint.Endpoint
	FinishEndpoint endpoint.Endpoint
}

func NewTaskEndpoints(srv *service.Task) *TaskEndpoints {
	return &TaskEndpoints{
		GetEndpoint:    MakeGetEndpoint(srv),
		BeginEndpoint:  MakeBeginEndpoint(srv),
		FinishEndpoint: MakeFinishEndpoint(srv),
	}
}

func MakeGetEndpoint(srv *service.Task) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var tasks []service.STask
		req := request.(TaskGetRequest)

		res, err := srv.Get(req.People, req.PeriodStart, req.PeriodEnd)

		if err != nil {
			log.Debugln(err)
		}

		for _, val := range res {
			tasks = append(tasks, val)
		}

		return TaskGetResponse{
			Tasks: tasks,
			Error: err,
		}, err
	}
}

func MakeBeginEndpoint(srv *service.Task) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(TaskBeginRequest)

		res, err := srv.Begin(req.People, req.Task)

		if err != nil {
			log.Debugln(err)
		}

		return TaskBeginResponse{
			Task:  res,
			Error: err,
		}, err
	}
}

func MakeFinishEndpoint(srv *service.Task) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(TaskFinishRequest)

		success, err := srv.Finish(req.Task)

		if err != nil {
			log.Debugln(err)
		}

		return BoolResponse{
			Success: success,
			Error:   err,
		}, err
	}
}
