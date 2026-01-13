package httpapi

import (
"encoding/json"
"net/http"
)

func (a *API) handleMe(w http.ResponseWriter, r *http.Request) {
uid, ok := UserIDFromContext(r.Context())
if !ok {
http.Error(w, "unauthorized", http.StatusUnauthorized)
return
}
w.Header().Set("Content-Type", "application/json")
_ = json.NewEncoder(w).Encode(map[string]any{"user_id": uid})
}
