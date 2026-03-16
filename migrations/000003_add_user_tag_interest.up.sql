CREATE TABLE user_tag_interest (
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  tag_id BIGINT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
  score INT NOT NULL DEFAULT 0,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, tag_id)
);

CREATE INDEX idx_user_tag_interest_user_score
ON user_tag_interest(user_id, score DESC);