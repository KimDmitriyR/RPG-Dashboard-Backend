CREATE Table users (
    id bigserial not null primary key,
    email varchar not null unique,
    name_user varchar not null, 
    encrypted_password varchar not null
);