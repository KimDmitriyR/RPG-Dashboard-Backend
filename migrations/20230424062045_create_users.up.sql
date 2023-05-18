CREATE TABLE users (
    id bigserial not null primary key,
    email varchar not null unique,
    name_user varchar not null, 
    encrypted_password varchar not null
);

CREATE TABLE users_information (
    id bigserial not null primary key,
    role varchar not null,
    user_level integer not null default 1
);
