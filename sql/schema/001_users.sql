-- +goose Up
create table users(
    id uuid primary key,
    email text not null unique,
    phone_number varchar(13) not null unique,
    password text not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table users;