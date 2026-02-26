-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
    id serial primary key,
    name varchar(255) not null,
    email varchar(255) not null unique,
    password varchar(255) not null,
    registered_at timestamp default current_timestamp
);

CREATE TABLE refresh_tokens(
    id serial not null unique,
    user_id int references users(id) on delete cascade not null,
    token varchar(255) not null unique,
    expires_at timestamp not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE refresh_tokens;

DROP TABLE users;
-- +goose StatementEnd
