-- +goose Up
ALTER TABLE games ADD COLUMN last_move_id BIGINT NOT NULL DEFAULT 0;
ALTER TABLE players ADD COLUMN rating INT NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE games DROP COLUMN last_move_id;
ALTER TABLE players DROP COLUMN rating;