-- name: CreateReview :one
insert into reviews(id, description, user_id, product_id, created_at, updated_at)
values(
    gen_random_uuid(),
    $1,
    $2,
    $3,
    NOW(),
    NOW()
)
returning *;

-- name: UpdateReview :one
update reviews set description = $1, updated_at = NOW() where id = $2 and user_id = $3
returning description, created_at, updated_at;

-- name: RemoveReview :exec
delete from reviews where id = $1 and user_id = $2;

-- name: GetReviewsByUserID :many
select * from reviews where user_id = $1;

-- name: GetReviewsByProductID :many
select * from reviews where product_id = $1;