-- name: CreateState :one
insert into states(id, name, created_at, updated_at)
values(
    gen_random_uuid(),
    $1,
    NOW(),
    NOW()
)
returning *;

-- name: UpdateState :one
update states set name = $1, updated_at = NOW() where id = $2
returning name, created_at, updated_at;

-- name: RemoveState :exec
delete from states where id = $1;

-- name: GetStateName :one
select name from states where id = $1;

-- name: GetAllStates :many
select * from states;