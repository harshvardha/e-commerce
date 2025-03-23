-- +goose Up
create table Sellers(
    id uuid primary key,
    email text,
    phone_number varchar(10),
    gst_number varchar(15) not null,
    pan_number varchar(10) not null,
    pickup_address text not null,
    bank_account_holder_name text not null,
    bank_account_number text not null,
    ifsc_code text not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table Sellers;