package httpapi

import (
"encoding/json"
"net/http"

"recruitflow/apps/api/internal/store"

"golang.org/x/crypto/bcrypt"
)

type authReq struct {
Email    string `json:"email"`
Password string `json:"password"`
}

type authResp struct {
Token string `json:"token"`
}

func (a *API) handleRegister(w http.ResponseWriter, r *http.Request) {
var req authReq
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, "bad json", http.StatusBadRequest)
return
}
if req.Email == "" || req.Password == "" {
http.Error(w, "email and password required", http.StatusBadRequest)
return
}

hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
if err != nil {
http.Error(w, "hash error", http.StatusInternalServerError)
return
}

u, err := a.q.CreateUser(r.Context(), store.CreateUserParams{
Email:        req.Email,
PasswordHash: string(hash),
})
if err != nil {
http.Error(w, "email already exists", http.StatusConflict)
return
}

tok, err := a.jwtSvc.Sign(u.ID)
if err != nil {
http.Error(w, "token error", http.StatusInternalServerError)
return
}

w.Header().Set("Content-Type", "application/json")
_ = json.NewEncoder(w).Encode(authResp{Token: tok})
}

func (a *API) handleLogin(w http.ResponseWriter, r *http.Request) {
var req authReq
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, "bad json", http.StatusBadRequest)
return
}
if req.Email == "" || req.Password == "" {
http.Error(w, "email and password required", http.StatusBadRequest)
return
}

u, err := a.q.GetUserByEmail(r.Context(), req.Email)
if err != nil {
http.Error(w, "invalid credentials", http.StatusUnauthorized)
return
}

if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
http.Error(w, "invalid credentials", http.StatusUnauthorized)
return
}

tok, err := a.jwtSvc.Sign(u.ID)
if err != nil {
http.Error(w, "token error", http.StatusInternalServerError)
return
}

w.Header().Set("Content-Type", "application/json")
_ = json.NewEncoder(w).Encode(authResp{Token: tok})
}
