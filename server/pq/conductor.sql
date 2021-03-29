CREATE SCHEMA IF NOT EXISTS conductor;

CREATE TYPE OperationStatus AS ENUM ('UNKNOWN', 'SUCCESS', 'FAILURE');

CREATE TABLE IF NOT EXISTS conductor.Operations (
    operation_id uuid,
    done bool,
    payload bytea,
    status OperationStatus
);

CREATE TABLE IF NOT EXISTS conductor.Queue (

);
