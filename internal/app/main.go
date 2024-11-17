package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/hmcalister/GoDockerComposeTemplate/internal/app/database"
	"github.com/jackc/pgx/v5"
)

func getDatabase() (*database.Queries, error) {
	ctx := context.Background()

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPasswordFile := os.Getenv("DB_PASSWORD_FILE")
	dbPassword, err := os.ReadFile(dbPasswordFile)
	if err != nil {
		slog.Error("error when reading password file", "error", err)
		return nil, err
	}

	conn, err := pgx.Connect(ctx, fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", dbUser, dbPassword, dbHost, dbPort, dbName))
	if err != nil {
		slog.Error("error when connecting to db", "error", err)
		return nil, err
	}
	return database.New(conn), nil
}

func main() {
	debugFlag := os.Getenv("DEBUG")
	if debugFlag == "" {
		logFile, err := os.OpenFile("/app/logs/log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			slog.Error("could not create log file", "error", err)
		}
		defer logFile.Close()
		defaultLoggerHandler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		slog.SetDefault(slog.New(defaultLoggerHandler))
	} else {
		defaultLoggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		slog.SetDefault(slog.New(defaultLoggerHandler))
	}

	db, err := getDatabase()
	if err != nil {
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Get("/newAuthor/{authorName}", func(w http.ResponseWriter, r *http.Request) {
		authorName := chi.URLParam(r, "authorName")
		slog.Info("new author request", "authorName", authorName)

		ctx := context.Background()
		_, err := db.CreateAuthor(ctx, database.CreateAuthorParams{
			Name: authorName,
		})
		if err != nil {
			slog.Error("error when creating new author", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	router.Get("/allAuthors", func(w http.ResponseWriter, r *http.Request) {

		slog.Info("all author request")

		ctx := context.Background()
		allAuthors, err := db.ListAuthors(ctx)
		if err != nil {
			slog.Error("error when requesting authors", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		for _, author := range allAuthors {
			w.Write([]byte(fmt.Sprintf("%v: %v\n", author.ID, author.Name)))
		}
	})
	slog.Info("ready to serve")
	if err := http.ListenAndServe(":8080", router); err != nil {
		slog.Error("error during listen and serve", "error", err)
	}
}
