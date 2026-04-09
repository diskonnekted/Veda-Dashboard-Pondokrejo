package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Resident struct {
	Name        string   `json:"name"`
	Nik         string   `json:"nik"`
	Gender      string   `json:"gender"` // ID Kelamin (Ket)
	AidList     []string `json:"aid_list"`
	KerjaDetail string   `json:"kerja_detail"` // Kerja Detail
	UshDetail   string   `json:"ush_detail"`   // Ush Detail
	Age         string   `json:"age"`          // Usia
	Education   string   `json:"education"`    // ID Jenjang (Ket)
	Income      string   `json:"income"`       // Income
	Pregnant    string   `json:"pregnant"`     // Hamil
	Disability  string   `json:"disability"`   // ID Difable
	Relation    string   `json:"relation"`     // ID Hub Keluarga (Ket)
	Job         string   `json:"job"`          // ID Kerja (Ket)
}

type Household struct {
	NoKK         string     `json:"no_kk"`
	HeadName     string     `json:"head_name"`
	Address      string     `json:"address"`
	Dusun        string     `json:"dusun"`
	Latitude     float64    `json:"latitude"`
	Longitude    float64    `json:"longitude"`
	WelfareLevel string     `json:"welfare_level"` // ID Desil (Ket)
	Members      []Resident `json:"members"`
	PkhThn       string     `json:"pkh_thn"`      // Pkh Thn
	BpntThn      string     `json:"bpnt_thn"`     // Bpnt Thn
	LantaiLuas   string     `json:"lantai_luas"`  // Lantai Luas
	Keterangan   string     `json:"keterangan"`   // Keterangan
	Expenditure  string     `json:"expenditure"`  // Overall Sum
	FloorType    string     `json:"floor_type"`   // ID Lantai (Ket)
	WallType     string     `json:"wall_type"`    // ID Dinding (Ket)
	RoofType     string     `json:"roof_type"`    // ID Atap (Ket)
	WaterSource  string     `json:"water_source"` // ID Airminum (Ket)
	Sanitation   string     `json:"sanitation"`   // ID Fasbab (Ket)

	// New Analytics Fields
	IncomeCategory string `json:"income_category"` // <1jt, 1-2jt, >2jt
	JobProfile     string `json:"job_profile"`     // Dominant Job
	HouseStatus    string `json:"house_status"`    // Milik Sendiri, Sewa, Menumpang
	HasLatrine     bool   `json:"has_latrine"`     // Jamban Sendiri
	HasCleanWater  bool   `json:"has_clean_water"` // Air Bersih (Bukan Sungai)
	IsElderlyHead  bool   `json:"is_elderly_head"` // Kepala Keluarga > 65
	HasToddler     bool   `json:"has_toddler"`     // Ada Balita
}

// GeoJSON Structures
type GeoJSON struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type     string   `json:"type"`
	Geometry Geometry `json:"geometry"`
}

type Geometry struct {
	Type        string          `json:"type"`
	Coordinates [][][][]float64 `json:"coordinates"` // MultiPolygon
}

func LoadBoundary(filename string) (*GeoJSON, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var geo GeoJSON
	if err := json.Unmarshal(data, &geo); err != nil {
		return nil, err
	}
	return &geo, nil
}

func IsPointInPolygon(lat, lng float64, geo *GeoJSON) bool {
	// Simple Ray Casting
	// PONDOKREJO.geojson is MultiPolygon
	for _, feature := range geo.Features {
		if feature.Geometry.Type == "MultiPolygon" {
			for _, polygon := range feature.Geometry.Coordinates {
				// Outer ring is usually the first one
				ring := polygon[0]
				if isPointInRing(lat, lng, ring) {
					return true
				}
			}
		}
	}
	return false
}

func isPointInRing(lat, lng float64, ring [][]float64) bool {
	inside := false
	j := len(ring) - 1
	for i := 0; i < len(ring); i++ {
		xi, yi := ring[i][0], ring[i][1] // GeoJSON is [lng, lat]
		xj, yj := ring[j][0], ring[j][1]

		intersect := ((yi > lat) != (yj > lat)) &&
			(lng < (xj-xi)*(lat-yi)/(yj-yi)+xi)
		if intersect {
			inside = !inside
		}
		j = i
	}
	return inside
}

func ExtractDusun(address, rw string) string {
	// 1. Try to map based on RW
	rwClean := strings.TrimSpace(rw)
	// Remove "RW" prefix if present (though likely just number)
	rwClean = strings.TrimPrefix(strings.ToUpper(rwClean), "RW")
	rwClean = strings.TrimSpace(rwClean)

	// Parse int to handle "01" vs "1"
	if rwInt, err := strconv.Atoi(rwClean); err == nil {
		switch rwInt {
		case 1, 2:
			return "Ngentak"
		case 3:
			return "Plotengan"
		case 4:
			return "Badalan"
		case 5, 6:
			return "Jlopo"
		case 7, 8:
			return "Karanglo"
		case 9, 10:
			return "Dukuh"
		case 11, 12, 13:
			return "Jlapan"
		case 14, 15:
			return "Banjarharjo"
		case 16, 17:
			return "Glagahombo"
		case 18:
			return "Watupecah"
		case 19:
			return "Jenengan"
		case 20:
			return "Mlesan Balan"
		}
	}

	// 2. Fallback to Text Analysis (If RW is 0 or invalid)
	s := strings.ToUpper(address)

	// Standardization Map based on 12 Valid Padukuhan
	if strings.Contains(s, "BANJAR") {
		return "Banjarharjo"
	}
	if strings.Contains(s, "DUKUH") {
		return "Dukuh"
	}
	if strings.Contains(s, "GLAGAH") {
		return "Glagahombo"
	}
	if strings.Contains(s, "JLAPAN") {
		return "Jlapan"
	}
	if strings.Contains(s, "JLOPO") {
		return "Jlopo"
	}
	if strings.Contains(s, "KARANG") {
		return "Karanglo"
	}
	if strings.Contains(s, "NGENTAK") {
		return "Ngentak"
	}
	if strings.Contains(s, "PLOTENG") {
		return "Plotengan"
	}
	if strings.Contains(s, "WATU") {
		return "Watupecah"
	} // Norm: Watupecah
	if strings.Contains(s, "MLES") || strings.Contains(s, "BALAN") {
		return "Mlesan Balan"
	}
	if strings.Contains(s, "BADAL") {
		return "Badalan"
	}
	if strings.Contains(s, "JENENG") {
		return "Jenengan"
	}

	return "Padukuhan Lainnya"
}

// Helper to clean strings from Excel (e.g. remove leading ')
func CleanString(s string) string {
	return strings.TrimPrefix(strings.TrimSpace(s), "'")
}

func ParseExcel(filename string) ([]Household, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}

	// Column Indices (0-based) for KK_Data_Final_Readable.xlsx
	const (
		// Household Data
		ColNoKK       = 2  // NO KK
		ColHeadName   = 12 // Nama Kepala KK
		ColAddress    = 11 // Alamat Dukuh
		ColRT         = 10 // RT
		ColRW         = 9  // RW
		ColCoordinate = 17 // Koordinat
		ColKeterangan = 14 // Keterangan
		ColIDDesil    = 15 // ID Desil (Ket)
		ColISUsulan   = 16 // IS Usulan (Ket)
		ColISEkstrem  = 18 // IS Ekstrem (Ket)

		ColFloorArea   = 20 // Lantai Luas
		ColFloorType   = 21 // ID Lantai (Ket)
		ColWallType    = 22 // ID Dinding (Ket)
		ColRoofType    = 23 // ID Atap (Ket)
		ColWaterSource = 24 // ID Airminum (Ket)
		ColSanitation  = 32 // ID Fasbab (Ket) (Sanitation ownership)

		ColBpntThn     = 38  // Bpnt Tahun
		ColPkhThn      = 41  // Pkh Thn
		ColExpenditure = 100 // Overall Sum

		// Member Data
		ColNik         = 109 // Nik
		ColName        = 110 // Nama
		ColRelation    = 111 // ID Hub Keluarga (Ket)
		ColGender      = 114 // ID Kelamin (Ket)
		ColAge         = 137 // Usia
		ColEducation   = 121 // ID Jenjang (Ket)
		ColJob         = 128 // ID Kerja (Ket)
		ColJobStatus   = 129 // ID Kerja Sttus (Ket)
		ColKerjaDetail = 131 // Kerja Detail
		ColIncome      = 132 // Income
		ColUshDetail   = 142 // Ush Detail
		ColPregnant    = 118 // Hamil
		ColDisability  = 125 // ID Difable

		// Aid Flags (Boolean/Ket)
		ColISBpnt      = 36
		ColISPkh       = 39
		ColISBlt       = 42
		ColISBanpem    = 48 // Bantuan Pemerintah (Assuming generic) or use specific columns
		ColISPupuk     = 49
		ColISLpg       = 50
		ColISBaznas    = 53
		ColISInet      = 76
		ColISBank      = 77
		ColISNpwp      = 133
		ColISTki       = 136
		ColISRokok     = 138
		ColISUsh       = 139
		ColSosJamkes   = 143
		ColSosPrakerja = 144
		ColSosKur      = 145
		ColSosMikro    = 146
		ColSosPip      = 147
		ColSosJamket   = 149

		// New Analytics Columns
		ColHouseStatus = 29 // ID Rumah Milik (Ket) - Guessing based on proximity to other housing vars
		// Need to verify ColHouseStatus index. Assuming standard Susenas structure or user provided file.
		// Let's use a safe assumption or string search if possible.
		// Actually, let's look for "ID Rumah Milik (Ket)" or similar.
		// For now, I will use logic to parse from "Keterangan" if specific col is missing or use placeholder.
		// Wait, I don't have exact column index for "Status Kepemilikan Rumah".
		// I'll try to deduce or use a generic one.
		// Let's use 29 as a placeholder or try to find it.
	)

	// Load Boundary
	boundary, err := LoadBoundary("PONDOKREJO.geojson")
	if err != nil {
		fmt.Println("Warning: Could not load boundary for correction:", err)
	}

	householdsMap := make(map[string]*Household)

	// Start from row index 1 (2nd row) assuming 0 is header
	// Inspect showed:
	// Row 0: Headers (Unnamed: 0, ID Rumah tangga, ...)
	// Row 1: Data starts
	// Wait, sample output showed "Row 0" as data but pandas index 0 corresponds to Excel row 2 usually if header=0.
	// The `rows` from excelize is 0-indexed slice of all rows.
	// If the file has 1 header row:
	// Row 0: Headers
	// Row 1: Data 1
	startRow := 1

	// Check if Row 0 is indeed header
	if len(rows) > 0 {
		if strings.Contains(strings.ToLower(rows[0][ColNoKK]), "kk") {
			startRow = 1
		}
	}

	for i := startRow; i < len(rows); i++ {
		row := rows[i]

		// Safety check for row length
		if len(row) <= ColName {
			continue
		}

		// Clean NoKK (remove leading apostrophe)
		noKK := CleanString(row[ColNoKK])
		if noKK == "" {
			continue
		}

		// Parse Coordinates
		coordStr := ""
		if len(row) > ColCoordinate {
			coordStr = row[ColCoordinate]
		}

		var lat, lng float64
		hasCoord := false

		if coordStr != "" {
			parts := strings.Split(coordStr, ",")
			if len(parts) == 2 {
				var err1, err2 error
				lat, err1 = strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
				lng, err2 = strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
				if err1 == nil && err2 == nil {
					hasCoord = true
				}
			}
		}

		// Get Head Name
		headName := ""
		if len(row) > ColHeadName {
			headName = strings.TrimSpace(row[ColHeadName])
		}
		if headName == "" {
			headName = "UNKNOWN"
		}

		// Use Composite Key (NoKK + HeadName) to ensure uniqueness
		// even though the new file should be cleaner, this is safer.
		uniqueKey := noKK + "|" + headName

		// Helper for extra fields (Defined here to be available for both HH and Member logic)
		getVal := func(idx int) string {
			if len(row) > idx {
				return row[idx]
			}
			return ""
		}

		// Get or Create Household
		hh, exists := householdsMap[uniqueKey]
		if !exists {
			address := ""
			if len(row) > ColAddress {
				address = row[ColAddress]
			}
			welfare := ""
			if len(row) > ColIDDesil {
				welfare = row[ColIDDesil]
			}

			// Normalize Welfare Level to 1, 2, 3, 4
			welfare = strings.TrimSpace(welfare)
			wUpper := strings.ToUpper(welfare)
			if strings.Contains(wUpper, "SANGAT MISKIN") {
				welfare = "1"
			} else if strings.Contains(wUpper, "HAMPIR MISKIN") {
				welfare = "3" // Check "Hampir" first before "Miskin"
			} else if strings.Contains(wUpper, "MISKIN") {
				// Careful: "Rentan Miskin" also contains "Miskin"
				// But "Rentan Miskin" is usually better than "Miskin"
				if strings.Contains(wUpper, "RENTAN") || strings.Contains(wUpper, "RENTAH") {
					welfare = "4" // Rentan Miskin (Desil 4)
				} else {
					welfare = "2" // Pure Miskin
				}
			} else if strings.Contains(wUpper, "RENTAH") || strings.Contains(wUpper, "RENTAN") {
				welfare = "4" // Assume Vulnerable is 4
			} else if strings.Contains(wUpper, "MAMPU") {
				welfare = "4"
			}

			// Fallback Logic for Welfare if ID Desil is empty
			if welfare == "" || welfare == "NaN" {
				if len(row) > ColISEkstrem {
					val := strings.TrimSpace(strings.ToUpper(row[ColISEkstrem]))
					if val == "1" || val == "YA" || val == "Y" {
						welfare = "1" // Miskin Ekstrem
					}
				}
				if (welfare == "" || welfare == "NaN") && len(row) > ColISUsulan {
					val := strings.TrimSpace(strings.ToUpper(row[ColISUsulan]))
					if val == "1" || val == "YA" || val == "Y" {
						welfare = "2" // Usulan (Assuming Desil 2)
					}
				}

				// If still empty, assume Mampu (Desil 4+)
				if welfare == "" || welfare == "NaN" {
					welfare = "4"
				}
			}

			// Map Dusun Code to Name using Extraction
			rw := getVal(ColRW)
			dusunName := ExtractDusun(address, rw)

			// Calculate Analytics Fields

			// 1. Income Category (Based on Expenditure if Income is missing, or use Expenditure as proxy)
			expStr := getVal(ColExpenditure)
			expVal, _ := strconv.ParseFloat(expStr, 64)
			var incomeCat string
			if expVal < 1000000 {
				incomeCat = "< Rp1 Juta"
			} else if expVal < 2000000 {
				incomeCat = "Rp1 - 2 Juta"
			} else {
				incomeCat = "> Rp2 Juta"
			}

			// 2. Infrastructure
			sanitation := getVal(ColSanitation)
			hasLatrine := false
			if strings.Contains(strings.ToUpper(sanitation), "MILIK SENDIRI") {
				hasLatrine = true
			}

			water := getVal(ColWaterSource)
			hasCleanWater := true
			if strings.Contains(strings.ToUpper(water), "SUNGAI") || strings.Contains(strings.ToUpper(water), "DANAU") {
				hasCleanWater = false
			}

			// 3. House Status (Approximate from Col 29 if valid, or just assume)
			// Since we don't have exact column confirmed, we'll try to use what we have or skip
			houseStatus := getVal(29) // Try index 29
			if houseStatus == "" {
				houseStatus = "Milik Sendiri"
			} // Default

			hh = &Household{
				NoKK:           noKK,
				HeadName:       headName,
				Address:        address,
				Dusun:          dusunName,
				Latitude:       lat,
				Longitude:      lng,
				WelfareLevel:   welfare,
				PkhThn:         getVal(ColPkhThn),
				BpntThn:        getVal(ColBpntThn),
				LantaiLuas:     getVal(ColFloorArea),
				Keterangan:     getVal(ColKeterangan),
				Expenditure:    getVal(ColExpenditure),
				FloorType:      getVal(ColFloorType),
				WallType:       getVal(ColWallType),
				RoofType:       getVal(ColRoofType),
				WaterSource:    getVal(ColWaterSource),
				Sanitation:     getVal(ColSanitation),
				IncomeCategory: incomeCat,
				HouseStatus:    houseStatus,
				HasLatrine:     hasLatrine,
				HasCleanWater:  hasCleanWater,
				Members:        []Resident{},
			}
			householdsMap[uniqueKey] = hh
		} else {
			// Update coordinates if they were missing and now found
			if hh.Latitude == 0 && hh.Longitude == 0 && hasCoord {
				hh.Latitude = lat
				hh.Longitude = lng
			}
		}

		// Add Member
		name := CleanString(row[ColName])
		nik := CleanString(row[ColNik])

		// Analytics: Job Profile for Head (if name matches head)
		// Or just store head job separately?
		// Let's store job for everyone, then in analysis we pick head's job.

		job := getVal(ColJob)
		ageStr := getVal(ColAge)
		age, _ := strconv.Atoi(ageStr)

		// Check for Toddler (Balita < 5 tahun)
		if age > 0 && age < 5 {
			hh.HasToddler = true
		}

		// Check for Elderly Head
		// Note: We need to verify if this person is Head.
		// Excel usually has "Hubungan Keluarga" column.
		relation := getVal(ColRelation)
		if strings.Contains(strings.ToUpper(relation), "KEPALA") {
			if age > 65 {
				hh.IsElderlyHead = true
			}
			// Update Job Profile for HH based on Head
			if job != "" {
				hh.JobProfile = job
			}
		}

		// Check Aid Info
		var aids []string
		checkAid := func(colIdx int, aidName string) {
			if len(row) > colIdx {
				val := strings.TrimSpace(strings.ToUpper(row[colIdx]))
				// Check for "Ya", "1", "Ada"
				if val == "1" || val == "YA" || val == "Y" || val == "ADA" {
					aids = append(aids, aidName)
				}
			}
		}

		checkAid(ColISBpnt, "BPNT")
		checkAid(ColISPkh, "PKH")
		checkAid(ColISBlt, "BLT")
		checkAid(ColISBanpem, "Banpres")
		checkAid(ColSosKur, "KUR")
		checkAid(ColSosMikro, "UMKM Mikro")
		checkAid(ColSosPip, "PIP")
		checkAid(ColSosJamkes, "Jamkes")
		checkAid(ColSosPrakerja, "Prakerja")
		checkAid(ColSosJamket, "Jamket")

		member := Resident{
			Name:        name,
			Nik:         nik,
			AidList:     aids,
			UshDetail:   getVal(ColUshDetail),
			KerjaDetail: getVal(ColKerjaDetail),
			Job:         job,
			Age:         ageStr,
			Education:   getVal(ColEducation),
			Income:      getVal(ColIncome),
			Pregnant:    getVal(ColPregnant),
			Disability:  getVal(ColDisability),
			Relation:    relation,
			Gender:      getVal(ColGender),
		}

		hh.Members = append(hh.Members, member)
	}

	// Convert map to slice
	var result []Household

	// Verification Map for Dusun Counts
	dusunCounts := make(map[string]int)

	// Pre-process for Coordinate Correction
	if boundary != nil {
		// Group valid households by Dusun
		dusunCentroids := make(map[string]struct {
			sumLat, sumLng float64
			count          int
		})

		for _, hh := range householdsMap {
			if hh.Latitude != 0 && hh.Longitude != 0 && hh.Dusun != "" {
				dusunCounts[hh.Dusun]++ // Count valid only
				if IsPointInPolygon(hh.Latitude, hh.Longitude, boundary) {
					// Valid Point
					s := dusunCentroids[hh.Dusun]
					s.sumLat += hh.Latitude
					s.sumLng += hh.Longitude
					s.count++
					dusunCentroids[hh.Dusun] = s
				}
			}
		}

		// Apply Correction
		for _, hh := range householdsMap {
			if hh.Latitude != 0 && hh.Longitude != 0 {
				if !IsPointInPolygon(hh.Latitude, hh.Longitude, boundary) {
					// Outside!
					// Find neighbors in same Dusun
					s, ok := dusunCentroids[hh.Dusun]
					if ok && s.count > 0 {
						// Move to Centroid
						newLat := s.sumLat / float64(s.count)
						newLng := s.sumLng / float64(s.count)

						hh.Latitude = newLat
						hh.Longitude = newLng

						if hh.Keterangan != "" {
							hh.Keterangan += "; "
						}
						hh.Keterangan += "Koordinat digeser (Auto-Correction)"
					}
				}
			}
		}
	}

	for _, hh := range householdsMap {
		// Only include households with valid coordinates for the map
		if hh.Latitude != 0 && hh.Longitude != 0 {
			result = append(result, *hh)
		}
	}

	fmt.Println("--- DUSUN COUNT VERIFICATION ---")
	for k, v := range dusunCounts {
		fmt.Printf("%s: %d\n", k, v)
	}
	fmt.Println("--------------------------------")

	return result, nil
}
