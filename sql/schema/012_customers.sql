-- +goose Up
create table customers(
    id uuid not null references users(id) on delete cascade,
    delivery_address text not null,
    pincode varchar(6) not null,
    city uuid not null references cities(id) on delete cascade,
    state uuid not null references states(id) on delete cascade,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table customers;