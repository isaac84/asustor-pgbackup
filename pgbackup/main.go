package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

func main() {
	var d = false
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to debug
	if ok {
		d = "debug" == lvl
	}

	file, err := os.Create("photogallery-tmp.db")

	if err != nil {
		log.Fatal(err.Error())
	}

	outlog, err := os.Create("backupfiles.txt")

	if err != nil {
		log.Fatal(err.Error())
	}

	in, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}
	defer in.Close()

	bytes, err := io.Copy(file, in)
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()

	fmt.Printf("Copied database %d bytes\n", bytes)

	sqliteDatabase, _ := sql.Open("sqlite3", "./photogallery-tmp.db")
	defer sqliteDatabase.Close()

	out := bufio.NewWriter(outlog)

	findPhotosWithTag(sqliteDatabase, os.Args[2], out, d)
	out.Flush()
}

func findPhotosWithTag(db *sql.DB, tag string, out *bufio.Writer, d bool) {
	log.Println("Finding tagged Photos")

	var count = 0
	taggedPhotoSQL := `select a.storage, a.inFile, a.inPar, a.name, a.type from PG_Media a, PG_Tag b, PG_Tag_List c where c.tag = ? and b.tag = c.id and a.inFile = b.inode`

	row, err := db.Query(taggedPhotoSQL, tag)
	if err != nil {
		log.Fatalln(err.Error())
	}

	defer row.Close()
	for row.Next() {
		var storage string
		var fileId int
		var parentId int
		var fileName string
		var fileType string

		row.Scan(&storage, &fileId, &parentId, &fileName, &fileType)

		if d {
			fmt.Printf("Photo: %s, %d, %d, %s, %s\n", storage, fileId, parentId, fileName, fileType)
		}

		var pathStr = getPath(db, storage, parentId)

		if d {
			fmt.Printf("Photo with path %s\n", pathStr+fileName)
		}

		out.WriteString(pathStr + fileName + "\n")
		count++
	}

	log.Println("Found ", count, " photos")
}

func getPath(db *sql.DB, storage string, parentId int) string {
	// extracted from download_photogallery.cgi
	pathSQL := `WITH RECURSIVE get_path(storage, inFile, inPar, name) AS (
			SELECT storage, inFile, inPar, name FROM PG_Folder WHERE storage = ? AND inFile = ?
				UNION ALL SELECT info.storage, info.inFile, info.inPar, info.name FROM PG_Folder AS info
				INNER JOIN get_path ON info.inFile = get_path.inPar AND info.storage = ? AND get_path.storage = ?
			) SELECT name FROM get_path;`

	path, err := db.Query(pathSQL, storage, parentId, storage, storage)
	if err != nil {
		log.Fatalln(err.Error())
	}

	var pathStr string
	defer path.Close()
	for path.Next() {
		var pathName string
		path.Scan(&pathName)

		pathStr = pathName + "/" + pathStr
	}

	pathStr = "/" + storage + "/" + pathStr
	return pathStr
}
