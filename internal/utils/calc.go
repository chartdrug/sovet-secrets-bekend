package utils

import (
	"github.com/qiangxue/sovet-secrets-bekend/internal/entity"
	"math"
)

type Drug struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Drug     string  `json:"drug"`
	Mdrug    string  `json:"mdrug"`
	Type_    int     `json:"type_"`
	Formula  string  `json:"formula"`
	Mmass    float64 `json:"mmass"`
	Mmassm   float64 `json:"mmassm"`
	Out      float64 `json:"out"`
	Outt     float64 `json:"outt"`
	Halflife int     `json:"halflife"`
}

type InjBall struct {
	what     string
	drugId   string
	dose     float64
	volume   float64
	halflife int
	outK     float64
	outKT    float64
	CO       float64
	COT      float64
	R        float64
	skin     float64
	pending  bool
}

type BloodVolume struct {
	Dt int64
	V  float64
}

const ZERO = 1e-6 //что считать нулем

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

	d["00000007-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000007-0003-0000-0000-ff00ff00ff00", "тестостерона изакапроат", "testosterone isocaproate", "testosterone",
		6, "c25h38o3", 386.576, 288.431, 74.6117234023326, 100, 5760,
	}

	d["00000008-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000008-0003-0000-0000-ff00ff00ff00", "тестостерона фенилпропионат", "testosterone phenylpropionate", "testosterone",
		6, "c28h36o3", 420.593, 288.431, 68.5772230412252, 100, 2160,
	}

	d["00000009-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000009-0003-0000-0000-ff00ff00ff00", "тестостерона ундеканоат", "testosterone undecanoate", "testosterone",
		6, "c30h48o3", 456.711, 288.431, 63.1539418589191, 100, 30096,
	}

	d["00000012-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000012-0003-0000-0000-ff00ff00ff00", "халотестин", "fluoxymesterone", "fluoxymesterone",
		8, "c20h29fo3", 336.447, 336.447, 100, 85, 420,
	}

	d["00000013-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000013-0003-0000-0000-ff00ff00ff00", "провирон", "mesterolone", "mesterolone",
		8, "c20h32o2", 304.474, 304.474, 100, 35, 750,
	}

	d["00000014-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000014-0003-0000-0000-ff00ff00ff00", "станозолол", "stanozolol", "stanozolol",
		9, "c21h32n2o", 328.5, 328.5, 100, 30, 510,
	}

	d["00000015-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000015-0003-0000-0000-ff00ff00ff00", "нандролона деканоат", "nandrolone decanoate", "nandrolone",
		6, "c28h44o3", 428.657, 274.404, 64.0148129980334, 37, 10800,
	}

	d["00000016-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000016-0003-0000-0000-ff00ff00ff00", "нандролона фенилпропионат", "nandrolone phenpropionate", "nandrolone",
		6, "c27h34o3", 406.566, 274.404, 67.4930995158196, 37, 2160,
	}

	d["00000017-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000017-0003-0000-0000-ff00ff00ff00", "нандролона ацетат", "nandrolone acetate", "nandrolone",
		6, "c20h28o3", 316.441, 274.404, 86.7156820787761, 37, 724,
	}

	d["00000018-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000018-0003-0000-0000-ff00ff00ff00", "тренболон", "trenbolone", "trenbolone",
		6, "c18h22o2", 270.372, 270.372, 100, 50, 420,
	}

	d["00000019-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000019-0003-0000-0000-ff00ff00ff00", "тренболона ацетат", "trenbolone acetate", "trenbolone",
		6, "c20h24o3", 312.409, 270.372, 86.5442431564281, 50, 724,
	}

	d["00000020-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000020-0003-0000-0000-ff00ff00ff00", "тренболона энантат", "trenbolone enanthate", "trenbolone",
		6, "c25h34o3", 382.544, 270.372, 70.677364114122, 50, 6480,
	}

	d["00000021-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000021-0003-0000-0000-ff00ff00ff00", "тренболона циклогексаметилкарбонат", "trenbolone cyclohexylmethylcarbonate", "trenbolone",
		6, "c26h34o4", 410.554, 270.372, 65.8554090236725, 50, 11520,
	}

	d["00000022-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000022-0003-0000-0000-ff00ff00ff00", "болденона ундецикленат", "boldenone undecylenate", "boldenone",
		6, "c30h44o3", 452.679, 286.415, 63.2711069235131, 50, 17280,
	}

	d["00000023-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000023-0003-0000-0000-ff00ff00ff00", "дростанолона пропионат", "drostanolone propionate", "drostanolone",
		6, "c23h36o3", 360.538, 304.474, 84.4499098229328, 35, 1152,
	}

	d["00000024-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000024-0003-0000-0000-ff00ff00ff00", "дростанолона энантат", "drostanolone enanthate", "drostanolone",
		6, "c27h44o3", 416.646, 304.474, 73.0773864618883, 35, 6480,
	}

	d["00000025-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000025-0003-0000-0000-ff00ff00ff00", "метенолона энантат", "methenolone enanthate", "methenolone",
		6, "c27h42o3", 414.63, 302.458, 72.9464852655823, 50, 6480,
	}
	d["00000026-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000026-0003-0000-0000-ff00ff00ff00", "метенолона ацетат", "methenolone acetate", "methenolone",
		14, "c22h32o3", 344.495, 302.458, 87.7975026414269, 50, 724,
	}

	d["00000027-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000027-0003-0000-0000-ff00ff00ff00", "оксиметолон", "oxymetholone", "oxymetholone",
		14, "c21h32o3", 332.484, 332.484, 100, 45, 840,
	}

	d["00000028-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000028-0003-0000-0000-ff00ff00ff00", "оксандролон", "oxandrolone", "oxandrolone",
		8, "c19h30o3", 306.446, 306.446, 100, 24, 600,
	}

	d["00000029-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000029-0003-0000-0000-ff00ff00ff00", "метандростенолон", "methandienone", "methandienone",
		14, "c20h28o2", 300.442, 300.442, 100, 50, 300,
	}

	d["00000030-0003-0000-0000-ff00ff00ff00"] = Drug{
		"00000030-0003-0000-0000-ff00ff00ff00", "пустышка", "placebo", "placebo",
		7, "", 0, 0, 0, 0, 0,
	}

	return d
}

func SkinStep(s string) float64 {

	// ключ - "растворитель - инъекция"

	d := make(map[string]float64)
	d["W1"] = 0.0009
	d["W2"] = 0.0003
	d["O1"] = 0.0000025
	d["O2"] = 0.000000833
	d["L1"] = 0.000018
	d["L2"] = 0.000006

	return d[s]
}

func (InjBall *InjBall) ballOut() (int64, int64) {

	var r = InjBall.R

	return 0, 0

	if r < ZERO {
		return 0, 0
	}

	var ri = r - InjBall.skin

	if ri < 0 {
		ri = 0
	}

	var depo = InjBall.dose

	var depoi = depo * ((4 / 3 * math.Pi * math.Pow(ri, 3)) / InjBall.volume) //считаем новый объем

	var dv = depo - depoi
	InjBall.R = ri
	if depoi < ZERO {
		depoi = 0
	}

	InjBall.dose = depoi
	if dv < ZERO {
		return 0, 0
	}

	InjBall.pending = true

	return int64(dv * InjBall.outK), int64(dv * InjBall.outKT)

}

func GetBloodVolume(sex string, antros []entity.Antro) []BloodVolume {
	var b []BloodVolume

	for _, antro := range antros {
		var V = 0.0
		if sex == "M" {
			V = math.Pow(antro.Result_fat, -0.114) * antro.General_weight * 98
		} else if sex == "F" {
			V = math.Pow(antro.Result_fat, -0.157) * antro.General_weight * 106
		} else {
			V = 0.0
		}
		b = append(b, BloodVolume{Dt: antro.Dt.Unix() * 1000, V: V})
	}
	return b
}
