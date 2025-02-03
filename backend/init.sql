CREATE TYPE user_status AS ENUM ('Confirmed', 'Unconfirmed');
CREATE TABLE users(
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    name VARCHAR(255) UNIQUE NOT NULL,
    status user_status DEFAULT 'Unconfirmed'
);
CREATE TABLE friendships  (
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    user_id1 BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_id2 BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY(user_id1, user_id2),
    CHECK (user_id1 < user_id2)
);
