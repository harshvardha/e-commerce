-- +goose Up
create table order_status(
    id uuid primary key, 
    status text not null, 
    created_at timestamp not null, 
    updated_at timestamp not null
);

-- +goose Down
drop table order_status;