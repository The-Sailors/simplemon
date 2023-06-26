CREATE TABLE IF NOT EXISTS monitors (
    monitor_id SERIAL NOT NULL,
    user_email TEXT NOT NULL,
    type TEXT NOT NULL,
    url TEXT NOT NULL,
    method TEXT NOT NULL,
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    body TEXT,
    headers TEXT,
    parameters TEXT,
    description TEXT,
    frequency_minutes INTEGER,
    threshold_minutes INTEGER,
    CONSTRAINT monitor_id PRIMARY KEY(user_email, type, url, method)
);