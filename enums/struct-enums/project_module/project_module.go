package project_module

var ModuleNameEnums = newModule()

func newModule() *moduleNames {
	return &moduleNames{
		TEMPORARY_ATTACHMENTS: "Temporary Attachments",
		BASE_USER:             "User",
		ROLE:                  "Role",
	}
}

type moduleNames struct {
	TEMPORARY_ATTACHMENTS string
	BASE_USER             string
	ROLE                  string
}
