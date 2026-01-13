package main

import (
"context"
"log"
"net/http"
"os"
"os/signal"
"syscall"
"time"

"recruitflow/apps/api/internal/auth"
"recruitflow/apps/api/internal/config"
"recruitflow/apps/api/internal/db"
"recruitflow/apps/api/internal/httpapi"

"github.com/go-chi/chi/v5"
"github.com/go-chi/cors"
)

func main() {
cfg := config.Load()

pool, err := db.NewPool(cfg.DBURL)
if err != nil {
log.Fatal(err)
}
defer pool.Close()

jwtSvc := auth.NewJWT(cfg.JWTSecret)

r := chi.NewRouter()
r.Use(cors.Handler(cors.Options{
AllowedOrigins:   []string{
      "https://YOUR-PROJECT.vercel.app","http://localhost:5173"},
AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
ExposedHeaders:   []string{"Authorization"},
AllowCredentials: true,
MaxAge:           300,
}))

api := httpapi.NewAPI(pool, jwtSvc)
api.Mount(r)

srv := &http.Server{
Addr:         ":" + cfg.Port,
Handler:      r,
ReadTimeout:  10 * time.Second,
WriteTimeout: 10 * time.Second,
IdleTimeout:  60 * time.Second,
}

go func() {
log.Printf("API listening on http://localhost:%s\n", cfg.Port)
if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
log.Fatalf("listen: %v", err)
}
}()

quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
_ = srv.Shutdown(ctx)
log.Println("API stopped")
}


