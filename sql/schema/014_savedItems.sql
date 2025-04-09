-- +goose Up
create table saved_items(
    user_id uuid not null references users(id) on delete cascade,
    product_id uuid not null references products(id) on delete cascade,
    created_at timestamp not null,
    updated_at timestamp not null,
    unique(user_id, product_id)
);

-- +goose Down
drop table saved_items;