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

CREATE INDEX "bookmarks_text_search_index" ON "bookmarks" (
	"title",
	"description",
	"url"
);

CREATE INDEX "tags_search_name" ON "tags" (
	"name"
);