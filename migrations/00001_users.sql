-- +goose Up
-- +goose StatementBegin
create extension if not exists pgcrypto;
-- таблица с пользователями
create table if not exists users (
                                     id bigserial primary key,
                                     login text not null unique,
                                     password_hash text not null,
                                     created_at timestamptz not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
drop extension pgcrypto;
-- +goose StatementEnd
