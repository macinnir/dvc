package main

import (
// "log"
)

const (
	// LogTypeError is an error log
	LogTypeError = "error"
	// LogTypeInfo is an info level log
	LogTypeInfo = "info"
)

// ToLogString parses the file name into a human readable string for logging the action
// func (l *logger) ParseChangeLogToLogString(f ChangeFile) (logString string) {

// 	// 001_{action}_{target}
// 	fileNameParts := strings.Split(f.Name, "_")

// 	// action
// 	fileAction := fileNameParts[1]

// 	// target
// 	fileTarget := strings.Join(fileNameParts[2:], "_")

// 	// Remove the `.sql` extension
// 	fileTarget = fileTarget[0 : len(fileTarget)-4]
// 	switch fileAction {
// 	case "createTable":
// 		logString = fmt.Sprintf("Creating table `%s`", fileTarget)
// 	case "alterTable":
// 		fileTargetParts := strings.Split(fileTarget, "__")
// 		fileActionParts := strings.Split(fileTargetParts[1], "_")
// 		logString = fmt.Sprintf("Altering table `%s` - %s %s", fileTargetParts[0], fileActionParts[0], fileActionParts[1])
// 	case "dropTable":
// 		logString = fmt.Sprintf("Dropping table `%s`", fileTarget)
// 	case "createView":
// 		logString = fmt.Sprintf("Creating view `%s`", fileTarget)
// 	case "alterView":
// 		logString = fmt.Sprintf("Altering view `%s`", fileTarget)
// 	case "dropView":
// 		logString = fmt.Sprintf("Dropping view `%s`", fileTarget)
// 	case "insert":
// 		logString = fmt.Sprintf("Inserting data into `%s`", fileTarget)
// 	}

// 	return
// }

// func (l *logger) LogFatal(msg string) {
// 	l.logs = append(l.logs, DVCLog{LogMessage: msg, LogType: LogTypeError})
// 	d.WriteLogs()
// 	log.Fatal(msg)
// }

// func (l *logger) Log(msg string) {
// 	log.Println(msg)
// 	l.logs = append(l.logs, DVCLog{LogMessage: msg, LogType: LogTypeInfo})
// }

// func (l *Logger) WriteLogs() {
// 	if len(d.logs) > 0 {
// 		tx, e := d.conn.Begin()
// 		if e != nil {
// 			log.Fatal(e)
// 		}

// 		for _, l := range d.logs {
// 			stmt, e := tx.Prepare("INSERT INTO dvcLog(dateCreated, logType, logMessage, dvcRunId) VALUES (UNIX_TIMESTAMP(), ?, ?, ?)")
// 			if e != nil {
// 				log.Fatal(e)
// 			}

// 			_, e = stmt.Exec(l.LogType, l.LogMessage, d.runID)
// 			if e != nil {
// 				log.Fatal(e)
// 			}
// 		}

// 		tx.Commit()
// 	}
// }
