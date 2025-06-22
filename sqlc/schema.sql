-- enable uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


-- Users table
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) UNIQUE NOT NULL,

)