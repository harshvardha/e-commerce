-- +goose Up
create table cities(
    id uuid primary key,
    name text not null unique,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table cities;