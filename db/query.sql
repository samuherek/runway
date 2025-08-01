
-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
select * from users where email = $1;

-- name: CreateUser :one
insert into users(id, email) values($1, $2) returning *;

-- name: SetUserVerified :one
update users set verified_at = $2 where id = $1 returning *;

-- name: CreateTempToken :one
insert into temp_tokens(id, expires_at, user_id, value) values($1, $2, $3, $4) returning *;

-- name: GetTempToken :one
select * from temp_tokens where value = $1;

-- name: SetTempTokenUsed :one
update temp_tokens set used = true where value = $1 returning *;

-- name: CreateSession :one
insert into sessions(id, user_id, token, ip_address, user_agent, expires_at) values($1, $2, $3, $4, $5, $6) returning *;

-- name: GetSessionByToken :one
select * from sessions where token = $1 and expires_at > $2;
