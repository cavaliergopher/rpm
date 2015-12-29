package yum

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

type RepoDatabaseLocation struct {
	Href string `xml:"href,attr"`
}

func (c *RepoDatabase) String() string {
	return c.Type
}
