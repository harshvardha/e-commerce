-- +goose Up
create table orders_users_products(
    order_id uuid not null references orders(id) on delete cascade,
    user_id uuid not null references users(id) on delete cascade,
    product_id uuid not null references products(id) on delete cascade,
    created_at timestamp not null,
    updated_at timestamp not null,
    unique(order_id, user_id, product_id)
);

-- +goose Down
drop table orders_users_products;