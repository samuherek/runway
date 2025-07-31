CREATE TABLE IF NOT EXISTS users(
    id UUID primary key,
    email varchar(128) not null unique,
    verified_at timestamptz,
    created_at timestamptz not null default current_timestamp
);

create unique index if not exists email_unique_idx on users(lower(email));
