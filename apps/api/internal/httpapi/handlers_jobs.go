package httpapi

import (
"encoding/json"
"net/http"
"strconv"
"time"

"recruitflow/apps/api/internal/store"

"github.com/go-chi/chi/v5"
"github.com/jackc/pgx/v5/pgtype"
)

type jobReq struct {
Company      string  `json:"company"`
Title        string  `json:"title"`
Link         string  `json:"link"`
Status       string  `json:"status"`
Salary       string  `json:"salary"`
Notes        string  `json:"notes"`
FollowUpDate *string `json:"follow_up_date"` // "YYYY-MM-DD" or null
}

type jobResp struct {
ID           int64   `json:"id"`
UserID       int64   `json:"user_id"`
Company      string  `json:"company"`
Title        string  `json:"title"`
Link         string  `json:"link"`
Status       string  `json:"status"`
Salary       string  `json:"salary"`
Notes        string  `json:"notes"`
FollowUpDate *string `json:"follow_up_date"`
CreatedAt    string  `json:"created_at"`
UpdatedAt    string  `json:"updated_at"`
}

func (a *API) handleListJobs(w http.ResponseWriter, r *http.Request) {
uid, ok := UserIDFromContext(r.Context())
if !ok {
http.Error(w, "unauthorized", http.StatusUnauthorized)
return
}

jobs, err := a.q.ListJobs(r.Context(), uid)
if err != nil {
http.Error(w, "db error", http.StatusInternalServerError)
return
}

out := make([]jobResp, 0, len(jobs))
for _, j := range jobs {
out = append(out, toJobResp(j))
}

w.Header().Set("Content-Type", "application/json")
_ = json.NewEncoder(w).Encode(out)
}

func (a *API) handleCreateJob(w http.ResponseWriter, r *http.Request) {
uid, ok := UserIDFromContext(r.Context())
if !ok {
http.Error(w, "unauthorized", http.StatusUnauthorized)
return
}

var req jobReq
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, "bad json", http.StatusBadRequest)
return
}
if req.Company == "" || req.Title == "" {
http.Error(w, "company and title required", http.StatusBadRequest)
return
}

fud, err := parseOptionalDate(req.FollowUpDate)
if err != nil {
http.Error(w, "follow_up_date must be YYYY-MM-DD or null", http.StatusBadRequest)
return
}

created, err := a.q.CreateJob(r.Context(), store.CreateJobParams{
UserID:       uid,
Company:      req.Company,
Title:        req.Title,
Link:         def(req.Link),
Status:       defStatus(req.Status),
Salary:       def(req.Salary),
Notes:        def(req.Notes),
FollowUpDate: fud,
})
if err != nil {
http.Error(w, "db error", http.StatusInternalServerError)
return
}

w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated)
_ = json.NewEncoder(w).Encode(toJobResp(created))
}

func (a *API) handleUpdateJob(w http.ResponseWriter, r *http.Request) {
uid, ok := UserIDFromContext(r.Context())
if !ok {
http.Error(w, "unauthorized", http.StatusUnauthorized)
return
}

idStr := chi.URLParam(r, "id")
id, err := strconv.ParseInt(idStr, 10, 64)
if err != nil {
http.Error(w, "invalid id", http.StatusBadRequest)
return
}

var req jobReq
if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
http.Error(w, "bad json", http.StatusBadRequest)
return
}
if req.Company == "" || req.Title == "" {
http.Error(w, "company and title required", http.StatusBadRequest)
return
}

fud, err := parseOptionalDate(req.FollowUpDate)
if err != nil {
http.Error(w, "follow_up_date must be YYYY-MM-DD or null", http.StatusBadRequest)
return
}

updated, err := a.q.UpdateJob(r.Context(), store.UpdateJobParams{
ID:           id,
Company:      req.Company,
Title:        req.Title,
Link:         def(req.Link),
Status:       defStatus(req.Status),
Salary:       def(req.Salary),
Notes:        def(req.Notes),
FollowUpDate: fud,
UserID:       uid,
})
if err != nil {
http.Error(w, "not found", http.StatusNotFound)
return
}

w.Header().Set("Content-Type", "application/json")
_ = json.NewEncoder(w).Encode(toJobResp(updated))
}

func (a *API) handleDeleteJob(w http.ResponseWriter, r *http.Request) {
uid, ok := UserIDFromContext(r.Context())
if !ok {
http.Error(w, "unauthorized", http.StatusUnauthorized)
return
}

idStr := chi.URLParam(r, "id")
id, err := strconv.ParseInt(idStr, 10, 64)
if err != nil {
http.Error(w, "invalid id", http.StatusBadRequest)
return
}

if err := a.q.DeleteJob(r.Context(), store.DeleteJobParams{ID: id, UserID: uid}); err != nil {
http.Error(w, "not found", http.StatusNotFound)
return
}

w.WriteHeader(http.StatusNoContent)
}

// helpers

func def(s string) string { return s }

func defStatus(s string) string {
if s == "" {
return "Saved"
}
return s
}

func parseOptionalDate(s *string) (pgtype.Date, error) {
var d pgtype.Date
if s == nil || *s == "" {
d.Valid = false
return d, nil
}
t, err := time.Parse("2006-01-02", *s)
if err != nil {
return pgtype.Date{}, err
}
d.Time = t
d.Valid = true
return d, nil
}

func toJobResp(j store.Job) jobResp {
var fud *string
if j.FollowUpDate.Valid {
s := j.FollowUpDate.Time.Format("2006-01-02")
fud = &s
}
return jobResp{
ID:           j.ID,
UserID:       j.UserID,
Company:      j.Company,
Title:        j.Title,
Link:         j.Link,
Status:       j.Status,
Salary:       j.Salary,
Notes:        j.Notes,
FollowUpDate: fud,
CreatedAt:    j.CreatedAt.Time.Format(time.RFC3339),
UpdatedAt:    j.UpdatedAt.Time.Format(time.RFC3339),
}
}
