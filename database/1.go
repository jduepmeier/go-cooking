package database

func init() {
	Versions = append(Versions, Version{
		Version: 1,
		Statement: `
			CREATE TABLE user(
				username TEXT PRIMARY KEY NOT NULL,
				password TEXT NOT NULL
			);
		`,
	})
}
