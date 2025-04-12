-- name: CreateCustomer :exec
insert into customers(id, created_at, updated_at)
values(
    $1,
    NOW(),
    NOW()
);

-- name: UpdateCustomerAddress :one
update customers set delivery_address = $1, pincode = $2, city = $3, state = $4, updated_at = NOW where id = $5
returning *;

-- name: GetCustomerAddress :one
select delivery_address, pincode, city, state from customers where id = $1;

-- name: GetCustomerInformation :one
select users.email,
users.phone_number,
customers.delivery_address,
customers.pincode,
customers.city,
customers.state
from customers join users on customers.id == users.id where customers.id = $1;