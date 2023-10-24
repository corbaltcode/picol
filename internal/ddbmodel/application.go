package ddbmodel

type Application int8

const (
	ApplicationAerial        Application = 1
	ApplicationGround        Application = 2
	ApplicationIrrigation    Application = 3
	ApplicationPlantDip      Application = 4
	ApplicationSeedTreatment Application = 5
)

func (a Application) Code() byte {
	switch a {
	case ApplicationAerial:
		return 'A'
	case ApplicationGround:
		return 'G'
	case ApplicationIrrigation:
		return 'I'
	case ApplicationPlantDip:
		return 'D'
	case ApplicationSeedTreatment:
		return 'S'
	}
	panic("unknown application")
}

func (a Application) Name() string {
	switch a {
	case ApplicationAerial:
		return "AERIAL"
	case ApplicationGround:
		return "GROUND"
	case ApplicationIrrigation:
		return "IRRIGATION"
	case ApplicationPlantDip:
		return "PLANT DIP"
	case ApplicationSeedTreatment:
		return "SEED TREATMENT"
	}
	panic("unknown application")
}

func (a Application) String() string {
	switch a {
	case ApplicationAerial:
		return "Aerial"
	case ApplicationGround:
		return "Ground"
	case ApplicationIrrigation:
		return "Irrigation"
	case ApplicationPlantDip:
		return "Plant Dip"
	case ApplicationSeedTreatment:
		return "Seed Treatment"
	}
	panic("unknown application")
}
