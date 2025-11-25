package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/hmcalister/Go-Compose-Template/cmd/initialization"
	"github.com/hmcalister/Go-Compose-Template/internal/database"
)

func main() {
	if err := initialization.SetupLogger(); err != nil {
		os.Exit(1)
	}

	db, err := initialization.GetDatabase()
	if err != nil {
		slog.Error("error while initializing database", "error", err)
		os.Exit(1)
	}

	// --------------------------------------------------------------------------------

	mux := http.NewServeMux()

	mux.HandleFunc("POST /authors/{authorName}", func(w http.ResponseWriter, r *http.Request) {
		authorName := r.PathValue("authorName")
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
		w.WriteHeader(http.StatusCreated)
	})

	mux.HandleFunc("GET /authors", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("all authors request")

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
	if err := http.ListenAndServe(":8080", mux); err != nil {
		slog.Error("error during listen and serve", "error", err)
	}
}
