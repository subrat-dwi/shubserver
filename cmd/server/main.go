package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/subrat-dwi/shubserver/internal/app"
	"github.com/subrat-dwi/shubserver/internal/db"
)

func main() {
	_ = godotenv.Load()

	if os.Getenv("DATABASE_URL") == "" {
		env := os.Getenv("APP_ENV")
		if env == "docker" {
			if v := os.Getenv("DATABASE_URL_DOCKER"); v != "" {
				_ = os.Setenv("DATABASE_URL", v)
			}
		} else {
			if v := os.Getenv("DATABASE_URL_LOCAL"); v != "" {
				_ = os.Setenv("DATABASE_URL", v)
			}
		}
	}

	dbPool := db.ConnectDB()
	defer dbPool.Close()

	s := app.Setup(dbPool)

	fmt.Printf("Listening to http://localhost%v\n", s.Addr)
	err := http.ListenAndServe(s.Addr, s.Router)

	if err != nil {
		fmt.Println("Server stopped:", err)
	}
}
