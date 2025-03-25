-- +goose Up
create table Products(
    id uuid primary key,
    name text not null,
    description json not null,
    price float not null,
    image_urls json not null,
    stock_amount int not null default 0,
    store_id uuid not null references stores(id) on delete cascade,
    category_id uuid not null references categories(id) on delete cascade,
    created_at timestamp not null,
    updated_at timestamp not null
);

-- +goose Down
drop table Products;