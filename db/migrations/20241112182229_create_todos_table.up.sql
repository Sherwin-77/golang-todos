CREATE TABLE todos (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    is_completed BOOLEAN NOT NULL,
    created_at TIMESTAMP(6) WITH TIME ZONE,
    updated_at TIMESTAMP(6) WITH TIME ZONE,

    FOREIGN KEY (user_id) REFERENCES users(id)
);