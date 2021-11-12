/*
This table is to keep API keys created for the user
The API keys can be used to create forms programmatically
*/
CREATE TABLE IF NOT EXISTS api_keys (
  key TEXT PRIMARY KEY ,
  user_id TEXT NOT NULL,
  created_at timestamp NOT NULL  DEFAULT  NOW()
)