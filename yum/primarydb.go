package yum

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

// TODO: Add support for XML primary dbs

const sqlCreateTables = `CREATE TABLE db_info (dbversion INTEGER, checksum TEXT);
CREATE TABLE packages ( pkgKey INTEGER PRIMARY KEY, pkgId TEXT, name TEXT, arch TEXT, version TEXT, epoch TEXT, release TEXT, summary TEXT, description TEXT, url TEXT, time_file INTEGER, time_build INTEGER, rpm_license TEXT, rpm_vendor TEXT, rpm_group TEXT, rpm_buildhost TEXT, rpm_sourcerpm TEXT, rpm_header_start INTEGER, rpm_header_end INTEGER, rpm_packager TEXT, size_package INTEGER, size_installed INTEGER, size_archive INTEGER, location_href TEXT, location_base TEXT, checksum_type TEXT);
CREATE TABLE files ( name TEXT, type TEXT, pkgKey INTEGER);
CREATE TABLE requires ( name TEXT, flags TEXT, epoch TEXT, version TEXT, release TEXT, pkgKey INTEGER , pre BOOLEAN DEFAULT FALSE);
CREATE TABLE provides ( name TEXT, flags TEXT, epoch TEXT, version TEXT, release TEXT, pkgKey INTEGER );
CREATE TABLE conflicts ( name TEXT, flags TEXT, epoch TEXT, version TEXT, release TEXT, pkgKey INTEGER );
CREATE TABLE obsoletes ( name TEXT, flags TEXT, epoch TEXT, version TEXT, release TEXT, pkgKey INTEGER );`

const sqlCreateTriggers = `CREATE TRIGGER removals AFTER DELETE ON packages  BEGIN    DELETE FROM files WHERE pkgKey = old.pkgKey;    DELETE FROM requires WHERE pkgKey = old.pkgKey;    DELETE FROM provides WHERE pkgKey = old.pkgKey;    DELETE FROM conflicts WHERE pkgKey = old.pkgKey;    DELETE FROM obsoletes WHERE pkgKey = old.pkgKey;  END;`

const sqlCreateIndexes = `CREATE INDEX packagename ON packages (name);
CREATE INDEX packageId ON packages (pkgId);
CREATE INDEX filenames ON files (name);
CREATE INDEX pkgfiles ON files (pkgKey);
CREATE INDEX pkgrequires on requires (pkgKey);
CREATE INDEX requiresname ON requires (name);
CREATE INDEX pkgprovides on provides (pkgKey);
CREATE INDEX providesname ON provides (name);
CREATE INDEX pkgconflicts on conflicts (pkgKey);
CREATE INDEX pkgobsoletes on obsoletes (pkgKey);`

const sqlSelectPackages = `SELECT
 pkgKey,
 , name
 , arch
 , epoch
 , version
 , release
 , size_package
 , size_installed
 , size_archive
 , location_href
 , pkgId
 , checksum_type
 , time_build
FROM packages;`

type PrimaryDatabase struct {
	dbpath string
}

func CreatePrimaryDB(path string) error {
	// create database file
	dbpath := "./primary_db.sqlite"
	os.Remove(dbpath)

	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return fmt.Errorf("Error creating Primary DB: %v", err)
	}
	defer db.Close()

	// create database tables
	_, err = db.Exec(sqlCreateTables)
	if err != nil {
		return fmt.Errorf("Error creating Primary DB tables: %v", err)
	}

	// create database indexes
	_, err = db.Exec(sqlCreateIndexes)
	if err != nil {
		return fmt.Errorf("Error creating Primary DB indexes: %v", err)
	}

	// create database triggers
	_, err = db.Exec(sqlCreateTriggers)
	if err != nil {
		return fmt.Errorf("Error creating Primary DB triggers: %v", err)
	}

	return nil
}

func OpenPrimaryDB(path string) (*PrimaryDatabase, error) {
	// open database file
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// TODO: Validate primary_db on open

	return &PrimaryDatabase{
		dbpath: path,
	}, nil
}

func (c *PrimaryDatabase) Packages() (PackageEntries, error) {
	// open database file
	db, err := sql.Open("sqlite3", c.dbpath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// select packages
	rows, err := db.Query(sqlSelectPackages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// parse each row as a package
	packages := make(PackageEntries, 0)
	for rows.Next() {
		p := PackageEntry{}

		// scan the values into the slice
		if err = rows.Scan(&p.key, &p.name, &p.architecture, &p.epoch, &p.version, &p.release, &p.package_size, &p.install_size, &p.archive_size, &p.locationhref, &p.checksum, &p.checksum_type, &p.time_build); err != nil {
			return nil, fmt.Errorf("Error scanning packages: %v", err)
		}

		packages = append(packages, p)
	}

	return packages, nil
}
