-- +goose Up
create table admins(
    id uuid primary key,
    name text not null, 
    email text not null,
    phonenumber text not null,
    password text not null, 
    created_at timestamp not null, 
    updated_at timestamp not null
);

-- +goose Down
drop table admins;