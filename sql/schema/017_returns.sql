-- +goose Up
create table returns(
    id uuid primary key,
    seller_id text not null references sellers(id) on delete cascade,
    total_value float not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table returns;