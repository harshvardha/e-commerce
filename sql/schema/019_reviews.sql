-- +goose Up
create table reviews(
    id uuid primary key,
    description text not null,
    user_id uuid not null references users(id) on delete cascade,
    product_id uuid not null references products(id) on delete cascade,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table reviews;