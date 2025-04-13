-- name: CreateOrder :one
insert into orders(id, total_value, seller_id, status, created_at, updated_at)
values(
    gen_random_uuid(),
    $1,
    $2,
    $3,
    NOW(),
    NOW()
)
returning *;

-- name: UpdateOrderStatus :exec
update orders set status = $1 where id = $2;

-- name: DeleteOrder :exec
delete from orders where id = $1;

-- name: AddProductToOrder :exec
insert into orders_users_products(order_id, user_id, product_id, quantity, created_at, updated_at)
values(
    $1,
    $2,
    $3,
    $4,
    NOW(),
    NOW()
);

-- name: GetOrderDetails :many
select products.name, products.description, products.price, products.image_urls, orders_users_products.quantity from orders_users_products join products on orders_users_products.product_id = products.id where orders_users_products.user_id = $1 and orders_users_products.order_id = $2; 

-- name: GetAllOrders :many
select orders_users_products.order_id, products.name, products.image_urls from orders_users_products join products on orders_users_products.product_id = products.id where user_id = $1;