-- +goose Up
create table Stores(
    id uuid primary key,
    name text not null unique,
    owner_id varchar(10) not null references sellers(id) on delete cascade,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table Stores;