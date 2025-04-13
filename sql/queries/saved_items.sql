-- name: AddProductToSavedItems :exec
insert into saved_items(user_id, product_id, created_at, updated_at)
values(
    $1,
    $2,
    NOW(),
    NOW()
);

-- name: RemoveProductFromSavedItems :exec
delete from saved_items where user_id = $1 and product_id = $2;

-- name: GetAllSavedItems :many
select products.name, products.description, products.price, products.image_urls from saved_items join products on saved_items.product_id = products.id where saved_items.user_id = $1;