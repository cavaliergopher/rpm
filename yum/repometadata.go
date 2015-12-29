package yum

import (
	"encoding/xml"
	"fmt"
	"io"
)

type RepoMetadata struct {
	XMLName  xml.Name `xml:"repomd"`
	XMLNS    string   `xml:"xmlns,attr"`
	XMLNSRPM string   `xml:"xmlns:rpm,attr"`

	Revision int                `xml:"revision"`
	Data     []RepoMetadataItem `xml:"data"`
}

type RepoMetadataItem struct {
	Type            string                   `xml:"type,attr"`
	Location        RepoMetadataItemLocation `xml:"location"`
	Timestamp       int                      `xml:"timestamp"`
	Size            int                      `xml:"size"`
	Checksum        RepoMetadataItemChecksum `xml:"checksum"`
	OpenSize        int                      `xml:"open-size"`
	OpenChecksum    RepoMetadataItemChecksum `xml:"open-checksum"`
	DatabaseVersion int                      `xml:"database_version"`
}

type RepoMetadataItemChecksum struct {
	Type string `xml:"type,attr"`
	Hash string `xml:",chardata"`
}

type RepoMetadataItemLocation struct {
	Href string `xml:"href,attr"`
}

func ReadRepoMetadata(r io.Reader) (*RepoMetadata, error) {
	md := RepoMetadata{
		Data: make([]RepoMetadataItem, 0),
	}

	decoder := xml.NewDecoder(r)
	err := decoder.Decode(&md)

	if err != nil {
		return nil, fmt.Errorf("Error decoding repository metadata: %v", err)
	}

	return &md, nil
}

func (c *RepoMetadata) Write(w io.Writer) error {
	c.XMLNS = "http://linux.duke.edu/metadata/repo"
	c.XMLNSRPM = "http://linux.duke.edu/metadata/rpm"

	encoder := xml.NewEncoder(w)
	err := encoder.Encode(c)
	if err != nil {
		return fmt.Errorf("Error encoding repository metadata: %v", err)
	}

	return nil
}
