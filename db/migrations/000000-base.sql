CREATE TABLE "organizations" (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX organizations_name_idx ON organizations (name) WHERE deleted_at IS NULL;


CREATE TABLE "projects" (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX projects_name_idx ON projects (name) WHERE deleted_at IS NULL;

CREATE TABLE "items" (
  id BIGSERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  state BOOLEAN DEFAULT true, -- false is closed
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE "discussions" (
  id BIGSERIAL PRIMARY KEY,
  body TEXT NOT NULL,
  item BIGINT REFERENCES items(id),
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE "comments" (
  id BIGSERIAL PRIMARY KEY,
  body TEXT NOT NULL,
  item BIGINT REFERENCES items(id),
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE "members" (
  id BIGSERIAL PRIMARY KEY,
  username TEXT NOT NULL,
  name TEXT NOT NULL,
  email TEXT NOT NULL,
  password BYTEA NOT NULL,
  salt BYTEA NOT NULL,
  email_verified BOOLEAN DEFAULT false,
  active BOOLEAN DEFAULT false,
  created_at TIMESTAMP WITH TIME ZONE,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX members_username_idx ON members (username) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX members_email_idx ON members (email) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX members_salt_idx ON members (salt) WHERE deleted_at IS NULL;
