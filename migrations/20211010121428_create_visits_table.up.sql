CREATE TABLE IF NOT EXISTS url_visits
(
    id         SERIAL PRIMARY KEY,
    url_id     INTEGER   NOT NULL,
    ip         INET      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),

    FOREIGN KEY (url_id) REFERENCES short_urls (id)
);

ALTER TABLE IF EXISTS short_urls
    DROP COLUMN IF EXISTS visits,
    ADD COLUMN enabled BOOLEAN NOT NULL DEFAULT true,
    ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT now();
