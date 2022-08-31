CREATE TABLE movies (
    id serial not null unique,
    title varchar(255) not null,
    release varchar(255) not null,
    streaming_service varchar(255) not null,
    saved_at timestamp not null default now()
);

CREATE TABLE users (
    id serial not null unique,
    name varchar(255) not null,
    email varchar(255) not null unique,
    password varchar(255) not null,
    registered_at timestamp not null default now()
);

CREATE TABLE refresh_tokens (
    id serial not null unique,
    user_id int not null unique,
    token varchar(255) not null,
    expires_at timestamp not null default now()
)