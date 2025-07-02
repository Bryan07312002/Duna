CREATE TABLE users (
    id UUID DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL PRIMARY KEY
);
