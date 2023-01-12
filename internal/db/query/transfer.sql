-- name: CreateTransfer :one
INSERT INTO transfers (from_account_id, to_account_id, amount) VALUES ($1, $2, $3) RETURNING *;

-- name: GetTransferById :one
SELECT * FROM transfers WHERE id = $1 LIMIT 1;

-- name: GetTransferByFromAccountId :one
SELECT * FROM transfers WHERE from_account_id = $1 LIMIT 1;

-- name: GetTransferByToAccountId :one
SELECT * FROM transfers WHERE to_account_id = $1 LIMIT 1;

-- name: ListsTransfers :many
SELECT * FROM transfers ORDER BY id LIMIT $1 OFFSET $2;

-- name: GetListsTransfers :many
SELECT * FROM transfers WHERE from_account_id = $1 OR to_account_id = $1 ORDER BY id LIMIT $2 OFFSET $3;

-- name: GetTotalPageListsTransfers :one
SELECT COUNT(*) FROM transfers;

-- name: GetTotalPageListsTransfersSpesific :one
SELECT COUNT(*) FROM transfers WHERE from_account_id = $1 OR to_account_id = $1;