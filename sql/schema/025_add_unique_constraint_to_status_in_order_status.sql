-- +goose Up
alter table order_status add constraint status_unique unique(status);

-- +goose Down
alter table order_status drop constraint status_unique;