-- +goose Up

create table "user"
(
    username   text primary key,
    email      text      not null,
    password   text      not null,
    role       integer   not null references user_role (id),
    created_at timestamp not null default now(),
    updated_at timestamp not null default now()
);

-- +goose Down

drop table if exists "user";