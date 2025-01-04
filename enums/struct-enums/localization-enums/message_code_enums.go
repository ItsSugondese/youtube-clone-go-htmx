package localization_enums

var MessageCodeEnums = newMessageCode()

func newMessageCode() *messageCode {
	return &messageCode{
		SAVE:                  "save",
		API_OPERATION:         "api.operation",
		COLUMN_NOT_EXISTS:     "column.not.exist",
		COLUMN_ALREADY_EXISTS: "column.already.exist",
	}
}

type messageCode struct {
	SAVE                  string
	API_OPERATION         string
	COLUMN_NOT_EXISTS     string
	COLUMN_ALREADY_EXISTS string
}
