CREATE TABLE  users (
  id TEXT PRIMARY KEY NOT NULL UNIQUE,
  email TEXT NOT NULL UNIQUE,
  avatar_url TEXT NOT NULL DEFAULT '',
  first_name TEXT NOT NULL DEFAULT '',
  last_name TEXT NOT NULL DEFAULT '',
  encrypted_password TEXT NOT NULL DEFAULT '',
  confirmed_at timestamp NULL DEFAULT NULL,
  invited_at timestamp NULL DEFAULT NULL,
  confirmation_token TEXT NOT NULL DEFAULT '',
  confirmation_sent_at timestamp NULL DEFAULT NULL,
  recovery_token TEXT NOT NULL DEFAULT '',
  recovery_sent_at timestamp NULL DEFAULT NULL,
  email_change_token TEXT NOT NULL DEFAULT '',
  email_change TEXT NOT NULL DEFAULT '',
  email_change_sent_at timestamp NULL DEFAULT NULL,
  last_sign_in_at timestamp NULL DEFAULT NULL,
  user_info jsonb NOT NULL DEFAULT '{}',
  created_at timestamp NOT NULL  DEFAULT  NOW(),
  updated_at timestamp NOT NULL  DEFAULT  NOW()
) ;
  
CREATE INDEX  users_email_idx  ON users (email);

CREATE OR REPLACE FUNCTION set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE set_timestamp();

