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

const SCRAPING_CONCURRENCY = 10
const SCRAPING_INTERVAL = 60

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

	db := database.New(conn)

	apiCfg := apiConfig{
		DB: db,
	}

	startScraping(db, SCRAPING_CONCURRENCY, SCRAPING_INTERVAL)

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
	v1Router.Get("/users", apiCfg.authMiddleware(apiCfg.handlerGetUser))

	v1Router.Post("/feeds", apiCfg.authMiddleware(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)

	v1Router.Get("/feed_follows", apiCfg.authMiddleware(apiCfg.handlerGetUserFollowedFeeds))
	v1Router.Post("/feed_follows", apiCfg.authMiddleware(apiCfg.handlerCreateFeedFollow))

	v1Router.Get("/posts", apiCfg.authMiddleware(apiCfg.handlerUserPosts))

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
