package gender_type_enums

var GenderType = newGenderType()

func newGenderType() *genderType {
	return &genderType{
		MALE:   "MALE",
		FEMALE: "FEMALE",
	}
}

type genderType struct {
	MALE   string
	FEMALE string
}
