package main

import (
	"net/http"

	"person-service/internal/config"
	"person-service/internal/handler"
	"person-service/internal/logger"
	"person-service/internal/repository"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()
	db := sqlx.MustConnect("postgres", cfg.PostgresDSN())
	defer db.Close()

	logger.InitLogger(true)

	repo := repository.NewPersonRepository(db)
	handler := handler.NewHandler(repo)

	r := mux.NewRouter()
	r.HandleFunc("/people", handler.CreatePerson).Methods("POST")
	r.HandleFunc("/people", handler.GetPeople).Methods("GET")
	r.HandleFunc("/people/{id}", handler.DeletePerson).Methods("DELETE")
	r.HandleFunc("/people/{id}", handler.UpdatePerson).Methods("PUT")

	logger.Log.Info("Starting server...", "port", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		logger.Log.Fatalw("Server failed", "error", err)
	}
}
