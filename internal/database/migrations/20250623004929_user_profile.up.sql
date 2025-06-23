CREATE TABLE user_profiles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  first_name VARCHAR(100),
  last_name VARCHAR(100),
  bio TEXT,
  avatar_url TEXT,
  website_url TEXT,
  -- location VARCHAR(100),
  -- is_public BOOLEAN DEFAULT TRUE,
  -- is_verified BOOLEAN DEFAULT FALSE,
  -- is_active BOOLEAN DEFAULT TRUE,
  -- is_deleted BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_user_profiles_user_id ON user_profiles(user_id);