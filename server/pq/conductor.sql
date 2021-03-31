CREATE TYPE OperationStatus AS ENUM ('UNKNOWN', 'SUCCESS', 'FAILURE');

CREATE TABLE IF NOT EXISTS Operations (
    operation_id uuid,
    done bool,
    payload bytea,
    status OperationStatus
);

CREATE TABLE IF NOT EXISTS Queue (

);
