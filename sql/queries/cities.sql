-- name: CreateCity :one
insert into cities(id, name, state, created_at, updated_at)
values(
    gen_random_uuid(),
    $1,
    $2,
    NOW(),
    NOW()
)
returning *;

-- name: UpdateCity :one
update cities set name = $1, state = $2, updated_at = NOW() where id = $3
returning *;

-- name: RemoveCity :exec
delete from cities where id = $1;

-- name: GetCityAndState :one
select cities.name as city_name, states.name as state_name from cities join states on cities.state = states.id where cities.id = $1;

-- name: GetCityById :one
select * from cities where id = $1;

-- name: GetAllCities :many
select * from cities;