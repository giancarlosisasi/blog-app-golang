CREATE TABLE posts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR(255) NOT NULL,
  slug VARCHAR(255) UNIQUE NOT NULL,
  content TEXT NOT NULL,
  -- excerpt (extracto) is a short description of the post
  excerpt TEXT,
  featured_image_url TEXT,
  status VARCHAR(20) DEFAULT 'draft', -- 'draft', 'published', 'archived'
  author_id UUID REFERENCES users(id) ON DELETE CASCADE,
  published_at TIMESTAMP WITH TIME ZONE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_posts_author_id ON posts(author_id);
CREATE INDEX idx_posts_slug ON posts(slug);
-- CREATE INDEX idx_posts_status ON posts(status);