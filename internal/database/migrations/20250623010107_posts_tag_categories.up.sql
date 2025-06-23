CREATE TABLE post_tags (
  post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
  tag_id UUID REFERENCES tags(id) ON DELETE CASCADE,
  PRIMARY KEY (post_id, tag_id)
);

CREATE TABLE post_categories (
  post_id UUID REFERENCES posts(id) ON DELETE CASCADE,
  category_id UUID REFERENCES categories(id) ON DELETE CASCADE,
  PRIMARY KEY (post_id, category_id)
);

-- Index for fast lookup for filtering
-- Will add later to test performance improvements
-- CREATE INDEX idx_post_tags_post_id ON post_tags(post_id);
-- CREATE INDEX idx_post_categories_post_id ON post_categories(post_id);
-- CREATE INDEX idx_post_tags_tag_id ON post_tags(tag_id);
-- CREATE INDEX idx_post_categories_category_id ON post_categories(category_id);
