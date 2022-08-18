create table movies (
id serial not null unique,
title varchar(255) not null,
release varchar(255) not null,
streaming_service varchar(255) not null,
saved_at timestamp not null default now()
);