-- +goose Up
create table Characteristics(
    id uuid primary key,
    description json not null,
    product_id uuid not null references products(id) on delete cascade,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table Characteristics;