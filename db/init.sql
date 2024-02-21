CREATE TABLE user_account(
    id bigserial PRIMARY KEY,
    identifier char(16) NOT NULL,
    display_name varchar(64) NOT NULL,
    hash char(64) NOT NULL,
    email varchar(128) NOT NULL,
    bio varchar(512)
);

CREATE TABLE room(
    id bigserial PRIMARY KEY,
    name varchar(16) NOT NULL
);

CREATE TABLE user_room(
    user_id bigserial,
    room_id bigserial,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES user_account(id),
    CONSTRAINT fk_room FOREIGN KEY(room_id) REFERENCES room(id),
    PRIMARY KEY(user_id, room_id)
);

CREATE TABLE room_message(
    id bigserial PRIMARY KEY,
    user_id bigserial NOT NULL,
    room_id bigserial NOT NULL,
    message_text varchar(512) NOT NULL,
    utc bigserial NOT NULL,
    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES user_account(id),
    CONSTRAINT fk_room FOREIGN KEY(room_id) REFERENCES room(id)
);
