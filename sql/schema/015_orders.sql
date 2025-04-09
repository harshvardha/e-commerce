-- +goose Up
create table orders(
    id uuid primary key,
    total_value float not null,
    seller_id text not null references sellers(id) on delete cascade,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table orders;