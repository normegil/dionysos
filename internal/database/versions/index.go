package versions

import "database/sql"

type Version interface {
	Apply() error
	Rollback() error
}

type Versions []Version

func NewVersions(db *sql.DB) Versions {
	versions := make([]Version, 0)

	versions = append(versions, VersionCreation{DB: db})
	versions = append(versions, SchemaCreation{DB: db})

	return versions
}

func (v Versions) HighestVersion() int {
	return len(v) - 1
}
