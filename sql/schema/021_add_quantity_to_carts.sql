-- +goose Up
alter table carts add column quantity int not null default 1;

-- +goose Down
alter table carts drop column quantity;