CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    club_id INTEGER NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    event_type VARCHAR(20) NOT NULL CHECK (event_type IN ('in_person', 'online')),
    event_date DATE NOT NULL,
    event_time TIME NOT NULL,
    location VARCHAR(255),
    online_link TEXT,
    max_attendees INTEGER,
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_events_club_id ON events(club_id);
CREATE INDEX idx_events_event_date ON events(event_date);
CREATE INDEX idx_events_event_time ON events(event_time);
CREATE INDEX idx_events_deleted_at ON events(deleted_at);

-- Create event_rsvps table
CREATE TABLE event_rsvps (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('going', 'maybe', 'not_going')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, event_id)
);

CREATE INDEX idx_event_rsvps_user_id ON event_rsvps(user_id);
CREATE INDEX idx_event_rsvps_event_id ON event_rsvps(event_id);