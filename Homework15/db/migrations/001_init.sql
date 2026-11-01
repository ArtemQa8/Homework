-- +goose Up
create table players (
    id bigint primary key,
    name text not null,
    created_at timestamptz not null default now()
);
create table games (
    id bigint primary key,
    player1_id bigint not null references players(id) on delete cascade,
    player1_name text not null,
    player1_color int not null,
    player2_id bigint not null references players(id) on delete cascade,
    player2_name text not null,
    player2_color int not null,
    board jsonb,
    current_color int not null,
    white_king_moved boolean not null default false,
    black_king_moved boolean not null default false,
    white_rook_a_moved boolean not null default false,
    white_rook_h_moved boolean not null default false,
    black_rook_a_moved boolean not null default false,
    black_rook_h_moved boolean not null default false,
    created_at timestamptz not null default now()
);
create table moves (
    id bigint primary key,
    game_id bigint not null references games(id) on delete cascade,
    from_row int not null,
    from_col int not null,
    to_row int not null,
    to_col int not null,
    captured jsonb,
    promotion int not null default 0,
    moved_piece jsonb,
    "check" boolean not null default false,
    mate boolean not null default false
);
create table counters (
    name text primary key,
    seq bigint not null default 0
);
-- +goose Down
drop table if exists moves;
drop table if exists games;
drop table if exists players;
drop table if exists counters;