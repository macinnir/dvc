package commands

// CommandSelect selects rows from the database
func (c *Cmd) CommandSelect(args []string) {
	if len(args) > 0 && args[0] == "help" {
		helpSelect()
		return
	}
	// database := c.loadDatabase()

	// reader := bufio.NewReader(os.Stdin)

	// tableName := ""

	// if len(args) > 0 {
	// 	tableName = args[0]
	// 	args = args[1:]
	// } else {
	// 	tableName = lib.ReadCliInput(reader, "Table Name:")
	// }

	// if _, ok := database.Tables[tableName]; !ok {
	// 	fmt.Printf("Unknown table `%s`\n", tableName)
	// 	return
	// }

	// query := fmt.Sprintf("SELECT * FROM `%s` LIMIT 100\n", tableName)

}
