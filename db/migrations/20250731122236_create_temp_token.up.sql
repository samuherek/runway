create table if not exists temp_tokens(
    id uuid primary key,
    expires_at timestamptz not null,
    user_id uuid references users(id),
    used bool default false
);

