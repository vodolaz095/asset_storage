-- +goose Up
-- +goose StatementBegin
-- добавляем в сессию IP адрес клиента
ALTER TABLE sessions ADD client_ip INET;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE sessions DROP COLUMN client_ip;
-- +goose StatementEnd
