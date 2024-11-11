CREATE TABLE role_users (
    role_id UUID NOT NULL,
    user_id UUID NOT NULL,

    PRIMARY KEY (role_id, user_id),
    FOREIGN KEY (role_id) REFERENCES roles(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);