-- +goose Up
create table users(
    id uuid primary key,
    email text,
    phone_number varchar(13),
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table users;