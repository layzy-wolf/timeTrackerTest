package service

import (
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"time"
)

type ITask interface {
	Get(people int, periodStart time.Time, periodEnd time.Time) ([]interface{}, error)
	Begin(people int, task string) (int64, error)
	Finish(task int) (bool, error)
}

type STask struct {
	Id           int       `json:"id" db:"id"`
	PeopleId     int       `json:"people_id" db:"people_id"`
	Name         string    `json:"task" db:"task"`
	TaskStart    time.Time `json:"task_start" db:"task_start"`
	TaskEnd      time.Time `json:"task_end" db:"task_end"`
	TaskInterval string    `json:"task_interval" db:"task_interval"`
}

type Task struct {
	db *sqlx.DB
}

func NewTask(db *sqlx.DB) *Task {
	return &Task{db: db}
}

// Get 				godoc
//
//	@Summary		get tasks
//	@Description	get tasks between dates and orders by their interval
//	@Tags			task
//	@Accept			json
//	@Produce		json
//	@Param			Params			body			endpoints.TaskGetRequest		true		"get request"
//	@Success		200				{object}		endpoints.TaskGetResponse					"Tasks, error"
//	@Failure		500				{object}		error
//	@Router /task [get]
func (t *Task) Get(id int, periodStart time.Time, periodEnd time.Time) ([]STask, error) {
	var res []STask

	err := t.db.Select(&res, "SELECT * FROM task WHERE people_id = $1 AND task_start >= $2 AND task_end >= $3 ORDER BY task_interval", id, periodStart, periodEnd)

	if err != nil {
		log.Debugln(err)
	}

	return res, err
}

// Begin 			godoc
//
//	@Summary		Begin task
//	@Description	begin task by user id
//	@Tags			task
//	@Accept			json
//	@Produce		json
//	@Param			Json			body			endpoints.TaskBeginRequest		true		"user_id, task title"
//	@Success		200				{object}		endpoints.TaskBeginResponse					"id, error"
//	@Failure		500				{object}		error
//	@Router /task/start [post]
func (t *Task) Begin(id int, task string) (int, error) {
	var taskId int

	res, err := t.db.Queryx("INSERT INTO task (people_id, task, task_start) VALUES ($1, $2, $3) RETURNING id", id, task, time.Now().UTC())
	if err != nil {
		log.Debugln(err)
		return 0, err
	}

	for res.Next() {
		res.Scan(&taskId)
	}

	return taskId, nil
}

// Finish 			godoc
//
//	@Summary		Finish task
//	@Description	finish task by task id
//	@Tags			task
//	@Accept			json
//	@Produce		json
//	@Param			Json			body			endpoints.TaskFinishRequest		true		"task_id"
//	@Success		200				{object}		endpoints.BoolResponse						"ok, error"
//	@Failure		500				{object}		error
//	@Router /task/end [post]
func (t *Task) Finish(task int) (bool, error) {
	if _, err := t.db.Exec("UPDATE task SET task_end = $1, task_interval = (SELECT ($1 - task.task_start) AS interval FROM task WHERE id = $2) WHERE id = $2",
		time.Now().UTC(),
		task,
	); err != nil {
		log.Debugln(err)
		return false, err
	}
	return true, nil
}
