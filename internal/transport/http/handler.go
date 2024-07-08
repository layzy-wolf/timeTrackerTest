package http

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/layzy-wolf/timeTrackerTest/docs"
	"github.com/layzy-wolf/timeTrackerTest/internal/env"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Тестовое задание Effective Mobile
//	@version		1.0
//	@description	Сервис тайм-трекер

//	@contact.name	Заула Никита
//	@contact.url	https://hh.ru/resume/2859ad12ff0d463d060039ed1f7a4352466859
//	@contact.email	nikita.zatula@mail.ru

func Handler(conf *env.Config, db *sqlx.DB) *gin.Engine {
	r := gin.New()

	if !conf.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	p := MakePeopleHandler(conf, db)
	t := MakeTaskHandler(db)

	peopleGr := r.Group("/people")
	{
		peopleGr.POST("", gin.WrapH(p.Create))
		peopleGr.GET("", gin.WrapH(p.Read))
		peopleGr.PATCH("", gin.WrapH(p.Update))
		peopleGr.PUT("", gin.WrapH(p.Update))
		peopleGr.DELETE("", gin.WrapH(p.Delete))
	}

	taskGr := r.Group("/task")
	{
		taskGr.GET("", gin.WrapH(t.Get))
		taskGr.POST("/start", gin.WrapH(t.Begin))
		taskGr.POST("/end", gin.WrapH(t.Finish))
	}

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
