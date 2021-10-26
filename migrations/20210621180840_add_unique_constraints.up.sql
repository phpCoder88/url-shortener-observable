ALTER TABLE IF EXISTS short_urls ADD CONSTRAINT short_urls_token_unique UNIQUE (token);
ALTER TABLE IF EXISTS short_urls ADD CONSTRAINT short_urls_long_url_unique UNIQUE (long_url);
