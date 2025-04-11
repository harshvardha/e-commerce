-- name: ListProduct :one
insert into products(id, name, description, price, image_urls, 
    stock_amount, store_id, category_id, created_at, updated_at)
    values(
        gen_random_uuid(),
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        NOW(),
        NOW()
    )
returning *;

-- name: UpdateProduct :one
update products
set
name = $1,
description = $2,
price = $3,
image_urls = $4,
stock_amount = $5,
category_id = $6,
updated_at = NOW()
where id = $7
returning *;

-- name: RemoveProduct :exec
delete from products where id = $1 and store_id = $2;

-- name: GetProductById :one
select * from products where id = $1;

-- name: GetProductsByCategory :many
select * from products where category_id = $1;

-- name: GetProductsByStoreId :many
select * from products where store_id = $1;