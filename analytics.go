package main

import "strings"

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
		
		// Check Poverty Status (Desil 1 or 2)
		isPoor := hh.WelfareLevel == "1" || hh.WelfareLevel == "2"
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
			// Job Profile of Head
			job := hh.JobProfile
			if job == "" { job = "Tidak Bekerja/Lainnya" }
			data.JobProfilePoor[job]++
			
			// Education of Head (Need to find head in members again or store it)
			// Simplification: Iterate members to find head
			foundHead := false
			for _, m := range hh.Members {
				if strings.Contains(strings.ToUpper(m.Relation), "KEPALA") {
					edu := m.Education
					if edu == "" { edu = "Tidak Sekolah/Belum Tamat SD" }
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
		
		// 2. Infrastructure
		// RTLH Logic: Status Sewa/Menumpang OR (Lantai Tanah/Kayu AND Dinding Bambu/Kayu)
		// For now, let's use the HouseStatus field we added.
		// Assume "Milik Sendiri" is good. Others are risk.
		isRTLH := false
		houseStatus := strings.ToUpper(hh.HouseStatus)
		if !strings.Contains(houseStatus, "MILIK SENDIRI") && houseStatus != "" {
			isRTLH = true
		}
		
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
		// Logic: Head > 65 AND (Member Count == 1 OR Member Count == 2)
		if hh.IsElderlyHead && len(hh.Members) <= 2 {
			data.ElderlySingleCount++
		}
	}
	
	return data
}
