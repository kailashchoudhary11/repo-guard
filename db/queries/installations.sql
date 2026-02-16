-- name: CreateInstallation :one
INSERT INTO installations (
    installation_id,
    config_data,
    updated_by,
    installed_by
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, installation_id, config_data, updated_by, installed_by;


-- name: GetInstallationByID :one
SELECT
    id,
    installation_id,
    config_data,
    updated_by,
    installed_by
FROM installations
WHERE installation_id = $1;


-- name: GetInstallationByInstallationID :one
SELECT
    id,
    installation_id,
    config_data,
    updated_by,
    installed_by
FROM installations
WHERE installation_id = $1;


-- name: ListInstallations :many
SELECT
    id,
    installation_id,
    config_data,
    updated_by,
    installed_by
FROM installations
ORDER BY id DESC
LIMIT $1 OFFSET $2;


-- name: UpdateInstallationConfig :one
UPDATE installations
SET
    config_data = $2,
    updated_by = $3
WHERE id = $1
RETURNING id, installation_id, config_data, updated_by, installed_by;


-- name: UpsertInstallationConfig :one
INSERT INTO installations (
    installation_id,
    config_data,
    updated_by,
    installed_by
) VALUES (
    $1, $2, $3, $3
)
ON CONFLICT (installation_id) DO UPDATE
SET
    config_data = EXCLUDED.config_data,
    updated_by = EXCLUDED.updated_by
RETURNING id, installation_id, config_data, updated_by, installed_by;


-- name: DeleteInstallation :exec
DELETE FROM installations
WHERE id = $1;
