DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'operationstatus') THEN
        CREATE TYPE OperationStatus AS ENUM ('UNKNOWN', 'SUCCESS', 'FAILURE');
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS Operations (
    operation_id uuid,
    done bool,
    payload bytea,
    status OperationStatus
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'boardshape') THEN
        CREATE TYPE BoardShape AS ENUM ('UNKNOWN', 'STANDARD_31_BY_30');
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS Queue (
    discoverable bool,
    game_code VARCHAR(64),
    min_players int,
    max_players int,
    board_shape BoardShape
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'gamestatus') THEN
        CREATE TYPE GameStatus AS ENUM ('UNKNOWN', 'PLAYING', 'DONE');
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS Games (
    game_id uuid,
    players bytea[],
    status GameStatus,
    round_leaders int[],
    board_shape BoardShape
);

CREATE TABLE IF NOT EXISTS Boards (
    game_id uuid,
    board bytea
);
