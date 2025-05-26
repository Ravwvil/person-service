package main

import (
	"log"
	"net/http"

	"person-service/internal/config"
	"person-service/internal/handler"
	"person-service/internal/logger"
	"person-service/internal/repository"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "person-service/docs" // Swagger generated docs

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// @title Person Service API
// @version 1.0
// @description API для управления информацией о людях с автоматическим обогащением данных
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.LoadConfig()

	dbURL := cfg.PostgresDSN()
	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatalf("Migration DB driver error: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	db = sqlx.MustConnect("postgres", dbURL)
	defer db.Close()

	logger.InitLogger(true)

	repo := repository.NewPersonRepository(db)
	h := handler.NewHandler(repo)

	r := mux.NewRouter()

	r.HandleFunc("/people", h.CreatePerson).Methods("POST")
	r.HandleFunc("/people", h.GetPeople).Methods("GET")
	r.HandleFunc("/people/{id}", h.UpdatePerson).Methods("PUT")
	r.HandleFunc("/people/{id}", h.DeletePerson).Methods("DELETE")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	logger.Log.Infow("Starting server", "port", cfg.Port)
	logger.Log.Infow("Swagger docs at", "url", "http://localhost:"+cfg.Port+"/swagger/index.html")

	if err := http.ListenAndServe(
		":"+cfg.Port,
		r,
	); err != nil {
		logger.Log.Fatalw("Server failed", "error", err)
	}
}
