CREATE TABLE IF NOT EXISTS otp_log
(
    id             SERIAL PRIMARY KEY,
    profile_id     INT   NOT NULL,
    otp            VARCHAR(10)   NOT NULL,
    validated      BOOL DEFAULT 'f',
    created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expiry         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() + interval '15 minute'
)