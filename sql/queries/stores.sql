-- name: CreateStore :one
insert into Stores(id, name, owner_id, created_at, updated_at)
values (
    gen_random_uuid(),
    $1,
    $2,
    NOW(),
    NOW()
)
returning *;

-- name: UpdateStoreInformation :one
update Stores set name = $1, updated_at = NOW() where id = $2
returning name, updated_at;

-- name: GetStoreInformation :one
select name, created_at, updated_at from Stores where id = $1;

-- name: DeleteStore :exec
delete from stores where id = $1;