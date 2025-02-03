CREATE TYPE user_status AS ENUM ('Confirmed', 'Unconfirmed');
CREATE TABLE users(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    name VARCHAR(255) UNIQUE NOT NULL,
    status user_status DEFAULT 'Unconfirmed'
);
CREATE TABLE friends(
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    user_id uuid REFERENCES users(id) ON DELETE CASCADE,
    friend_id uuid REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY(user_id, friend_id)
);