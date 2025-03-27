-- name: CreateRefreshToken :exec
insert into refresh_token(token, user_id, expires_at, created_at, updated_at)
values(
    $1,
    $2,
    $3,
    NOW(),
    NOW()
);

-- name: GetRefreshToken :one
select * from refresh_token;

-- name: UpdateRefreshToken :exec
update refresh_token set token = $1, updated_at = NOW() where user_id = $2;