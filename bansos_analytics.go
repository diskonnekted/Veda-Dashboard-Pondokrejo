package main

import "sort"

// BansosAnalytics holds aggregated data for the social protection dashboard
type BansosAnalytics struct {
	TotalHouseholds int `json:"total_households"`

	// Totals per program
	TotalBpnt    int `json:"total_bpnt"`
	TotalPkh     int `json:"total_pkh"`
	TotalBlt     int `json:"total_blt"`
	TotalListrik int `json:"total_listrik"`
	TotalBaznas  int `json:"total_baznas"`
	TotalLpg     int `json:"total_lpg"`
	TotalBanpem  int `json:"total_banpem"`
	TotalEkstrem int `json:"total_ekstrem"`
	TotalStun    int `json:"total_stun"`

	// Multi-program overlap
	PenerimaTiga int `json:"penerima_tiga"`  // Dapat 3+ program
	PenerimaDua  int `json:"penerima_dua"`   // Dapat 2 program
	PenerimaZero int `json:"penerima_zero"`  // Tidak dapat bantuan sama sekali

	// Per Dusun breakdown
	DusunBpnt    map[string]int `json:"dusun_bpnt"`
	DusunPkh     map[string]int `json:"dusun_pkh"`
	DusunEkstrem map[string]int `json:"dusun_ekstrem"`
	DusunStun    map[string]int `json:"dusun_stun"`

	// Top Dusun by recipient count (for charts)
	DusunNames []string `json:"dusun_names"` // sorted

	// Coverage rate (persen penerima dari total KK)
	BpntCoverage float64 `json:"bpnt_coverage"`
	PkhCoverage  float64 `json:"pkh_coverage"`
	BltCoverage  float64 `json:"blt_coverage"`
}

func CalculateBansos(households []Household) BansosAnalytics {
	data := BansosAnalytics{
		TotalHouseholds: len(households),
		DusunBpnt:       make(map[string]int),
		DusunPkh:        make(map[string]int),
		DusunEkstrem:    make(map[string]int),
		DusunStun:       make(map[string]int),
	}

	dusunSet := make(map[string]bool)

	for _, hh := range households {
		dusunSet[hh.Dusun] = true

		// Count programs per KK
		programs := 0
		if hh.IsBpnt {
			data.TotalBpnt++
			data.DusunBpnt[hh.Dusun]++
			programs++
		}
		if hh.IsPkh {
			data.TotalPkh++
			data.DusunPkh[hh.Dusun]++
			programs++
		}
		if hh.IsBlt {
			data.TotalBlt++
			programs++
		}
		if hh.IsListrik {
			data.TotalListrik++
		}
		if hh.IsBaznas {
			data.TotalBaznas++
		}
		if hh.IsLpg {
			data.TotalLpg++
		}
		if hh.IsBanpem {
			data.TotalBanpem++
		}
		if hh.IsEkstrem {
			data.TotalEkstrem++
			data.DusunEkstrem[hh.Dusun]++
		}
		if hh.IsStun {
			data.TotalStun++
			data.DusunStun[hh.Dusun]++
		}

		if programs == 0 {
			data.PenerimaZero++
		} else if programs == 2 {
			data.PenerimaDua++
		} else if programs >= 3 {
			data.PenerimaTiga++
		}
	}

	// Coverage rates
	if data.TotalHouseholds > 0 {
		data.BpntCoverage = float64(data.TotalBpnt) / float64(data.TotalHouseholds) * 100
		data.PkhCoverage = float64(data.TotalPkh) / float64(data.TotalHouseholds) * 100
		data.BltCoverage = float64(data.TotalBlt) / float64(data.TotalHouseholds) * 100
	}

	// Collect and sort dusun names
	for d := range dusunSet {
		data.DusunNames = append(data.DusunNames, d)
	}
	sort.Strings(data.DusunNames)

	return data
}
