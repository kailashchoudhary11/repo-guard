-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS metrics (
    id uuid DEFAULT gen_random_uuid () PRIMARY KEY,

    installation VARCHAR(255) NOT NULL,

    issues_closed BIGINT NOT NULL DEFAULT 0,
    duplicate_issue_found BIGINT NOT NULL DEFAULT 0,
    issues_colaborated_with BIGINT NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (installation)
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS metrics;
-- +goose StatementEnd
