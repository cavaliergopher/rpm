package yum

// RepoDatabase represents an entry in a repository metadata file for an
// individual database file such as primary_db or filelists_db.
type RepoDatabase struct {
	Type            string               `xml:"type,attr"`
	Location        RepoDatabaseLocation `xml:"location"`
	Timestamp       int                  `xml:"timestamp"`
	Size            int                  `xml:"size"`
	Checksum        RepoDatabaseChecksum `xml:"checksum"`
	OpenSize        int                  `xml:"open-size"`
	OpenChecksum    RepoDatabaseChecksum `xml:"open-checksum"`
	DatabaseVersion int                  `xml:"database_version"`
}

// RepoDatabaseLocation represents the URI, relative to a package repository,
// of a repository database.
type RepoDatabaseLocation struct {
	Href string `xml:"href,attr"`
}

func (c *RepoDatabase) String() string {
	return c.Type
}
