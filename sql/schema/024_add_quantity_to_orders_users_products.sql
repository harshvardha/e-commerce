-- +goose Up
alter table orders_users_products add column quantity int not null default 1;

-- +goose Down
alter table orders_users_products drop column quantity;