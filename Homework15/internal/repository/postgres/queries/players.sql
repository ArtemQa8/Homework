-- name: CreatePlayer :one
INSERT INTO players (id, name, rating) VALUES ($1, $2, $3) RETURNING *;

-- name: GetPlayer :one
SELECT * FROM players WHERE id = $1;

-- name: ListPlayers :many
SELECT * FROM players ORDER BY id;

-- name: UpdatePlayer :exec
UPDATE players SET name = $2 WHERE id = $1;

-- name: DeletePlayer :exec
DELETE FROM players WHERE id = $1;

-- name: GetPlayerByName :many
SELECT * FROM players WHERE name = $1;