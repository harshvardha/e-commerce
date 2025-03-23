-- +goose Up
create table Stores(
    id uuid primary key,
    name text not null,
    owner_id uuid not null references sellers(id) on delete cascade,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table Stores;