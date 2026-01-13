package httpapi

import (
"net/http"

"recruitflow/apps/api/internal/auth"
"recruitflow/apps/api/internal/store"

"github.com/go-chi/chi/v5"
"github.com/jackc/pgx/v5/pgxpool"
)

type API struct {
pool   *pgxpool.Pool
q      *store.Queries
jwtSvc *auth.JWT
}

func NewAPI(pool *pgxpool.Pool, jwtSvc *auth.JWT) *API {
return &API{
pool:   pool,
q:      store.New(pool),
jwtSvc: jwtSvc,
}
}

func (a *API) Mount(r chi.Router) {
r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
w.WriteHeader(http.StatusOK)
w.Write([]byte("ok"))
})

r.Route("/v1", func(r chi.Router) {
r.Post("/auth/register", a.handleRegister)
r.Post("/auth/login", a.handleLogin)

r.Group(func(r chi.Router) {
r.Use(AuthMiddleware(a.jwtSvc))
r.Get("/me", a.handleMe)

r.Get("/jobs", a.handleListJobs)
r.Post("/jobs", a.handleCreateJob)
r.Put("/jobs/{id}", a.handleUpdateJob)
r.Delete("/jobs/{id}", a.handleDeleteJob)
})
})
}
