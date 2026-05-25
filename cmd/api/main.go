package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"ebook-system/internal/handlers"
	localMiddleware "ebook-system/internal/middleware"
	"ebook-system/internal/repository"
	"ebook-system/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de sistema")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSL := os.Getenv("DB_SSLMODE")

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbSSL == "" {
		dbSSL = "disable"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPass, dbHost, dbPort, dbName, dbSSL)

	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	defer dbPool.Close()
	log.Println("Conexión exitosa a la base de datos")

	bookRepo := repository.NewBookRepository(dbPool)
	userRepo := repository.NewUserRepository(dbPool)
	purchaseRepo := repository.NewPurchaseRepository(dbPool)

	bookService := services.NewBookService(bookRepo)
	userService := services.NewUserService(userRepo)
	purchaseService := services.NewPurchaseService(purchaseRepo, userRepo, bookRepo)

	bookHandler := handlers.NewBookHandler(bookService)
	authHandler := handlers.NewAuthHandler(userService)
	userHandler := handlers.NewUserHandler(userService, purchaseService)
	purchaseHandler := handlers.NewPurchaseHandler(purchaseService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(localMiddleware.AuthMiddleware)

	r.Get("/", bookHandler.Home)
	r.Get("/catalog", bookHandler.Catalog)

	r.Get("/login", authHandler.LoginView)
	r.Post("/login", authHandler.LoginPost)
	r.Get("/register", authHandler.RegisterView)
	r.Post("/register", authHandler.RegisterPost)
	r.Post("/logout", authHandler.Logout)

	r.Group(func(r chi.Router) {
		r.Use(localMiddleware.RequireAuth)

		r.Post("/books", bookHandler.CreateBook)
		r.Get("/profile", userHandler.Profile)
		r.Post("/profile/balance", userHandler.AddBalance)
		r.Post("/books/{id}/buy", purchaseHandler.BuyBook)
		r.Get("/books/{id}/download", purchaseHandler.DownloadBook)
	})

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Servidor escuchando en http://localhost%s", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
