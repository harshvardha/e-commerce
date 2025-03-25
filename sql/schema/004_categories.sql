-- +goose Up
create table Categories(
    id uuid primary key,
    name text not null unique,
    description text not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table Categories;