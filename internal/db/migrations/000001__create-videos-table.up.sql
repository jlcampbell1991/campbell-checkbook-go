CREATE TABLE IF NOT EXISTS videos (
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    channel TEXT NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    label TEXT,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    PRIMARY KEY(url)
)