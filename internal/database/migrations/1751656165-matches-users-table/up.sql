CREATE TABLE users_matches (
    user_uuid UUID REFERENCES users(uuid) ON DELETE CASCADE,
    match_uuid UUID REFERENCES matches(uuid) ON DELETE CASCADE,
    PRIMARY KEY (user_uuid, match_uuid)
);
