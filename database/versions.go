package database

var (
	// Versions contains all known versions.
	// This will be filled with init functions.
	Versions = make([]Version, 0)
)

// Version is a defined version.
type Version struct {
	Version   uint64
	Statement string
}

// ByVersion implements the sort.Interface based on the Version field.
type ByVersion []Version

func (versions ByVersion) Len() int {
	return len(versions)
}

func (versions ByVersion) Less(i, j int) bool {
	return versions[i].Version < versions[j].Version
}

func (versions ByVersion) Swap(i, j int) {
	versions[i], versions[j] = versions[j], versions[i]
}
