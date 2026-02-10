-- Football Bot Database Schema
-- Version: 001_initial

-- +goose Up

-- Users table
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    username VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    is_global_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_active_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_telegram_id ON users(telegram_id);

-- Communities (group chats)
CREATE TABLE communities (
    id SERIAL PRIMARY KEY,
    telegram_chat_id BIGINT UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_by BIGINT REFERENCES users(id),
    settings JSONB DEFAULT '{
        "team_size": 5,
        "max_subs_per_team": 1,
        "timezone": "UTC",
        "rating_params": ["speed", "dribbling", "shot"]
    }',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_communities_chat_id ON communities(telegram_chat_id);

-- Community members
CREATE TABLE community_members (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    community_id INTEGER REFERENCES communities(id) ON DELETE CASCADE,
    role VARCHAR(20) DEFAULT 'player' CHECK (role IN ('admin', 'player')),
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, community_id)
);

CREATE INDEX idx_community_members_user ON community_members(user_id);
CREATE INDEX idx_community_members_community ON community_members(community_id);

-- Venues (playing locations)
CREATE TABLE venues (
    id SERIAL PRIMARY KEY,
    community_id INTEGER REFERENCES communities(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    contact_type VARCHAR(20) CHECK (contact_type IN ('telegram', 'whatsapp', 'phone')),
    contact_value VARCHAR(255),
    created_by BIGINT REFERENCES users(id),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_venues_community ON venues(community_id);

-- Events (games)
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    community_id INTEGER REFERENCES communities(id) ON DELETE CASCADE,
    venue_id INTEGER REFERENCES venues(id) ON DELETE SET NULL,
    event_date DATE NOT NULL,
    event_time TIME NOT NULL,
    team_size INTEGER DEFAULT 5,
    max_subs_per_team INTEGER DEFAULT 1,
    status VARCHAR(20) DEFAULT 'open' CHECK (status IN ('open', 'full', 'in_progress', 'finished', 'cancelled', 'not_held')),
    team_a INTEGER[] DEFAULT '{}',
    team_b INTEGER[] DEFAULT '{}',
    subs_a INTEGER[] DEFAULT '{}',
    subs_b INTEGER[] DEFAULT '{}',
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP
);

CREATE INDEX idx_events_community ON events(community_id);
CREATE INDEX idx_events_date ON events(event_date);
CREATE INDEX idx_events_status ON events(status);

-- Participants
CREATE TABLE participants (
    id SERIAL PRIMARY KEY,
    event_id INTEGER REFERENCES events(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    guest_name VARCHAR(255),
    guest_claim_token VARCHAR(64) UNIQUE,
    added_by BIGINT REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'confirmed' CHECK (status IN ('confirmed', 'cancelled')),
    registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT participant_user_or_guest CHECK (
        (user_id IS NOT NULL AND guest_name IS NULL) OR 
        (user_id IS NULL AND guest_name IS NOT NULL)
    )
);

CREATE INDEX idx_participants_event ON participants(event_id);
CREATE INDEX idx_participants_user ON participants(user_id);
CREATE INDEX idx_participants_claim_token ON participants(guest_claim_token);

-- Player ratings (per community)
CREATE TABLE player_ratings (
    id SERIAL PRIMARY KEY,
    community_id INTEGER REFERENCES communities(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    speed DECIMAL(3,2) CHECK (speed >= 1 AND speed <= 5),
    dribbling DECIMAL(3,2) CHECK (dribbling >= 1 AND dribbling <= 5),
    shot DECIMAL(3,2) CHECK (shot >= 1 AND shot <= 5),
    games_played INTEGER DEFAULT 0,
    ratings_received INTEGER DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(community_id, user_id)
);

CREATE INDEX idx_player_ratings_community ON player_ratings(community_id);
CREATE INDEX idx_player_ratings_user ON player_ratings(user_id);

-- Rating votes (individual ratings after games)
CREATE TABLE rating_votes (
    id SERIAL PRIMARY KEY,
    event_id INTEGER REFERENCES events(id) ON DELETE CASCADE,
    from_user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    to_user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    speed INTEGER CHECK (speed >= 1 AND speed <= 5),
    dribbling INTEGER CHECK (dribbling >= 1 AND dribbling <= 5),
    shot INTEGER CHECK (shot >= 1 AND shot <= 5),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(event_id, from_user_id, to_user_id)
);

CREATE INDEX idx_rating_votes_event ON rating_votes(event_id);
CREATE INDEX idx_rating_votes_to_user ON rating_votes(to_user_id);

-- +goose Down

DROP TABLE IF EXISTS rating_votes;
DROP TABLE IF EXISTS player_ratings;
DROP TABLE IF EXISTS participants;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS venues;
DROP TABLE IF EXISTS community_members;
DROP TABLE IF EXISTS communities;
DROP TABLE IF EXISTS users;
