package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/layzy-wolf/timeTrackerTest/internal/env"
	log "github.com/sirupsen/logrus"
	net "net/http"
	"strings"
)

const pageLimit int = 10

type IPeople interface {
	Create(pass string) (int64, error)
	Read(filter map[string]string, page int) ([]interface{}, error)
	Update(id int, update map[string]string) (bool, error)
	Delete(id int) (bool, error)
}

type SPeople struct {
	Id             int    `json:"id" db:"id"`
	PassportSerie  string `json:"passportSerie" db:"passport_serie"`
	PassportNumber string `json:"passportNumber" db:"passport_number"`
	Name           string `json:"name" db:"name"`
	Surname        string `json:"surname" db:"surname"`
	Patronymic     string `json:"patronymic" db:"patronymic"`
	Address        string `json:"address" db:"address"`
}

type People struct {
	conf *env.Config
	db   *sqlx.DB
}

func NewPeople(conf *env.Config, db *sqlx.DB) *People {
	return &People{conf: conf, db: db}
}

// Create 		godoc
//
//	@Summary		Create new user
//	@Description	get passport serie and number and add it to db
//	@Tags			people
//	@Accept			json
//	@Produce		json
//	@Param			passportNumber	body			endpoints.PeopleCreateRequest	true		"PassportNumber"
//	@Success		200				{object}		endpoints.PeopleCreateResponse				"id, error"
//	@Failure		500				{object}		error
//	@Router /people [post]
func (p *People) Create(pass string) (int, error) {
	var r []byte
	var people SPeople
	var userId int

	arr := strings.Split(pass, " ")

	people.PassportSerie = arr[0]
	people.PassportNumber = arr[1]

	if len([]rune(people.PassportSerie)) != 4 || len([]rune(people.PassportNumber)) != 6 {
		return 0, errors.New("passport format is incorrect")
	}

	resp, err := net.Get(fmt.Sprintf("%v/info?passportSerie=%v&passportNumber=%v",
		p.conf.ExternalAPI,
		people.PassportSerie,
		people.PassportNumber),
	)

	if err == nil {
		defer resp.Body.Close()

		_, err = resp.Body.Read(r)

		if err != nil {
			log.Debugln(err)
		}

		if err = json.Unmarshal(r, &people); err != nil {
			log.Debugln(err)
		}
	} else {
		log.Debugln(err)
	}

	res, err := p.db.NamedQuery("INSERT INTO people(passport_serie, passport_number, name, surname, patronymic, address) "+
		"VALUES(:passport_serie, :passport_number, :name, :surname, :patronymic, :address) RETURNING id", people)

	if err != nil {
		log.Debugln(err)
		return 0, err
	}

	for res.Next() {
		res.Scan(&userId)
	}

	return userId, err
}

// Read 			godoc
//
//	@Summary		get users
//	@Description	returns all users, allows filters and pagination
//	@Tags			people
//	@Accept			json
//	@Produce		json
//	@Param			Filter, Page	body			endpoints.PeopleReadRequest					true	"filter, page"
//	@Success		200				{object}		endpoints.PeopleReadResponse						"users, error"
//	@Failure		500				{object}		error
//	@Router /people [get]
func (p *People) Read(filter map[string]string, page int) ([]SPeople, error) {
	var res []SPeople
	var args []interface{}
	var i int = 1

	if page == 0 {
		page = 1
	}

	query := "SELECT * FROM people WHERE id = id"

	if len(filter) > 0 {
		for key, val := range filter {
			args = append(args, val)
			query += fmt.Sprintf(" AND %v = $%v", key, i)
			i++
		}
	}

	args = append(args, page*pageLimit)
	args = append(args, page*pageLimit-pageLimit)

	query += fmt.Sprintf(" LIMIT $%v OFFSET $%v", i, i+1)

	err := p.db.Select(&res, query, args...)

	if err != nil {
		log.Debugln(err)
		return nil, err
	}

	return res, err
}

// Update 			godoc
//
//	@Summary		update user
//	@Description	update user by his id
//	@Tags			people
//	@Accept			json
//	@Produce		json
//	@Param			Id, Update 			body		endpoints.PeopleUpdateRequest			true		"id, update"
//	@Success		200				{object}		endpoints.BoolResponse							"ok, error"
//	@Failure		500				{object}		error
//	@Router /people [put]
func (p *People) Update(id int, update map[string]string) (bool, error) {
	var args = make(map[string]interface{})
	var u string = "id = id"

	for key, val := range update {
		args[key] = val
		u += fmt.Sprintf(", %v = :%v", key, key)
	}

	if _, err := p.db.NamedExec(fmt.Sprintf("UPDATE people SET %v WHERE id = %d", u, id), args); err != nil {
		log.Debugln(err)
		return false, err
	}
	return true, nil
}

// Delete 			godoc
//
//	@Summary		delete user
//	@Description	delete user by his id
//	@Tags			people
//	@Accept			json
//	@Produce		json
//	@Param			Id				body			endpoints.PeopleUpdateRequest				true		"id"
//	@Success		200				{object}		endpoints.BoolResponse									"ok, error"
//	@Failure		500				{object}		error
//	@Router /people [delete]
func (p *People) Delete(id int) (bool, error) {
	_, err := p.db.Query("DELETE FROM people WHERE id = $1", id)
	if err != nil {
		log.Debugln(err)
		return false, err
	}
	return true, nil
}
