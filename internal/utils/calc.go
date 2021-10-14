package utils

type Drug struct {
	ID       string
	name     string
	drug     string
	mdrug    string
	type_    int
	formula  string
	mmass    float64
	mmassm   float64
	Out      float64
	outt     int
	Halflife int
}

func GetDrugs() map[string]Drug {
	d := make(map[string]Drug)

	d["00000010-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000010-0003-0000-0000-ff00ff00ff00", "туринабол", "chlorodehydromethyltestosterone", "testosterone",
		8, "c20h27clo2", 334.884, 288.431, 86.1286283213281, 100, 960,
	}
	d["00000011-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000011-0003-0000-0000-ff00ff00ff00", "метилтестостерон", "methyltestosterone", "testosterone",
		8, "c20h30o2", 302.458, 288.431, 95.3623255625719, 100, 420,
	}

	return d
}
