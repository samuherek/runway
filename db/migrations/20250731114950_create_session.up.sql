create table if not exists sessions(
    id UUID primary key,
    user_id UUID not null references users(id) on delete cascade,
    token text unique not null,
    ip_address text,
    user_agent text,
    last_seen_at timestamptz,
    expires_at timestamptz not null
);

create index if not exists session_unique_idx on sessions(token);

