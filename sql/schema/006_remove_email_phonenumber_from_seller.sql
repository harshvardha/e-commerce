-- +goose Up
alter table sellers drop column email, drop column phone_number;

-- +goose Down
alter table sellers add column email text, add column phone_number varchar(10);