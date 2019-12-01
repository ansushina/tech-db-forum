
DROP TABLE IF EXISTS "forums" CASCADE;
DROP TABLE IF EXISTS "posts" CASCADE;
DROP TABLE IF EXISTS "threads" CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;
DROP TABLE IF EXISTS "votes" CASCADE;
DROP TABLE IF EXISTS "forum_users" CASCADE;

CREATE TABLE users (
  "nickname" TEXT UNIQUE PRIMARY KEY,
  "fullname" TEXT NOT NULL,
  "about"    TEXT,
  "email"    TEXT UNIQUE NOT NULL,
);

CREATE TABLE forums (
  "title"   TEXT    NOT NULL,
  "user"    TEXT    NOT NULL REFERENCES users ("nickname"),
  "slug"    TEXT    UNIQUE NOT NULL,
  "posts"   BIGINT  DEFAULT 0,
  "threads" INTEGER DEFAULT 0,
);

CREATE TABLE threads (
  "id"      SERIAL         UNIQUE PRIMARY KEY,
  "author"  TEXT           NOT NULL REFERENCES users ("nickname"),
  "created" TIMESTAMPTZ(3) DEFAULT now(),
  "forum"   TEXT           NOT NULL REFERENCES forums ("slug"),
  "message" TEXT           NOT NULL,
  "slug"    TEXT,
  "title"   TEXT           NOT NULL,
  "votes"   INTEGER        DEFAULT 0
); 

CREATE TABLE posts (
  "id"       BIGSERIAL      UNIQUE PRIMARY KEY,
  "author"   TEXT           NOT NULL REFERENCES users ("nickname"),
  "created"  TIMESTAMPTZ(3) DEFAULT now(),
  "forum"    TEXT           NOT NULL REFERENCES forums ("slug"),
  "isEdited" BOOLEAN        DEFAULT FALSE,
  "message"  TEXT           NOT NULL,
  "parent"   INTEGER        DEFAULT 0,
  "thread"   INTEGER        NOT NULL REFERENCES threads ("id")
);

CREATE TABLE votes (
  "thread"   INT NOT NULL REFERENCES threads("id"),
  "voice"    INTEGER NOT NULL,
  "nickname" CITEXT   NOT NULL
);