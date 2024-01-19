CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE channeltype AS ENUM ('user', 'channel');

CREATE TABLE IF NOT EXISTS channel (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email_id VARCHAR(256) NOT NULL,
    channel_type channeltype NOT NULL
);

CREATE TABLE IF NOT EXISTS membership (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    subscriber_id UUID NOT NULL REFERENCES channel(id),
    subscribed_to_id UUID NOT NULL REFERENCES channel(id)
);

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content VARCHAR(2048), 
    sender_id UUID NOT NULL REFERENCES channel(id),
    receiver_id UUID NOT NULL REFERENCES channel(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);