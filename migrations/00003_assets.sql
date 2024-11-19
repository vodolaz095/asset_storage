-- +goose Up
-- +goose StatementBegin

-- таблица с файлами
create table if not exists assets (
                                      name text not null,
                                      uid bigint not null, -- user id
                                      data bytea not null,
                                      created_at timestamptz not null default now(),
                                      primary key (name, uid),
                                      FOREIGN KEY (uid) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table assets;
-- +goose StatementEnd
