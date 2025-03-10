package main

import (
	"bigyohann/apidownloader/api"
	"bigyohann/apidownloader/internal/database"
	"bigyohann/apidownloader/internal/service"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	env := os.Getenv("FOO_ENV")
	if "" == env {
		env = "development"
	}

	godotenv.Load(".env." + env + ".local")
	if "test" != env {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load() // The Original .env

	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.Warn("This is a warning")
	r := gin.Default()
	r.Use(cors.Default())
	api.HandleRouter(r)
	database.InitDatabase()

	service.CreateJwt()

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
