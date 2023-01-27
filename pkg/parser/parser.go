package parser

import "strings"

type PROCEDURE_COMMAND int64

const (
	INSERT PROCEDURE_COMMAND = iota
	FETCH
	DELETE
	UPDATE
	UNKNOWN
)

type Extract struct{}

var Extracter Extract = Extract{}

// parse commands
func Parse(command string) PROCEDURE_COMMAND {

	significant := strings.Split(command, " ")[0]

	switch significant {
	case "fetch":
		return FETCH
	case "insert":
		return INSERT
	case "delete":
		return DELETE
	case "update":
		return UPDATE
	}

	return UNKNOWN

}
