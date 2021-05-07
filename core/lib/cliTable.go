package lib

import (
	"bufio"
	"fmt"
	"strings"
)

type CLITable struct {
	sizes      []int
	colNames   []string
	rows       [][]string
	currentRow int
	currentCol int
}

func NewCLITable(colNames []string) *CLITable {

	sizes := make([]int, len(colNames))
	rows := [][]string{
		[]string{},
	}
	for k := range colNames {
		sizes[k] = len(colNames[k])
		rows[0] = append(rows[0], colNames[k])
	}

	return &CLITable{
		sizes,
		colNames,
		rows,
		0,
		0,
	}
}

// row := t.NewRow()

// row[0] = col.Name
// row[1] = col.DataType
// row[2] = fmt.Sprintf("%d", col.MaxLength)

// if col.IsNullable {
// 	row[3] = "YES"
// } else {
// 	row[3] = "NO"
// }

// row[4] = col.Default
// row[5] = col.Extra

// t.AddRow(row)

func (c *CLITable) NewRow() []string {
	return make([]string, len(c.sizes))
}

func (c *CLITable) AddRow(row []string) {

	if len(row) != len(c.sizes) {
		panic(fmt.Sprintf("New column (%d) count does not match initial column count (%d)", len(row), len(c.sizes)))
	}

	c.rows = append(c.rows, row)

	for k := range row {
		if len(row[k]) > c.sizes[k] {
			c.sizes[k] = len(row[k])
		}
	}
}

func (c *CLITable) Row() {
	c.currentRow++
	c.currentCol = 0
	c.rows = append(c.rows, make([]string, len(c.sizes)))
}

func (c *CLITable) Col(val string) {
	c.rows[c.currentRow][c.currentCol] = val
	if len(val) > c.sizes[c.currentCol] {
		c.sizes[c.currentCol] = len(val)
	}
	c.currentCol++
}

func (c *CLITable) Colf(template string, args ...interface{}) {
	c.Col(fmt.Sprintf(template, args...))
}

func (c *CLITable) String() string {

	totalWidth := 0
	for k := range c.sizes {
		totalWidth += c.sizes[k] + 3
	}

	// Last pipe
	totalWidth++

	s := ""
	s += fmt.Sprintf("+" + strings.Repeat("-", totalWidth-2) + "+\n")
	for n := range c.rows {

		for col := range c.rows[n] {
			s += "| " + fmt.Sprintf(fmt.Sprintf("%%-%dv", c.sizes[col]), c.rows[n][col]) + " "
		}
		s += "|"
		s += fmt.Sprint("\n")

		if n == 0 {

			s += fmt.Sprintf("+" + strings.Repeat("-", totalWidth-2) + "+")

			s += fmt.Sprint("\n")
		}
	}
	s += fmt.Sprintf("+" + strings.Repeat("-", totalWidth-2) + "+")

	return s
}

func ReadCliInput(reader *bufio.Reader, title string) string {
	fmt.Print("> " + title)
	val, _ := reader.ReadString('\n')
	val = strings.Replace(val, "\n", "", -1)
	return val
}
