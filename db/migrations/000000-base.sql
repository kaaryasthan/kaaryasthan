CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

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

CREATE FUNCTION update_items_num_column()
RETURNS TRIGGER AS $$
DECLARE
nextValue BIGINT;
BEGIN
  SELECT COALESCE(MAX(num), 0)+1 INTO nextValue FROM items;
  NEW.num = nextValue;
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE FUNCTION update_item_discussion_comment_search_tsv_column()
RETURNS TRIGGER AS $$
DECLARE
rec RECORD;
BEGIN
  IF TG_TABLE_NAME = 'items' THEN
    SELECT i.title, i.description, string_agg(c.body, ' ') AS comments, string_agg(DISTINCT d.body, ' ') AS discussions, i.id INTO rec
      FROM items i LEFT JOIN discussions d ON i.id=d.item_id LEFT JOIN comments c ON c.discussion_id = d.id
      WHERE i.id=(SELECT DISTINCT i2.id FROM items i2 LEFT JOIN discussions d2 ON i2.id=d2.item_id LEFT JOIN comments c2 ON c2.discussion_id=d2.id WHERE i2.id=NEW.id)
      GROUP BY i.id;
  ELSIF TG_TABLE_NAME = 'discussions' THEN
    SELECT i.title, i.description, string_agg(c.body, ' ') AS comments, string_agg(DISTINCT d.body, ' ') AS discussions, i.id INTO rec
      FROM items i LEFT JOIN discussions d ON i.id=d.item_id LEFT JOIN comments c ON c.discussion_id = d.id
      WHERE i.id=(SELECT DISTINCT i2.id FROM items i2 LEFT JOIN discussions d2 ON i2.id=d2.item_id LEFT JOIN comments c2 ON c2.discussion_id=d2.id WHERE d2.id=NEW.id)
      GROUP BY i.id;
  ELSIF TG_TABLE_NAME = 'comments' THEN
    SELECT i.title, i.description, string_agg(c.body, ' ') AS comments, string_agg(DISTINCT d.body, ' ') AS discussions, i.id INTO rec
      FROM items i LEFT JOIN discussions d ON i.id=d.item_id LEFT JOIN comments c ON c.discussion_id = d.id
      WHERE i.id=(SELECT DISTINCT i2.id FROM items i2 LEFT JOIN discussions d2 ON i2.id=d2.item_id LEFT JOIN comments c2 ON c2.discussion_id=d2.id WHERE c2.id=NEW.id)
      GROUP BY i.id;
  END IF;
  INSERT INTO item_discussion_comment_search
    (item_id, tsv) VALUES (rec.id,
      setweight(to_tsvector('pg_catalog.english', rec.title), 'A') ||
      setweight(to_tsvector('pg_catalog.english', rec.description), 'B') ||
      setweight(to_tsvector('pg_catalog.english', rec.discussions), 'C') ||
      setweight(to_tsvector('pg_catalog.english', rec.comments), 'D')
      )
  ON CONFLICT (item_id)
  DO UPDATE SET tsv =
    setweight(to_tsvector('pg_catalog.english', rec.title), 'A') ||
    setweight(to_tsvector('pg_catalog.english', rec.description), 'B') ||
    setweight(to_tsvector('pg_catalog.english', rec.discussions), 'C') ||
    setweight(to_tsvector('pg_catalog.english', rec.comments), 'D')
    WHERE item_discussion_comment_search.item_id=rec.id;
  RETURN NEW;
END
$$ LANGUAGE 'plpgsql';

CREATE FUNCTION notify_item_change() RETURNS TRIGGER AS $$
BEGIN
   RAISE NOTICE 'NEW.id is currently %', NEW.id::text;
   -- Execute pg_notify(channel, notification)
   PERFORM pg_notify('item_change', NEW.id::text);

   -- Result is ignored since this is an AFTER trigger
   RETURN NULL;
END;
$$ LANGUAGE 'plpgsql';

CREATE TYPE role AS ENUM ('admin', 'manager', 'member');

CREATE TABLE "users" (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
  username TEXT NOT NULL,
  name TEXT NOT NULL,
  email TEXT NOT NULL,
  personal_note TEXT DEFAULT '' NOT NULL,
  user_role role DEFAULT 'member' NOT NULL,
  password BYTEA NOT NULL,
  salt BYTEA NOT NULL,
  email_verified BOOLEAN DEFAULT false NOT NULL,
  active BOOLEAN DEFAULT false NOT NULL,
  tsv TSVECTOR,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_unq_users_id ON users (id); -- TODO: is this required?
CREATE UNIQUE INDEX idx_unq_users_username ON users (username);
CREATE UNIQUE INDEX idx_unq_users_email ON users (email) WHERE deleted_at IS NULL;
CREATE UNIQUE INDEX idx_unq_users_salt ON users (salt);

CREATE TRIGGER trgr_update_users_updated_at_column BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE INDEX idx_gin_users_fulltext_search ON users USING GIN (tsv);

CREATE FUNCTION update_users_tsv_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.tsv :=
    setweight(to_tsvector('pg_catalog.english', NEW.username), 'A') ||
    setweight(to_tsvector('pg_catalog.english', NEW.name), 'B') ||
    setweight(to_tsvector('pg_catalog.english', NEW.email),'C');
  RETURN NEW;
END
$$ LANGUAGE 'plpgsql';

CREATE TRIGGER trgr_update_users_tsv_column BEFORE INSERT OR UPDATE OF username, name, email ON users FOR EACH ROW EXECUTE PROCEDURE update_users_tsv_column();

CREATE TABLE "configurations" (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  value TEXT DEFAULT '' NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_unq_configurations_name ON configurations (name) WHERE deleted_at IS NULL;

CREATE TRIGGER trgr_update_configurations_updated_at_column BEFORE UPDATE ON configurations FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE TABLE "projects" (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT DEFAULT '' NOT NULL,
  item_template TEXT DEFAULT '' NOT NULL,
  archived BOOLEAN DEFAULT false NOT NULL,
  created_by UUID REFERENCES users(id) NOT NULL,
  updated_by UUID REFERENCES users(id),
  deleted_by UUID REFERENCES users(id),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_unq_projects_name ON projects (name);

CREATE TRIGGER trgr_update_projects_updated_at_column BEFORE UPDATE ON projects FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE TABLE "labels" (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  color TEXT NOT NULL CHECK(color ~ '^#[a-f0-9]{6}$'),
  project_id BIGINT REFERENCES projects(id) NOT NULL,
  created_by UUID REFERENCES users(id) NOT NULL,
  updated_by UUID REFERENCES users(id),
  deleted_by UUID REFERENCES users(id),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_unq_labels_name_project ON labels (name, project_id) WHERE deleted_at IS NULL;

CREATE TRIGGER trgr_update_labels_updated_at_column BEFORE UPDATE ON labels FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE TABLE "items" (
  id BIGSERIAL PRIMARY KEY,
  num BIGINT NOT NULL,
  title TEXT NOT NULL,
  description TEXT DEFAULT '' NOT NULL,
  open_state BOOLEAN DEFAULT true NOT NULL, -- false is closed
  project_id BIGINT REFERENCES projects(id) NOT NULL,
  lock_conversation BOOLEAN DEFAULT false NOT NULL,
  created_by UUID REFERENCES users(id) NOT NULL,
  updated_by UUID REFERENCES users(id),
  deleted_by UUID REFERENCES users(id),
  assignees UUID[],
  subscribers UUID[],
  labels TEXT[],
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_unq_items_num_id ON items (num);

CREATE TRIGGER trgr_update_items_num_column BEFORE INSERT ON items FOR EACH ROW EXECUTE PROCEDURE update_items_num_column();
CREATE TRIGGER trgr_update_items_updated_at_column BEFORE UPDATE ON items FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER trgr_update_items_fulltext_search AFTER INSERT OR UPDATE OF title, description ON items FOR EACH ROW EXECUTE PROCEDURE update_item_discussion_comment_search_tsv_column();
CREATE TRIGGER trgr_update_items_notify AFTER INSERT OR UPDATE OF title, description ON items FOR EACH ROW EXECUTE PROCEDURE notify_item_change();

CREATE TABLE "milestones" (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT DEFAULT '' NOT NULL,
  items BIGINT[],
  project_id BIGINT REFERENCES projects(id) NOT NULL,
  created_by UUID REFERENCES users(id) NOT NULL,
  updated_by UUID REFERENCES users(id),
  deleted_by UUID REFERENCES users(id),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_unq_milestones_name_project_id ON milestones (name, project_id) WHERE deleted_at IS NULL;

CREATE TRIGGER trgr_update_milestones_updated_at_column BEFORE UPDATE ON milestones FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE TABLE "discussions" (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
  body TEXT NOT NULL,
  item_id BIGINT REFERENCES items(id) NOT NULL,
  created_by UUID REFERENCES users(id) NOT NULL,
  updated_by UUID REFERENCES users(id),
  deleted_by UUID REFERENCES users(id),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_unq_discussions_id ON discussions (id); -- TODO: is this required?

CREATE TRIGGER trgr_update_discussions_updated_at_column BEFORE UPDATE ON discussions FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER trgr_update_discussions_fulltext_search AFTER INSERT OR UPDATE OF body ON discussions FOR EACH ROW EXECUTE PROCEDURE update_item_discussion_comment_search_tsv_column();

CREATE TABLE "comments" (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
  body TEXT NOT NULL,
  discussion_id UUID REFERENCES discussions(id) NOT NULL,
  created_by UUID REFERENCES users(id) NOT NULL,
  updated_by UUID REFERENCES users(id),
  deleted_by UUID REFERENCES users(id),
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_unq_comments_id ON comments (id); -- TODO: is this required?

CREATE TRIGGER trgr_update_comments_updated_at_column BEFORE UPDATE ON comments FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();
CREATE TRIGGER trgr_update_comments_fulltext_search AFTER INSERT OR UPDATE OF body ON comments FOR EACH ROW EXECUTE PROCEDURE update_item_discussion_comment_search_tsv_column();

CREATE TABLE "attachments" (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
  content BYTEA NOT NULL,
  name TEXT NOT NULL,
  description TEXT DEFAULT '' NOT NULL,
  item_id BIGINT REFERENCES items(id) NOT NULL,
  attached_by UUID REFERENCES users(id) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_unq_attachments_id ON attachments (id); -- TODO: is this required?
CREATE UNIQUE INDEX idx_unq_attachments_name_item_id ON attachments (name, item_id) WHERE deleted_at IS NULL;

CREATE TRIGGER trgr_update_attachments_updated_at_column BEFORE UPDATE ON attachments FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE TYPE response AS ENUM ('+1', '-1', 'Neutral', 'Laugh', 'Hooray', 'Confused', 'Love');

CREATE TABLE "item_reactions" (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
  reaction response NOT NULL,
  item_id BIGINT REFERENCES items(id) NOT NULL,
  reacted_by UUID REFERENCES users(id) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_unq_item_reactions_id ON item_reactions (id); -- TODO: is this required?

CREATE TRIGGER trgr_update_item_reactions_updated_at_column BEFORE UPDATE ON item_reactions FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE TABLE "discussion_reactions" (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
  reaction response NOT NULL,
  discussion_id UUID REFERENCES discussions(id) NOT NULL,
  reacted_by UUID REFERENCES users(id) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_unq_discussion_reactions_id ON discussion_reactions (id); -- TODO: is this required?

CREATE TRIGGER trgr_update_discussion_reactions_updated_at_column BEFORE UPDATE ON discussion_reactions FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE TABLE "comment_reactions" (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
  reaction response NOT NULL,
  comment_id UUID REFERENCES comments(id) NOT NULL,
  reacted_by UUID REFERENCES users(id) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE UNIQUE INDEX idx_unq_comment_reactions_id ON comment_reactions (id); -- TODO: is this required?

CREATE TRIGGER trgr_update_comment_reactions_updated_at_column BEFORE UPDATE ON comment_reactions FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();

CREATE TABLE "item_discussion_comment_search" (
  item_id BIGINT REFERENCES items(id) NOT NULL,
  tsv TSVECTOR
);

CREATE UNIQUE INDEX idx_unq_item_discussion_comment_search_item_id ON item_discussion_comment_search (item_id);

CREATE INDEX idx_gin_item_discussion_comment_fulltext_search ON item_discussion_comment_search USING GIN (tsv);
