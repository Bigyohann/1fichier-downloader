package main

import (
	"bigyohann/apidownloader/api"
	"bigyohann/apidownloader/internal/database"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	r := gin.Default()
	api.HandleRouter(r)
	database.InitDatabase()

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
