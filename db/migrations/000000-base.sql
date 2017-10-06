CREATE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   IF row(NEW.*) IS DISTINCT FROM row(OLD.*) THEN
      NEW.updated_at = now();
      RETURN NEW;
   ELSE
      RETURN OLD;
   END IF;
END;
$$ language 'plpgsql';

CREATE TABLE "organizations" (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX organizations_name_idx ON organizations (name) WHERE deleted_at IS NULL;

CREATE TRIGGER update_organizations_updated_at_column BEFORE UPDATE ON organizations FOR EACH ROW EXECUTE PROCEDURE  update_updated_at_column();

CREATE TABLE "projects" (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX projects_name_idx ON projects (name) WHERE deleted_at IS NULL;

CREATE TRIGGER update_projects_updated_at_column BEFORE UPDATE ON projects FOR EACH ROW EXECUTE PROCEDURE  update_updated_at_column();

CREATE TABLE "milestones" (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TRIGGER update_milestones_updated_at_column BEFORE UPDATE ON milestones FOR EACH ROW EXECUTE PROCEDURE  update_updated_at_column();

CREATE TABLE "labels" (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  color TEXT NOT NULL CHECK(color ~ '^#[a-f0-9]{6}$'),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TRIGGER update_labels_updated_at_column BEFORE UPDATE ON labels FOR EACH ROW EXECUTE PROCEDURE  update_updated_at_column();

CREATE TABLE "items" (
  id BIGSERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  state BOOLEAN DEFAULT true, -- false is closed
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TRIGGER update_items_updated_at_column BEFORE UPDATE ON items FOR EACH ROW EXECUTE PROCEDURE  update_updated_at_column();

CREATE TABLE "discussions" (
  id BIGSERIAL PRIMARY KEY,
  body TEXT NOT NULL,
  item BIGINT REFERENCES items(id),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TRIGGER update_discussions_updated_at_column BEFORE UPDATE ON discussions FOR EACH ROW EXECUTE PROCEDURE  update_updated_at_column();

CREATE TABLE "comments" (
  id BIGSERIAL PRIMARY KEY,
  body TEXT NOT NULL,
  item BIGINT REFERENCES items(id),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TRIGGER update_comments_updated_at_column BEFORE UPDATE ON comments FOR EACH ROW EXECUTE PROCEDURE  update_updated_at_column();

CREATE TABLE "members" (
  id BIGSERIAL PRIMARY KEY,
  username TEXT NOT NULL,
  name TEXT NOT NULL,
  email TEXT NOT NULL,
  password BYTEA NOT NULL,
  salt BYTEA NOT NULL,
  email_verified BOOLEAN DEFAULT false,
  active BOOLEAN DEFAULT false,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX members_username_idx ON members (username) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX members_email_idx ON members (email) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX members_salt_idx ON members (salt) WHERE deleted_at IS NULL;

CREATE TRIGGER update_members_updated_at_column BEFORE UPDATE ON members FOR EACH ROW EXECUTE PROCEDURE  update_updated_at_column();
