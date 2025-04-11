-- name: CreateCharacteristics :one
insert into characteristics(id, description, product_id, created_at, updated_at)
values(
    gen_random_uuid(),
    $1,
    $2,
    NOW(),
    NOW()
)
returning *;

-- name: UpdateCharacteristics :one
update characteristics set description = $1, updated_at = NOW() where id = $2 and product_id = $3
returning description;

-- name: GetProductCharacteristic :one
select description from characteristics where id = $1 and product_id = $2;

-- name: GetAllProductCharacteristics :many
select * from characteristics where product_id = $1;