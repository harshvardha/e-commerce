-- name: CreateCategory :one
insert into categories(id, name, description, created_at, updated_at)
values(
    gen_random_uuid(),
    $1,
    $2,
    NOW(),
    NOW()
)
returning *;

-- name: UpdateCategory :one
update categories set name = $1, description = $2, updated_at = NOW() where id = $3
returning name, description, created_at, updated_at;

-- name: RemoveCategory :exec
delete from categories where id = $1;

-- name: GetCateogryInformation :one
select name, description, created_at, updated_at from categories where id = $1;