-- +goose Up

alter table "user" add column if not exists deleted_at timestamp default null;

-- +goose Down

alter table "user" drop column if exists deleted_at;
