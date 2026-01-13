package httpapi

import (
"context"
"net/http"
"strings"

"recruitflow/apps/api/internal/auth"
)

type ctxKey string

const userIDKey ctxKey = "user_id"

func UserIDFromContext(ctx context.Context) (int64, bool) {
v := ctx.Value(userIDKey)
id, ok := v.(int64)
return id, ok
}

func AuthMiddleware(jwtSvc *auth.JWT) func(http.Handler) http.Handler {
return func(next http.Handler) http.Handler {
fn := func(w http.ResponseWriter, r *http.Request) {
h := r.Header.Get("Authorization")
if h == "" || !strings.HasPrefix(h, "Bearer ") {
http.Error(w, "missing bearer token", http.StatusUnauthorized)
return
}
token := strings.TrimPrefix(h, "Bearer ")
claims, err := jwtSvc.Parse(token)
if err != nil {
http.Error(w, "invalid token", http.StatusUnauthorized)
return
}
ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
next.ServeHTTP(w, r.WithContext(ctx))
}
return http.HandlerFunc(fn)
}
}
