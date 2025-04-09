-- +goose Up
alter table users add column firstname text not null, add column lastname text not null;

-- +goose Down
alter table users drop column firstname, drop column lastname;