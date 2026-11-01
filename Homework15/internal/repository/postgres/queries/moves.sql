-- name: CreateMove :one
INSERT INTO moves (
    id, game_id, from_row, from_col, to_row, to_col, captured, promotion, moved_piece, "check", mate
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetMove :one
SELECT * FROM moves WHERE id = $1;

-- name: ListMoves :many
SELECT * FROM moves ORDER BY id;

-- name: UpdateMove :exec
UPDATE moves SET 
    game_id = $2, from_row = $3, from_col = $4, to_row = $5, to_col = $6, captured = $7, promotion = $8, moved_piece = $9, "check" = $10, mate = $11
WHERE id = $1;

-- name: DeleteMove :exec
DELETE FROM moves WHERE id = $1;

-- name: ListMovesByGameID :many
SELECT * FROM moves WHERE game_id = $1 ORDER BY id;