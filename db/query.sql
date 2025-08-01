
-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
select * from users where email = $1;

-- name: CreateUser :one
insert into users(id, email) values($1, $2) returning *;

-- name: CreateTempToken :one
insert into temp_tokens(id, expires_at, user_id, value) values($1, $2, $3, $4) returning *;
