-- name: RegisterSeller :one
insert into sellers(
    id,
    gst_number,
    pan_number,
    pickup_address,
    bank_account_holder_name,
    bank_account_number,
    ifsc_code,
    user_id,
    created_at,
    updated_at
) values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    NOW(),
    NOW()
)
returning *;

-- name: UpdateSellerTaxAndAddress :one
update sellers set gst_number = $1, pan_number = $2, pickup_address = $3, updated_at = NOW() where id = $4
returning id, gst_number, pan_number, pickup_address, created_at, updated_at;

-- name: UpdateSellerBankDetails :one
update sellers set bank_account_holder_name = $1, bank_account_number = $2, ifsc_code = $3, updated_at = NOW() where id = $4
returning id, bank_account_holder_name, bank_account_number, ifsc_code, created_at, updated_at;

-- name: DeleteSellerAccount :exec
delete from sellers where id = $1;

-- name: GetSellerContactInfo :one
select sellers.id,
users.email, 
users.phone_number, 
sellers.created_at, 
sellers.updated_at 
from sellers join users on sellers.user_id = users.id where sellers.id = $1;

-- name: GetSellerTaxAndAddressInfo :one
select id, gst_number, pan_number, pickup_address, created_at, updated_at from sellers where id = $1;

-- name: GetSellerBankDetails :one
select id, bank_account_holder_name, bank_account_number, ifsc_code, created_at, updated_at from sellers where id = $1;