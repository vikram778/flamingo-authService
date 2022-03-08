CREATE TABLE IF NOT EXISTS profile
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(50)   NOT NULL,
    mobile      VARCHAR(20)   NOT NULL,
    dob         DATE NOT NULL,
    location    VARCHAR(50)  NOT NULL,
    is_verified BOOL DEFAULT 'f',
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
)