-- Role Tables
CREATE TABLE roles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) UNIQUE NOT NULL,
  description TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE permissions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(100) UNIQUE NOT NULL, -- 'posts.create', 'posts.update', 'posts.delete', 'posts.read'
  resource VARCHAR(100) NOT NULL, -- 'posts', 'comments'
  action VARCHAR(100) NOT NULL, --'create', 'update', 'delete', 'read'
  description TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- junction table
CREATE TABLE role_permissions(
  role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
  permission_id UUID REFERENCES permissions(id) ON DELETE CASCADE,
  -- can't create two times the same permission for a role
  PRIMARY KEY (role_id, permission_id)
);

-- User-role, junction table
CREATE TABLE user_roles (
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
  assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  assigned_by UUID REFERENCES users(id),
  PRIMARY KEY (user_id, role_id)
);

-- Index for fast permission checking
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_permissions_name ON permissions(name);
CREATE INDEX idx_roles_name ON roles(name);

-- Composite index for fastest permission lookup
CREATE INDEX idx_user_permission_lookup ON user_roles(user_id)
  INCLUDE (role_id);
