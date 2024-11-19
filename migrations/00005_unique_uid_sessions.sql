-- +goose Up
-- +goose StatementBegin
-- строго по одной сессии для пользователя
CREATE UNIQUE INDEX sessions_unique_uid ON sessions (uid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX sessions_unique_uid;
-- +goose StatementEnd
