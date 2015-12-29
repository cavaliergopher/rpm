package yum

import (
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCloneRepo(t *testing.T) {
	baseurl := "http://mirror.centos.org/centos/7/extras/x86_64/"

	repomd_url := baseurl + "/repodata/repomd.xml"
	repohash := fmt.Sprintf("%x", sha1.Sum([]byte(baseurl)))

	cache_dir := fmt.Sprintf("./.yumcache/%v", repohash)

	// create cache folder
	if err := os.MkdirAll(cache_dir, 0750); err != nil && os.IsNotExist(err) {
		t.Fatalf("Error creating cache directory: %v", err)
	}

	// download repository metadata
	repomd := GetRepoMetadata(t, cache_dir, repomd_url)

	// validate or download databases
	for _, d := range repomd.Databases {
		path := GetDatabase(t, baseurl, cache_dir, &d)
		path = DecompressDatabase(t, &d, path)
	}
}

// GetRepoMetadata downloads repository metadata file for the given repository
// URL and caches it in the given cache directory.
func GetRepoMetadata(t *testing.T, cache_dir string, repomd_url string) *RepoMetadata {
	cache_repomd_path := fmt.Sprintf("%s/repomd.xml", cache_dir)

	// open repo metadata from URL
	t.Logf("Downloading repo metadata...")
	resp, err := http.Get(repomd_url)
	if err != nil {
		t.Fatalf("Error retrieving repo metadata from URL: %v", err)
	}
	defer resp.Body.Close()

	// read repometadata into byte buffer
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading repo metadata: %v", err)
	}

	// decode repo metadata into struct
	repomd, err := ReadRepoMetadata(bytes.NewReader(b))
	if err != nil {
		t.Fatalf("Error decoding repo metadata: %v", err)
	}

	// read existing cache
	update_mdcache := false
	f, err := os.Open(cache_repomd_path)
	if err == nil {
		defer f.Close()

		// decode existing cache
		cache_repomd, err := ReadRepoMetadata(f)
		if err != nil {
			t.Fatalf("Error decoding cached repo metadata: %v", err)
		}

		// update cache if online version is newer
		if repomd.Revision > cache_repomd.Revision {
			update_mdcache = true
		} else {
			t.Logf("Repo metadata cache is still valid")
			update_mdcache = false
		}
	} else if err != nil && !os.IsNotExist(err) {
		t.Fatalf("Error reading precached repo metadata: %v", err)
	} else {
		t.Logf("Repo metadata is not yet cached")
		update_mdcache = true
	}

	// cache metadata locally
	if update_mdcache {
		t.Logf("Caching repo metadata...")
		if err = ioutil.WriteFile(cache_repomd_path, b, 0640); err != nil {
			t.Fatalf("Error caching repo metadata: %v", err)
		}
	}

	return repomd
}

// GetDatabase downloads a repository database to the given cache directory.
func GetDatabase(t *testing.T, baseurl string, cache_dir string, d *RepoDatabase) string {
	// append location to baseurl
	db_uri := fmt.Sprintf("%s%s", baseurl, d.Location.Href)

	// get file extension from db type
	db_ext := ""
	switch d.DatabaseVersion {
	case 0:
		db_ext = ".xml.gz"
		break

	case 10:
		db_ext = ".sqlite.bz2"
		break

	default:
		t.Fatalf("Unsupported database file version in %s: %d", d.Type, d.DatabaseVersion)
	}

	db_filename := fmt.Sprintf("%s%s", d.Type, db_ext)
	db_path := filepath.Join(cache_dir, db_filename)

	// should we download the file again?
	update_db := false

	// does file exist?
	f, err := os.Open(db_path)
	if err == nil {
		defer f.Close()

		err := d.Checksum.Check(f)
		if err == ErrChecksumMismatch {
			t.Logf("Checksum mismatch. Need to download again.")
			update_db = true
		} else if err != nil {
			t.Errorf("Error reading SHA256 checksum for file %s: %v", db_path, err)
			update_db = true
		} else {
			t.Logf("Checksum matches for precached file")
		}

	} else if os.IsNotExist(err) {
		update_db = true
	} else {
		t.Fatalf("Error opening cached file %s: %v", db_path, err)
	}

	if update_db {
		// download file
		t.Logf("Downloading file: %s", db_uri)
		if err = DownloadFile(db_uri, db_path); err != nil {
			t.Fatalf("%v", err)
		}
	}

	return db_path
}

func CheckChecksum(path string, checksum *RepoDatabaseChecksum) error {
	f, err := os.Open(path)
	if err == nil {
		defer f.Close()

		err := checksum.Check(f)
		if err != nil {
			return err
		}

		return nil
	}

	return err
}

func DecompressDatabase(t *testing.T, d *RepoDatabase, path string) string {
	// determine decompressed file path
	outpath := ""
	switch d.DatabaseVersion {
	case 0:
		outpath = path[:len(path)-3]
		break

	case 10:
		outpath = path[:len(path)-4]
		break

	default:
		t.Fatalf("Unsupported database file version in %s: %d", d.Type, d.DatabaseVersion)
	}

	// do we need to decompress the file?
	decompress := true

	// does the file exist and match checksum?
	if err := CheckChecksum(outpath, &d.OpenChecksum); err == nil {
		t.Logf("Database already extracted and valid")
		decompress = false
	}

	// decompress
	if decompress {
		switch d.DatabaseVersion {
		case 0:
			if err := GzipDecompress(path, outpath); err != nil {
				t.Fatalf("Error decompressing file %s: %v", path, err)
			}
			t.Logf("Extracted %s database with gzip", d.Type)
			break

		case 10:
			outpath := path[:len(path)-4]
			if err := Bzip2Decompress(path, outpath); err != nil {
				t.Fatalf("Error decompressing file %s: %v", path, err)
			}
			t.Logf("Extracted %s database with bzip2", d.Type)
			break

		}
	}

	return outpath
}

func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func DownloadFile(url string, path string) error {
	// get resource
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error retrieving %s: %v", url, err)
	}
	defer resp.Body.Close()

	// check response code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Bad response: %s", resp.Status)
	}

	// open file for writing
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Error creating file %s: %v", path, err)
	}
	defer f.Close()

	// write to file
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("Error writing file %s: %v", path, err)
	}

	return nil
}

func Bzip2Decompress(path string, out string) error {
	if !strings.HasSuffix(path, ".bz2") {
		return fmt.Errorf("File does not have the .bz2 extension: %s", path)
	}

	// open the file
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// open output file
	o, err := os.Create(out)
	if err != nil {
		return err
	}
	defer o.Close()

	// read the bzip2 file
	z := bzip2.NewReader(f)
	_, err = io.Copy(o, z)
	if err != nil {
		return err
	}

	return nil
}

func GzipDecompress(path string, out string) error {
	if !strings.HasSuffix(path, ".gz") {
		return fmt.Errorf("File does not have the .gz extension: %s", path)
	}

	// open the file
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// open output file
	o, err := os.Create(out)
	if err != nil {
		return err
	}
	defer o.Close()

	// read the gzip file
	z, err := gzip.NewReader(f)
	if err != nil {
		return err
	}

	_, err = io.Copy(o, z)
	if err != nil {
		return err
	}

	return nil
}
