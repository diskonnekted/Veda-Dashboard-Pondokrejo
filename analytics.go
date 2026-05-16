package main

import "strings"

// Lookup: ID Kerja PBDT → Label
var jobLabel = map[string]string{
	"1":  "Petani/Pekebun",
	"2":  "Nelayan",
	"3":  "Peternak",
	"4":  "Industri Pengolahan",
	"5":  "Konstruksi/Bangunan",
	"6":  "Perdagangan",
	"7":  "Transportasi",
	"8":  "Karyawan Swasta",
	"9":  "TNI/Polri",
	"10": "PNS/ASN",
	"11": "Wiraswasta",
	"12": "Buruh Harian Lepas",
	"13": "Pembantu Rumah Tangga",
	"14": "Pedagang",
	"15": "Jasa Lainnya",
	"16": "Karyawan BUMN/Bank",
	"17": "Karyawan Swasta Lainnya",
	"18": "Jasa Profesional",
	"19": "PNS/ASN",
	"20": "Ibu Rumah Tangga",
	"21": "Pelajar/Mahasiswa",
	"22": "Pensiunan",
	"23": "Tidak Bekerja",
	"24": "Lainnya",
	"25": "Pekerja Jasa/Outsourcing",
	"26": "Wiraswasta Usaha Sewa",
}

// Lookup: ID Jenjang PBDT → Label
var eduLabel = map[string]string{
	"1":  "Tidak/Belum Sekolah",
	"2":  "Belum Tamat SD",
	"3":  "Tamat SD/Sederajat",
	"4":  "SMP/Sederajat",
	"5":  "SMA/Sederajat",
	"6":  "SMK/Sederajat",
	"7":  "SD/Sederajat",
	"8":  "SMP/MTs",
	"9":  "SMA/MA",
	"10": "Diploma/D3",
	"11": "S1/D4",
	"12": "S2",
	"13": "S3",
	"15": "Tidak Sekolah",
}

// AnalyticsData holds the aggregated data for the dashboard
type AnalyticsData struct {
	TotalHouseholds int `json:"total_households"`
	TotalResidents  int `json:"total_residents"`
	
	// 1. Kesejahteraan
	IncomeDistribution map[string]int `json:"income_distribution"`
	JobProfilePoor     map[string]int `json:"job_profile_poor"` // Jobs of Desil 1 & 2
	EducationPoor      map[string]int `json:"education_poor"`     // Education of Desil 1 & 2
	IncomeComparison   map[string]int `json:"income_comparison"`  // New: Miskin vs Mampu Income
	
	// 2. Infrastruktur
	RTLHCount           int `json:"rtlh_count"`            // Rumah Tidak Layak Huni (Sewa/Numpang + Dinding/Lantai Buruk)
	NoLatrineCount      int `json:"no_latrine_count"`      // Tidak punya jamban
	NoCleanWaterCount   int `json:"no_clean_water_count"`  // Air sungai/danau
	
	// 3. Rentan
	ElderlySingleCount  int `json:"elderly_single_count"`  // Lansia Tunggal
	PoorWithToddlerCount int `json:"poor_with_toddler_count"` // Miskin + Balita
	StuntingCount        int `json:"stunting_count"`          // Resiko Stunting
}

func CalculateAnalytics(households []Household) AnalyticsData {
	data := AnalyticsData{
		TotalHouseholds: len(households),
		IncomeDistribution: make(map[string]int),
		JobProfilePoor:     make(map[string]int),
		EducationPoor:      make(map[string]int),
		IncomeComparison:   make(map[string]int),
	}
	
	for _, hh := range households {
		data.TotalResidents += len(hh.Members)
		
		// 1. Income Distribution
		if hh.IncomeCategory == "" { hh.IncomeCategory = "Tidak Diketahui" }
		data.IncomeDistribution[hh.IncomeCategory]++
		
		// Check Poverty Status (Desil 1, 2, or 3)
		isPoor := hh.WelfareLevel == "1" || hh.WelfareLevel == "2" || hh.WelfareLevel == "3"
		isMampu := hh.WelfareLevel == "4"

		// New: Income Comparison (Miskin vs Mampu)
		// We use IncomeCategory: <1jt, 1-2jt, >2jt
		// We want to count how many Miskin are in each category vs Mampu
		if isPoor {
			data.IncomeComparison["Miskin - "+hh.IncomeCategory]++
		} else if isMampu {
			data.IncomeComparison["Mampu - "+hh.IncomeCategory]++
		}
		
		if isPoor {
			// Job Profile of Head — translate kode numerik ke label
			jobCode := hh.JobProfile
			job := jobLabel[jobCode]
			if job == "" {
				if jobCode == "" {
					job = "Tidak Bekerja/Lainnya"
				} else {
					job = "Kode " + jobCode
				}
			}
			data.JobProfilePoor[job]++

			// Education of Head — translate kode numerik ke label
			foundHead := false
			for _, m := range hh.Members {
				// Kepala = ID Hub Keluarga = 1
				if strings.TrimSpace(m.Relation) == "1" {
					eduCode := m.Education
					edu := eduLabel[eduCode]
					if edu == "" { edu = "Tidak Diketahui" }
					data.EducationPoor[edu]++
					foundHead = true
					break
				}
			}
			if !foundHead {
				data.EducationPoor["Tidak Diketahui"]++
			}

			// 3. Rentan: Poor with Toddler
			if hh.HasToddler {
				data.PoorWithToddlerCount++
			}
		}

		// Stunting count (now calculated in parser, so we can just use the flag)
		if hh.IsStun {
			data.StuntingCount++
		}
		
		// 2. Infrastructure
		// RTLH: Parser sudah set HouseStatus="RTLH" jika dinding/atap buruk
		isRTLH := hh.HouseStatus == "RTLH"
		if isRTLH {
			data.RTLHCount++
		}

		if !hh.HasLatrine {
			data.NoLatrineCount++
		}
		if !hh.HasCleanWater {
			data.NoCleanWaterCount++
		}

		// 3. Rentan: Elderly Single
		if hh.IsElderlyHead && len(hh.Members) <= 2 {
			data.ElderlySingleCount++
		}
	}
	
	return data
}
