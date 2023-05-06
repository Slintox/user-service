-- +goose Up

insert into user_role (name) values ('user'), ('admin');

-- +goose Down

truncate table user_role;