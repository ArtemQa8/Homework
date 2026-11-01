-- name: GetNextSequence :one
INSERT INTO counters (name, seq) VALUES ($1, 1)
ON CONFLICT (name) DO UPDATE SET seq = counters.seq + 1
RETURNING seq;

-- name: SetCounter :exec
INSERT INTO counters (name, seq) VALUES ($1, $2) 
ON CONFLICT (name) DO UPDATE SET seq = EXCLUDED.seq;

-- name: ListCounters :many
SELECT * FROM counters;