CREATE TABLE IF NOT EXISTS otp_log
(
    id             SERIAL PRIMARY KEY,
    profile_id     INT   NOT NULL,
    otp            INTEGER   NOT NULL,
    validated      BOOL DEFAULT 'f',
    created_at     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    expiry         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() + interval '15 minute'
)