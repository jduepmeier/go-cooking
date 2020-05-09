package database

func init() {
	Versions = append(Versions, Version{
		Version: 0,
		Statement: `
			CREATE TABLE meta(
				key TEXT PRIMARY KEY NOT NULL,
				value TEXT NOT NULL
			);

			CREATE TABLE recipes(
				id INTEGER PRIMARY KEY NOT NULL,
				name TEXT UNIQUE NOT NULL,
				length TEXT NOT NULL,
				freshness INTEGER,
				source TEXT
			);

			INSERT INTO meta(key, value) VALUES ('version', '0');
		`,
	})
}
