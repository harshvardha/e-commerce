-- name: AddProductToCart :exec
insert into carts(user_id, product_id, created_at, updated_at)
values(
    $1,
    $2,
    NOW(),
    NOW()
);

-- name: UpdateProductQuantity :exec
update carts set quantity = $1 where user_id = $2 and product_id = $3;

-- name: RemoveProductFromCart :exec
delete from carts where user_id = $1 and product_id = $2;

-- name: GetAllProductsInCart :many
select products.name, products.description, products.price, products.image_urls from carts join products on carts.product_id = products.id where carts.user_id = $1;

-- name: EmptyCart :exec
delete from carts where user_id = $1;