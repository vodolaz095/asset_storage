-- +goose Up
-- +goose StatementBegin
-- таблица сессий
create table if not exists sessions (
                                        id text primary key default encode(gen_random_bytes(16),'hex'),
                                        uid bigint not null, -- user id
                                        created_at timestamptz not null default now(),
                                        FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table sessions;
-- +goose StatementEnd
