package database

func init() {
	Versions = append(Versions, Version{
		Version: 2,
		Statement: `
			ALTER TABLE recipes ADD description TEXT;
			UPDATE recipes SET description = '';
		`,
	})
}
