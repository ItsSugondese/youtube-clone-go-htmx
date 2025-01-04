package response_status_enum

import (
	"encoding/json"
	"fmt"
)

type ResponseStatusEnum interface {
	fmt.Stringer
	private()
}

type responseStatus struct {
	Str string
}

func (rs responseStatus) String() string {
	return rs.Str
}

// MarshalJSON method to handle custom marshaling to JSON
func (rs responseStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(rs.Str)
}

func (rs responseStatus) private() {
}

func Success() ResponseStatusEnum {
	return responseStatus{"SUCCESS"}
}

func Fail() ResponseStatusEnum {
	return responseStatus{"FAIL"}
}
