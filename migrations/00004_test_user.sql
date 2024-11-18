-- +goose Up
-- +goose StatementBegin
-- тестовый пользователь
insert into users
(login, password_hash)
values
    ('alice', encode(digest('secret', 'md5'),'hex'))
on conflict do nothing;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE login='alice';
-- +goose StatementEnd
