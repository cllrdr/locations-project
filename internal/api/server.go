package api

import (
	"log"
	//"time"
	//"net/http"
	"github.com/gin-gonic/gin"
  	"github.com/sirupsen/logrus"
  	"locations-project/internal/app/handler"
  	"locations-project/internal/app/repository"
)

//var lines = []string{"first line", "second line", "third line", "fourth line"}


func StartServer() {
  log.Println("Starting server")

  repo, err := repository.NewRepository()
  if err != nil {
    logrus.Error("ошибка инициализации репозитория")
  }

  handler := handler.NewHandler(repo)

  r := gin.Default()
  // добавляем наш html/шаблон
  r.LoadHTMLGlob("templates/*")
  r.Static("/static", "./resources")

  r.GET("/all-locations", handler.GetLocations)
  r.GET("/location/:id", handler.GetLocation)
  r.GET("/fav-locations", handler.GetFavorites)

  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
  log.Println("Server down")

}
