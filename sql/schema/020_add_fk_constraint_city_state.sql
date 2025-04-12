-- +goose Up
alter table cities add column state uuid not null references states(id) on delete cascade;

-- +goose Down
alter table drop column state;