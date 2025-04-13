-- +goose Up
alter table orders add column status uuid not null references order_status(id) on delete cascade;

-- +goose Down
alter table orders drop column status;