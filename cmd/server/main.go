package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/neelRajesh95/contacts-api/internal/database"
	"github.com/neelRajesh95/contacts-api/internal/handler"
	"github.com/neelRajesh95/contacts-api/internal/repository"
	"github.com/neelRajesh95/contacts-api/internal/service"
)

func main() {

	databaseURL := os.Getenv(
		"DATABASE_URL",
	)

	if databaseURL == "" {
		log.Fatal(
			"DATABASE_URL environment variable is required",
		)
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	db, err := database.NewPostgresPool(
		databaseURL,
	)

	if err != nil {
		log.Fatalf(
			"database connection failed: %v",
			err,
		)
	}

	defer db.Close()

	contactRepository :=
		repository.NewContactRepository(db)

	contactService :=
		service.NewContactService(
			contactRepository,
		)

	contactHandler :=
		handler.NewContactHandler(
			contactService,
		)

	router := chi.NewRouter()

	router.Use(
		middleware.RequestID,
	)

	router.Use(
		middleware.RealIP,
	)

	router.Use(
		middleware.Logger,
	)

	router.Use(
		middleware.Recoverer,
	)

	router.Use(
		middleware.Timeout(
			10 * time.Second,
		),
	)

	router.Get(
		"/health",
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			w.Header().Set(
				"Content-Type",
				"application/json",
			)

			w.WriteHeader(
				http.StatusOK,
			)

			_, _ = w.Write(
				[]byte(
					`{"status":"ok"}`,
				),
			)
		},
	)

	router.Post(
		"/contacts",
		contactHandler.CreateContact,
	)

	router.Get(
		"/contacts",
		contactHandler.GetContacts,
	)

	router.Post(
		"/contacts/{id}/enrich",
		contactHandler.EnrichContact,
	)

	server := &http.Server{
		Addr: ":" + port,

		Handler: router,

		ReadHeaderTimeout:
			5 * time.Second,

		ReadTimeout:
			10 * time.Second,

		WriteTimeout:
			10 * time.Second,

		IdleTimeout:
			60 * time.Second,
	}

	go func() {

		log.Printf(
			"server running on port %s",
			port,
		)

		if err := server.ListenAndServe();
		err != nil &&
			err != http.ErrServerClosed {

			log.Fatalf(
				"server failed: %v",
				err,
			)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-stop

	log.Println(
		"shutting down server",
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf(
			"server shutdown failed: %v",
			err,
		)
	}

	log.Println(
		"server stopped",
	)
}