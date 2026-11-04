-- name: CreateGame :one
INSERT INTO games (
    id, player1_id, player1_name, player1_color, 
    player2_id, player2_name, player2_color, board, current_color,
    white_king_moved, black_king_moved, white_rook_a_moved, white_rook_h_moved,
    black_rook_a_moved, black_rook_h_moved, last_move_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
) RETURNING *;

-- name: GetGame :one
SELECT * FROM games WHERE id = $1;

-- name: ListGames :many
SELECT * FROM games ORDER BY id;

-- name: UpdateGame :exec
UPDATE games SET
    player1_id = $2, player1_name = $3,
    player2_id = $4, player2_name = $5
WHERE id = $1;

-- name: OverwriteGame :exec
UPDATE games SET
    player1_id = $2, player1_name = $3, player1_color = $4,
    player2_id = $5, player2_name = $6, player2_color = $7,
    board = $8, current_color = $9,
    white_king_moved = $10, black_king_moved = $11,
    white_rook_a_moved = $12, white_rook_h_moved = $13,
    black_rook_a_moved = $14, black_rook_h_moved = $15,
    last_move_id = $16
WHERE id = $1;

-- name: DeleteGame :exec
DELETE FROM games WHERE id = $1;