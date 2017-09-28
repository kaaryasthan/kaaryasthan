CREATE TABLE "projects" (
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE, 
  deleted_at TIMESTAMP WITH TIME ZONE,
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT NOT NULL
);

CREATE UNIQUE INDEX projects_name_idx ON projects (name) WHERE deleted_at IS NULL;

CREATE TABLE "items" (
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE, 
  deleted_at TIMESTAMP WITH TIME ZONE,
  id BIGSERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT NOT NULL
);

CREATE TABLE "comments" (
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE, 
  deleted_at TIMESTAMP WITH TIME ZONE,
  id BIGSERIAL PRIMARY KEY,
  body TEXT NOT NULL,
  item BIGINT REFERENCES items(id)
);

CREATE TABLE "members" (
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE, 
  deleted_at TIMESTAMP WITH TIME ZONE,
  id BIGSERIAL PRIMARY KEY,
  username TEXT NOT NULL,
  name TEXT NOT NULL,
  email TEXT NOT NULL,
  password BYTEA NOT NULL,
  salt BYTEA NOT NULL
);

CREATE UNIQUE INDEX members_username_idx ON members (username) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX members_email_idx ON members (email) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX members_salt_idx ON members (salt) WHERE deleted_at IS NULL;
