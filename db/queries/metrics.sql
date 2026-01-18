-- name: CreateMetrics :one
INSERT INTO metrics (
    installation,
    issues_closed,
    duplicate_issue_found,
    issues_colaborated_with
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;


-- name: GetMetricsByInstallation :one
SELECT *
FROM metrics
WHERE installation = $1;

-- name: ListMetrics :many
SELECT *
FROM metrics
ORDER BY created_at DESC;

-- name: IncrementIssuesClosed :one
UPDATE metrics
SET issues_closed = issues_closed + $2,
    updated_at = NOW()
WHERE installation = $1
RETURNING *;

-- name: IncrementDuplicateIssuesFound :one
UPDATE metrics
SET duplicate_issue_found = duplicate_issue_found + $2,
    updated_at = NOW()
WHERE installation = $1
RETURNING *;

-- name: IncrementIssuesCollaboratedWith :one
UPDATE metrics
SET issues_colaborated_with = issues_colaborated_with + $2,
    updated_at = NOW()
WHERE installation = $1
RETURNING *;

-- name: UpsertMetrics :one
INSERT INTO metrics (
    installation,
    issues_closed,
    duplicate_issue_found,
    issues_colaborated_with
) VALUES (
    $1, $2, $3, $4
)
ON CONFLICT (installation)
DO UPDATE SET
    issues_closed = metrics.issues_closed + EXCLUDED.issues_closed,
    duplicate_issue_found = metrics.duplicate_issue_found + EXCLUDED.duplicate_issue_found,
    issues_colaborated_with = metrics.issues_colaborated_with + EXCLUDED.issues_colaborated_with,
    updated_at = NOW()
RETURNING *;

-- name: DeleteMetricsByInstallation :exec
DELETE FROM metrics
WHERE installation = $1;

-- name: GetTotalMetrics :one
SELECT
    COALESCE(SUM(issues_closed), 0)::BIGINT AS issues_closed,
    COALESCE(SUM(duplicate_issue_found), 0)::BIGINT AS duplicate_issue_found,
    COALESCE(SUM(issues_colaborated_with), 0)::BIGINT AS issues_colaborated_with
FROM metrics;
