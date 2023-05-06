-- +goose Up

create table user_role
(
    id serial primary key,
    name text not null
);

-- +goose Down

drop table if exists user_role;