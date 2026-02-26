ALTER TABLE users
DROP COLUMN email_verified;

DROP INDEX IF EXISTS idx_email_verification_token_hash;
DROP TABLE IF EXISTS email_verification_tokens;