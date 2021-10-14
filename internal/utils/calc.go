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
	d["00000001-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000001-0003-0000-0000-ff00ff00ff00", "тестостерон", "testosterone", "testosterone",
		7, "c19h28o2", 288.431, 288.431, 100, 100, 180,
	}

	d["00000002-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000002-0003-0000-0000-ff00ff00ff00", "тестостерона ацетат", "testosterone acetate", "testosterone",
		6, "c21h30o3", 330.468, 288.431, 87.2795617277104, 100, 724,
	}

	d["00000003-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000003-0003-0000-0000-ff00ff00ff00", "тестостерона пропионат", "testosterone propionate", "testosterone",
		6, "c22h32o3", 344.495, 288.431, 83.7257451979459, 100, 1152,
	}

	d["00000004-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000004-0003-0000-0000-ff00ff00ff00", "тестостерона ципионат", "testosterone cypionate", "testosterone",
		6, "c27h40o3", 412.614, 288.431, 69.9033482093565, 100, 7200,
	}

	d["00000005-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000005-0003-0000-0000-ff00ff00ff00", "тестостерона деканоат", "testosterone decanoate", "testosterone",
		6, "c29h46o3", 442.67, 288.431, 65.1571113972251, 100, 10800,
	}

	d["00000006-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000006-0003-0000-0000-ff00ff00ff00", "тестостерона энантат", "testosterone enanthate", "testosterone",
		6, "c26h40o3", 400.603, 288.431, 71.9992117395871, 100, 6480,
	}
	return d
}
