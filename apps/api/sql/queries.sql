-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING id, email, password_hash, created_at;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, created_at
FROM users
WHERE email = $1;

-- name: CreateJob :one
INSERT INTO jobs (user_id, company, title, link, status, salary, notes, follow_up_date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, user_id, company, title, link, status, salary, notes, follow_up_date, created_at, updated_at;

-- name: ListJobs :many
SELECT id, user_id, company, title, link, status, salary, notes, follow_up_date, created_at, updated_at
FROM jobs
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateJob :one
UPDATE jobs
SET company = $2,
    title = $3,
    link = $4,
    status = $5,
    salary = $6,
    notes = $7,
    follow_up_date = $8,
    updated_at = now()
WHERE id = $1 AND user_id = $9
RETURNING id, user_id, company, title, link, status, salary, notes, follow_up_date, created_at, updated_at;

-- name: DeleteJob :exec
DELETE FROM jobs
WHERE id = $1 AND user_id = $2;
