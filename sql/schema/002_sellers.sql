-- +goose Up
create table Sellers(
    id varchar(10) primary key,
    user_id uuid not null references users(id) on delete cascade,
    gst_number varchar(15) not null unique,
    pan_number varchar(10) not null unique,
    pickup_address text not null,
    bank_account_holder_name text not null,
    bank_account_number text not null,
    ifsc_code varchar(11) not null unique,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table Sellers;