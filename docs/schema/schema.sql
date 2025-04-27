CREATE TABLE "bookmarks" (
	"id"	INTEGER NOT NULL,
	"title"	TEXT NOT NULL,
	"description"	TEXT,
	"thumbnail"	TEXT,
	"url"	TEXT,
	"created_at"	INTEGER,
	"updated_at"	INTEGER,
	PRIMARY KEY("id" AUTOINCREMENT)
)

CREATE TABLE "tags" (
	"id"	INTEGER NOT NULL,
	"name"	TEXT NOT NULL,
	"created_at"	INTEGER,
	"updated_at"	INTEGER,
	PRIMARY KEY("id" AUTOINCREMENT)
);

CREATE TABLE "bookmarks_tags" (
	"bookmark_id"	INTEGER NOT NULL,
	"tag_id"	INTEGER NOT NULL,
	"created_at"	INTEGER,
	FOREIGN KEY("bookmark_id") REFERENCES "bookmarks"("id") ON DELETE CASCADE,
	FOREIGN KEY("tag_id") REFERENCES "tags" ON DELETE CASCADE
);

-- Users table
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at INTEGER,
	updated_at INTEGER
);

-- Access tokens table
CREATE TABLE tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    token TEXT NOT NULL UNIQUE,
	expires_at INTEGER NOT NULL,
    created_at INTEGER,
	updated_at INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Refresh tokens table (optional, for refresh token support)
CREATE TABLE refresh_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    refresh_token TEXT NOT NULL UNIQUE,
    expires_at INTEGER NOT NULL,
    created_at INTEGER,
    updated_at INTEGER,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE INDEX "bookmarks_text_search_index" ON "bookmarks" (
	"title",
	"description",
	"url"
);

CREATE INDEX "tags_search_name" ON "tags" (
	"name"
);