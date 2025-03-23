-- +goose Up
alter table sellers add column user_id uuid not null references users(id) on delete cascade;

-- +goose Down
alter table sellers drop column user_id;