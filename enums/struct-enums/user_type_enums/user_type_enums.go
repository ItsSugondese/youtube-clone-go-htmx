package user_type_enums

var UserType = newUserType()

func newUserType() *userType {
	return &userType{
		CUSTOMER: "CUSTOMER",
		ADMIN: "ADMIN",

	}
}

type userType struct {
	CUSTOMER string
	ADMIN    string
}
