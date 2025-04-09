-- +goose Up
create table returns_users_products(
    return_id uuid not null references returns(id) on delete cascade,
    user_id uuid not null references users(id) on delete cascade,
    product_id uuid not null references products(id) on delete cascade,
    created_at timestamp not null,
    updated_at timestamp not null,
    unique(return_id, user_id, product_id)
);

-- +goose Down
drop table returns_users_products;