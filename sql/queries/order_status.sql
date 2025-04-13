-- name: AddOrderStatus :one
insert into order_status(id, status, created_at, updated_at)
values(
    gen_random_uuid(),
    $1,
    NOW(),
    NOW()
)
returning *;

-- name: UpdateStatus :one
update order_status set status = $1 where id = $2
returning *;

-- name: GetStatusID :one
select id from order_status where status = $1;

-- name: RemoveStatus :exec
delete from order_status where id = $1;