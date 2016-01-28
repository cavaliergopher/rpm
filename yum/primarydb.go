package yum

import (
	"database/sql"
	"fmt"
	"github.com/cavaliercoder/go-rpm"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
)

// TODO: Add support for XML primary dbs

// Queries to create primary_db schema
const (
	sqlCreateTables = `CREATE TABLE db_info (dbversion INTEGER, checksum TEXT);
CREATE TABLE packages ( pkgKey INTEGER PRIMARY KEY, pkgId TEXT, name TEXT, arch TEXT, version TEXT, epoch TEXT, release TEXT, summary TEXT, description TEXT, url TEXT, time_file INTEGER, time_build INTEGER, rpm_license TEXT, rpm_vendor TEXT, rpm_group TEXT, rpm_buildhost TEXT, rpm_sourcerpm TEXT, rpm_header_start INTEGER, rpm_header_end INTEGER, rpm_packager TEXT, size_package INTEGER, size_installed INTEGER, size_archive INTEGER, location_href TEXT, location_base TEXT, checksum_type TEXT);
CREATE TABLE files ( name TEXT, type TEXT, pkgKey INTEGER);
CREATE TABLE requires ( name TEXT, flags TEXT, epoch TEXT, version TEXT, release TEXT, pkgKey INTEGER , pre BOOLEAN DEFAULT FALSE);
CREATE TABLE provides ( name TEXT, flags TEXT, epoch TEXT, version TEXT, release TEXT, pkgKey INTEGER );
CREATE TABLE conflicts ( name TEXT, flags TEXT, epoch TEXT, version TEXT, release TEXT, pkgKey INTEGER );
CREATE TABLE obsoletes ( name TEXT, flags TEXT, epoch TEXT, version TEXT, release TEXT, pkgKey INTEGER );`

	sqlCreateTriggers = `CREATE TRIGGER removals AFTER DELETE ON packages  BEGIN    DELETE FROM files WHERE pkgKey = old.pkgKey;    DELETE FROM requires WHERE pkgKey = old.pkgKey;    DELETE FROM provides WHERE pkgKey = old.pkgKey;    DELETE FROM conflicts WHERE pkgKey = old.pkgKey;    DELETE FROM obsoletes WHERE pkgKey = old.pkgKey;  END;`

	sqlCreateIndexes = `CREATE INDEX packagename ON packages (name);
CREATE INDEX packageId ON packages (pkgId);
CREATE INDEX filenames ON files (name);
CREATE INDEX pkgfiles ON files (pkgKey);
CREATE INDEX pkgrequires on requires (pkgKey);
CREATE INDEX requiresname ON requires (name);
CREATE INDEX pkgprovides on provides (pkgKey);
CREATE INDEX providesname ON provides (name);
CREATE INDEX pkgconflicts on conflicts (pkgKey);
CREATE INDEX pkgobsoletes on obsoletes (pkgKey);`
)

// Queries to insert packages
const (
	sqlInsertPackage = `INSERT INTO packages (pkgId, name, arch, version, epoch, release, summary, description, url, time_file, time_build, rpm_license, rpm_vendor, rpm_group, rpm_buildhost, rpm_sourcerpm, rpm_header_start, rpm_header_end, rpm_packager, size_package, size_installed, size_archive, location_href, location_base, checksum_type) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10, ?11, ?12, ?13, ?14, ?15, ?16, ?17, ?18, ?19, ?20, ?21, ?22, ?23, ?24, ?25);`

	sqlInsertRequires = `INSERT INTO requires (pkgKey, name, flags, epoch, version, release, pre) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7);`
)

// queries to select packages
const sqlSelectPackages = `SELECT
 pkgKey
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

// PrimaryDatabase is an SQLite database which contains package data for a
// yum package repository.
type PrimaryDatabase struct {
	dbpath string
}

// CreatePrimaryDB initializes a new and empty primary_db SQLite database on
// disk.
func CreatePrimaryDB(path string) (*PrimaryDatabase, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("Error creating Primary DB: %v", err)
	}
	defer db.Close()

	// create database tables
	_, err = db.Exec(sqlCreateTables)
	if err != nil {
		os.Remove(path)
		return nil, fmt.Errorf("Error creating Primary DB tables: %v", err)
	}

	// create database indexes
	_, err = db.Exec(sqlCreateIndexes)
	if err != nil {
		os.Remove(path)
		return nil, fmt.Errorf("Error creating Primary DB indexes: %v", err)
	}

	// create database triggers
	_, err = db.Exec(sqlCreateTriggers)
	if err != nil {
		os.Remove(path)
		return nil, fmt.Errorf("Error creating Primary DB triggers: %v", err)
	}

	// insert db_info data
	_, err = db.Exec(`INSERT INTO db_info (dbversion) VALUES (10);`)
	if err != nil {
		os.Remove(path)
		return nil, fmt.Errorf("Error setting database version info: %v", err)
	}

	return &PrimaryDatabase{
		dbpath: path,
	}, nil
}

// OpenPrimaryDB opens a primary_db SQLite database from file and return a
// pointer to the resulting struct.
func OpenPrimaryDB(path string) (*PrimaryDatabase, error) {
	// open database file
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// TODO: Validate primary_db on open, maybe with the db_info table

	return &PrimaryDatabase{
		dbpath: path,
	}, nil
}

// Packages returns all packages listed in the primary_db.
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
		p := PackageEntry{
			db: c,
		}

		// scan the values into the slice
		if err = rows.Scan(&p.key, &p.name, &p.architecture, &p.epoch, &p.version, &p.release, &p.package_size, &p.install_size, &p.archive_size, &p.locationhref, &p.checksum, &p.checksum_type, &p.time_build); err != nil {
			return nil, fmt.Errorf("Error scanning packages: %v", err)
		}

		packages = append(packages, p)
	}

	return packages, nil
}

func (c *PrimaryDatabase) AddPackages(packages []rpm.Package) error {
	// open database file
	db, err := sql.Open("sqlite3", c.dbpath)
	if err != nil {
		return err
	}
	defer db.Close()

	// start driver transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// prepare statements
	stmtInsPackage, err := tx.Prepare(sqlInsertPackage)
	if err != nil {
		return err
	}
	defer stmtInsPackage.Close()

	stmtInsRequires, err := tx.Prepare(sqlInsertRequires)
	if err != nil {
		return err
	}
	defer stmtInsRequires.Close()

	// insert packages
	for _, p := range packages {
		// checksum package
		sum, err := p.Checksum()
		if err != nil {
			return fmt.Errorf("Error computing checksum for %s: %v", p, err)
		}

		// TODO: Compute relative path of primary_db packages

		// insert package table entry
		res, err := stmtInsPackage.Exec(sum, p.Name(), p.Architecture(), p.Version(), p.Epoch(), p.Release(), p.Summary(), p.Description(), p.URL(), p.FileTime().Unix(), p.BuildTime().Unix(), p.License(), p.Vendor(), strings.Join(p.Groups(), ","), p.BuildHost(), p.SourceRPM(), p.HeaderStart(), p.HeaderEnd(), p.Packager(), p.FileSize(), p.Size(), p.ArchiveSize(), p.Path(), nil, "sha256")
		if err != nil {
			return fmt.Errorf("Error inserting package %v: %v", p, err)
		}

		// get package key
		pkgKey, err := res.LastInsertId()
		if err != nil {
			return err
		}

		// insert requires table entries
		reqs := p.Requires()
		for _, req := range reqs {
			// nil out values by default
			var flags, epoch, version, release interface{}

			if req.Flags() != 0 {
				flags = req.Flags()
			}

			if req.Epoch() > 0 {
				epoch = req.Epoch()
			}

			if req.Version() != "" {
				version = req.Version()
			}

			if req.Release() != "" {
				release = req.Release()
			}

			// insert
			_, err := stmtInsRequires.Exec(pkgKey, req.Name(), flags, epoch, version, release, false)
			if err != nil {
				return fmt.Errorf("Error inserting dependency for %v: %v", p, err)
			}
		}
	}

	// commit
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// DependenciesByPackage returns all package dependencies of the given type for
// the given package key. The dependency type may be one of 'requires',
// 'provides', 'conflicts' or 'obsoletes'.
func (c *PrimaryDatabase) DependenciesByPackage(pkgKey int, typ string) (rpm.Dependencies, error) {
	q := fmt.Sprintf("SELECT name, flags, epoch, version, release FROM %s WHERE pkgKey = %d", typ, pkgKey)

	// open database file
	db, err := sql.Open("sqlite3", c.dbpath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// select packages
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// parse results
	deps := make(rpm.Dependencies, 0)
	for rows.Next() {
		var flgs, name, version, release string
		var epoch, iflgs int

		if err = rows.Scan(&name, &flgs, &epoch, &version, &release); err != nil {
			return nil, fmt.Errorf("Error reading dependencies: %v", err)
		}

		switch flgs {
		case "EQ":
			iflgs = rpm.DepFlagEqual

		case "LT":
			iflgs = rpm.DepFlagLesser

		case "LE":
			iflgs = rpm.DepFlagLesserOrEqual

		case "GE":
			iflgs = rpm.DepFlagGreaterOrEqual

		case "GT":
			iflgs = rpm.DepFlagGreater
		}

		deps = append(deps, rpm.NewDependency(iflgs, name, epoch, version, release))
	}

	return deps, nil
}

// FilesByPackage returns all known files included in the package of the given
// package key.
func (c *PrimaryDatabase) FilesByPackage(pkgKey int) ([]string, error) {
	q := fmt.Sprintf("SELECT name FROM files WHERE pkgKey = %d", pkgKey)

	// open database file
	db, err := sql.Open("sqlite3", c.dbpath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// select packages
	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// parse results
	files := make([]string, 0)
	for rows.Next() {
		var file string
		if err := rows.Scan(&file); err != nil {
			return nil, fmt.Errorf("Error reading files: %v", err)
		}

		files = append(files, file)
	}

	return files, nil
}
