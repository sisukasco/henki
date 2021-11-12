CREATE TABLE IF NOT EXISTS refresh_tokens (
  id BIGSERIAL PRIMARY KEY,
  token TEXT NOT NULL DEFAULT '',
  user_id TEXT NOT NULL,
  revoked BOOLEAN NOT NULL DEFAULT FALSE,
  created_at timestamp NOT NULL  DEFAULT  NOW(),
  updated_at timestamp NOT NULL  DEFAULT  NOW()
) ;

CREATE INDEX  IF NOT EXISTS refresh_tokens_user_id_idx  ON refresh_tokens (user_id);

CREATE INDEX  IF NOT EXISTS refresh_tokens_token_idx  ON refresh_tokens (token);


CREATE TRIGGER refresh_tokens_set_timestamp
BEFORE UPDATE ON refresh_tokens
FOR EACH ROW
EXECUTE PROCEDURE set_timestamp();