CREATE TABLE matches (
    uuid UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    match_state INTEGER NOT NULL,
	created_by_user VARCHAR(255) NOT NULL
);

ALTER TABLE matches
ADD CONSTRAINT match_state_valid
CHECK (match_state IN (0, 1, 2));
