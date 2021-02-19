create table if not exists animator (
    id bigserial primary key,
    uuid text not null,
    username text not null,
    password text not null,
    display_name text not null,
    active boolean default true not null,
    created_at timestamp with time zone default now() not null,
    last_login_at timestamp with time zone
);

create table if not exists animation (
    id bigserial primary key,
    uuid text not null,
    animator_id bigint not null references animator(id),
    title text not null,
    details text,
    published boolean not null default false,
    visible boolean not null default false,
    animation text not null,
    created_at timestamp with time zone default now() not null,
    published_at timestamp with time zone,
    updated_at timestamp with time zone default now() not null
);