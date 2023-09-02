package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/quintui/rssagg/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("PORT Environment variable is required")
	}

	databaseUrl := os.Getenv("DB_URL")

	if databaseUrl == "" {
		log.Fatal("DB_URL Environment variable is required")
	}

	conn, err := sql.Open("postgres", databaseUrl)

	if err != nil {
		log.Fatal("Something went wrong during connection to database: ", err.Error())
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlerReadiness)
	router.Mount("/v1", v1Router)

	v1Router.Post("/users", apiCfg.handlerCreateUser)

	srv := &http.Server{
		Addr:    ":" + portString,
		Handler: router,
	}

	log.Printf("Server is running on the PORT: %v", portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Something went wrong during server startup: ", err.Error())
	}

}
