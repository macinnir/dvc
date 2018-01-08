package dvc

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"sort"
)

// DVC data version control
type DVC struct {
	Connections map[string]Server
	ChangeLogs  map[string]ChangeLog
	Databases   map[string]Database
}

// Run runs the dvc
func (d *DVC) Run(changesetsPath string, dbHost string, dbName string, dbUser string, dbPass string) {

	var e error
	var server Server
	var localChangeLogs LocalChangeLogs

	// 1. LocalChangeLogs
	log.Println("1. LocalChangeLogs...")
	localChangeLogs, e = NewLocalChangeLogs(changesetsPath)
	if e != nil {
		log.Fatal("Local changelog error ", e)
	}
	e = localChangeLogs.FetchLocalChangesetFiles()
	if e != nil {
		log.Fatal("FetchLocalChangesetFiles", e)
	}
	log.Printf("\t Found %d changesets\n", len(localChangeLogs.ChangeSets))

	// for changesetName, changeset := range localChangeLogs.ChangeSets {
	// 	log.Printf("\t Changeset Dir: %s\n", changesetName)

	// 	for _, k := range changeset.Files {
	// 		log.Printf("\t\t Change File: %s\n", k)
	// 	}
	// }

	// 2. Server
	log.Println("2. Server")
	server = NewServer(dbHost, dbUser, dbPass)
	server.FetchDatabases()
	log.Printf("\t Found %d databases\n", len(server.Databases))

	// 3. Database
	log.Printf("3. Database `%s`\n", dbName)
	database, ok := server.Databases[dbName]
	if !ok {
		log.Printf("Database `%s` does not exist. Trying to create...", dbName)
		database, e = server.CreateDatabase(dbName)
		if e != nil {
			log.Fatal(e)
		}
	}

	e = database.Connect(dbUser, dbPass)
	if e != nil {
		log.Fatal(e)
	}

	// e = database.DropBaseTables()
	// if e != nil {
	// 	log.Fatal(e)
	// }

	e = database.FetchTables()
	if e != nil {
		database.logFatal("FetchTables(): " + e.Error())
	}

	log.Println("3.1 Database: Create Version Tables...")
	e = database.CreateBaseTablesIfNotExists()
	if e != nil {
		database.logFatal("CreateBaseTablesIfNotExists " + e.Error())
	}

	e = database.Start()
	log.Printf("Starting run #%d", database.runID)
	if e != nil {
		database.logFatal("Start()" + e.Error())
	}

	// 4. Database Changesets
	log.Printf("3.2 Fetching changes in `%s`\n", dbName)
	e = database.FetchSets()
	if e != nil {
		database.logFatal("FetchSets()" + e.Error())
	}

	// Resolve
	log.Println("3.3 Importing files...")

	// Delete sets not found
	for _, changesetName := range database.sortedSetKeys {

		dbChangeset := database.sets[changesetName]

		if _, ok := localChangeLogs.ChangeSets[changesetName]; !ok {
			database.DeleteSet(changesetName)
			continue
		}

		localChangeSet := localChangeLogs.ChangeSets[changesetName]

		for _, fileName := range dbChangeset.SortedFileKeys {

			if len(fileName) == 0 {
				continue
			}

			// fmt.Printf("Looking for %s\n", fileName)
			dbFile, _ := dbChangeset.Files[fileName]

			found := false

			for _, localFileName := range localChangeSet.Files {
				if localFileName == fileName {
					// found
					found = true
					break
				}
			}
			if found == false {
				database.log(fmt.Sprintf("File %s not found. Deleting...", fileName))
				database.DeleteFile(dbChangeset, dbFile)
			}
		}
	}

	for _, changesetName := range localChangeLogs.SortedChangesetKeys {

		changeset := localChangeLogs.ChangeSets[changesetName]

		if !database.SetExists(changesetName) {
			// log.Printf("Set Not Exist: %s", changesetName)
			if e = database.CreateSet(changesetName); e != nil {
				database.logFatal("CreateSet() " + e.Error())
			}
		}

		dbChangeset := database.sets[changesetName]

		for _, fileName := range changeset.Files {

			// fmt.Printf("File: %s\nOrdinal: %s\nAction: %s\nTarget: %s\n", fileName, fileOrdinal, fileAction, fileTarget)

			filePath := path.Join(changesetsPath, changesetName, fileName)
			if !database.FileExists(changesetName, fileName) {
				hash, _ := HashFileMd5(filePath)
				// log.Printf("File Not Exists: %s/%s", changesetName, fileName)

				var contentBytes []byte
				contentBytes, e = ioutil.ReadFile(filePath)
				if e != nil {
					database.logFatal("ReadFile() " + filePath + " " + e.Error())
				}

				content := string(contentBytes)

				_, e = database.CreateFile(dbChangeset, fileName, hash, content)
				if e != nil {
					database.logFatal("CreateFile() " + e.Error())
				}

			}
		}
	}

	log.Println("3.4 Applying database changes...")

	notRunChangeFiles := map[string]*DVCFile{}
	notRunChangeFileNames := []string{}

	for _, changesetName := range database.sortedSetKeys {

		for _, fileName := range database.sets[changesetName].SortedFileKeys {
			f := database.sets[changesetName].Files[fileName]
			if !f.IsRun {
				notRunChangeFileNames = append(notRunChangeFileNames, f.FullPath)
				notRunChangeFiles[f.FullPath] = f
			}
		}
	}

	sort.Strings(notRunChangeFileNames)

	database.log(fmt.Sprintf("Found %d changes to run \n", len(notRunChangeFiles)))

	for _, notRunChangeFileName := range notRunChangeFileNames {
		// fmt.Printf("\t\t %s\n", notRunChangeFileName)
		fileBytes, readFileErr := ioutil.ReadFile(path.Join(changesetsPath, notRunChangeFileName))
		if readFileErr != nil {
			log.Fatal(readFileErr)
		}
		fileString := string(fileBytes)
		e = database.RunChange(fileString)
		if e != nil {
			database.logFatal("RunChange() File: " + notRunChangeFileName + " - " + e.Error())
		}

		database.log(notRunChangeFiles[notRunChangeFileName].ToLogString())

		e = database.SetFileAsRun(notRunChangeFiles[notRunChangeFileName].ID)
		if e != nil {
			log.Fatal(e)
		}
	}
	database.finish()
}
