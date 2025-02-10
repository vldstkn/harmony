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
CREATE TABLE rooms(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    creator_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
CREATE TYPE room_role AS ENUM ('OWNER', 'MODERATOR', 'MEMBER', 'BANNED');
CREATE TABLE rooms_members(
    room_id BIGINT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    role INT DEFAULT 0,
    PRIMARY KEY (room_id, user_id),
    CHECK ( role in (-1,0,1,2) )
);

CREATE TYPE message_status AS ENUM ('READ', 'UNREAD');
CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    room_id BIGINT NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    sender_id BIGINT REFERENCES users(id) NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    status message_status DEFAULT 'UNREAD'
);
CREATE TYPE notification_status AS ENUM('SEEN', 'UNSEEN');
CREATE TYPE notification_type As ENUM('SYSTEM', 'COMMON');
CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    status notification_status DEFAULT 'UNSEEN',
    text TEXT NOT NULL,
    type notification_type NOT NULL
);