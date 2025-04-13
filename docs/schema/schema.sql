CREATE TABLE "bookmarks" (
	"id"	INTEGER NOT NULL,
	"title"	TEXT NOT NULL,
	"description"	TEXT,
	"thumbnail"	TEXT,
	"url"	TEXT,
	"created_at"	INTEGER,
	"updated_at"	INTEGER,
	PRIMARY KEY("id" AUTOINCREMENT)
);

CREATE TABLE "tags" (
	"id"	INTEGER NOT NULL,
	"name"	TEXT NOT NULL,
	"created_at"	INTEGER,
	"updated_at"	INTEGER,
	PRIMARY KEY("id" AUTOINCREMENT)
);

CREATE TABLE "bookmark_tag" (
	"bookmark_id"	INTEGER NOT NULL,
	"tag_id"	INTEGER NOT NULL,
	"created_at"	INTEGER,
	"updated_at"	INTEGER,
	FOREIGN KEY("bookmark_id") REFERENCES "bookmarks"("id") ON DELETE CASCADE,
	FOREIGN KEY("tag_id") REFERENCES "tags" ON DELETE CASCADE
)