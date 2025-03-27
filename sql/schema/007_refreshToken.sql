-- +goose Up
create table refresh_token(
    token text not null,
    user_id uuid not null references users(id) on delete cascade,
    expires_at timestamp not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table refresh_token;