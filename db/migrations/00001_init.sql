-- +goose Up
-- +goose StatementBegin

CREATE TABLE owners (
id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
name STRING NOT NULL,
origin STRING NOT NULL,
service STRING NOT NULL,
created_at TIMESTAMPTZ NOT NULL,
updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE records (
id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
record STRING NOT NULL,
type STRING NOT NULL,
created_at TIMESTAMPTZ NOT NULL,
updated_at TIMESTAMPTZ NOT NULL,
UNIQUE INDEX idx_record_record_type (record, type)
);

CREATE TABLE answers (
id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
target STRING NOT NULL,
type STRING NOT NULL,
ttl INT DEFAULT 3600 NOT NULL CHECK (ttl >= 0),
has_details BOOL NOT NULL,
owner_id UUID NOT NULL REFERENCES owners(id) ON DELETE CASCADE ON UPDATE CASCADE,
record_id UUID NOT NULL REFERENCES records(id) ON DELETE CASCADE ON UPDATE CASCADE,
created_at TIMESTAMPTZ NOT NULL,
updated_at TIMESTAMPTZ NOT NULL,
INDEX idx_record_owner (record_id, owner_id),
INDEX idx_record_target_type (record_id, target, type),
UNIQUE INDEX idx_record_owner_target_type (record_id, owner_id, target, type)
);

CREATE TABLE answer_details (
id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
answer_id UUID NOT NULL REFERENCES answers(id) ON DELETE CASCADE ON UPDATE CASCADE,
port INT CHECK (port >= 0),
priority INT CHECK (priority >= 0),
protocol STRING,
weight INT CHECK (weight >= 0),
created_at TIMESTAMPTZ NOT NULL,
updated_at TIMESTAMPTZ NOT NULL,
UNIQUE INDEX idx_answer_id (answer_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE answers;
DROP TABLE records;
DROP TABLE owners;
DROP TABLE answer_details;

-- +goose StatementEnd
