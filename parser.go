package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Resident struct {
	Name        string   `json:"name"`
	Nik         string   `json:"nik"`
	Gender      string   `json:"gender"` // ID Kelamin
	AidList     []string `json:"aid_list"`
	KerjaDetail string   `json:"kerja_detail"` // Kerja Detail
	UshDetail   string   `json:"ush_detail"`   // Ush Detail
	Age         string   `json:"age"`          // Usia
	Education   string   `json:"education"`    // ID Ijazah
	Income      string   `json:"income"`       // Income
	Pregnant    string   `json:"pregnant"`     // Hamil
	Disability  string   `json:"disability"`   // ID Difable
}

type Household struct {
	NoKK         string     `json:"no_kk"`
	HeadName     string     `json:"head_name"`
	Address      string     `json:"address"`
	Dusun        string     `json:"dusun"`
	Latitude     float64    `json:"latitude"`
	Longitude    float64    `json:"longitude"`
	WelfareLevel string     `json:"welfare_level"` // ID Desil
	Members      []Resident `json:"members"`
	PkhThn       string     `json:"pkh_thn"`       // Pkh Thn
	BpntThn      string     `json:"bpnt_thn"`      // Bpnt Thn
	LantaiLuas   string     `json:"lantai_luas"`   // Lantai Luas
	Keterangan   string     `json:"keterangan"`    // Keterangan
	Expenditure  string     `json:"expenditure"`   // Overall Sum
	FloorType    string     `json:"floor_type"`    // ID Lantai
	WallType     string     `json:"wall_type"`     // ID Dinding
	RoofType     string     `json:"roof_type"`     // ID Atap
	WaterSource  string     `json:"water_source"`  // ID Airminum
	Sanitation   string     `json:"sanitation"`    // ID Fasbab
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

func ExtractDusun(address string) string {
	// Normalize
	s := strings.ToUpper(address)
	
	// Remove common prefixes/suffixes for hamlet
	// PADUKUHAN, DUSUN, DSN, DK
	rePrefix := regexp.MustCompile(`\b(PADUKUHAN|DUSUN|DSN|DK)\b`)
	s = rePrefix.ReplaceAllString(s, " ")

	// Remove RT/RW and numbers
	// Regex for RT/RW followed by optional numbers/slash
	re := regexp.MustCompile(`(RT|RW)\s*[\d\/\.\-]+`)
	s = re.ReplaceAllString(s, "")

	// Remove Roman Numerals (I, II, III, IV, V) - common in hamlet sections
	reRoman := regexp.MustCompile(`\b(I|II|III|IV|V|VI|VII|VIII|IX|X)\b`)
	s = reRoman.ReplaceAllString(s, "")

	// Remove non-alphabetic chars except spaces
	reNonAlpha := regexp.MustCompile(`[^A-Z\s]`)
	s = reNonAlpha.ReplaceAllString(s, " ")

	// Trim spaces and extra whitespace
	s = strings.TrimSpace(s)
	reSpace := regexp.MustCompile(`\s+`)
	s = reSpace.ReplaceAllString(s, " ")

	// Check for "DUKU" or "DUKUH" specifically as the only content or explicit word
	if s == "DUKU" || s == "DUKUH" {
		return "Dukuh" // Standardize to "Dukuh" as per user request
	}

	// Convert to Title Case for display
	s = strings.Title(strings.ToLower(s))
	
	// Standardization / Correction Map
	// Ngentak, Plotengan, Jlopo, Karanglo, Dukuh, Jlapan, Banjarharjo, Glagahombo, Watu Pecah
	switch s {
	case "Ngentak":
		return "Ngentak"
	case "Plotengan":
		return "Plotengan"
	case "Jlopo":
		return "Jlopo"
	case "Karanglo":
		return "Karanglo"
	case "Dukuh":
		return "Dukuh"
	case "Jlapan":
		return "Jlapan"
	case "Banjarharjo":
		return "Banjarharjo"
	case "Glagahombo", "Glagah Ombo":
		return "Glagahombo"
	case "Watu Pecah", "Watupecah":
		return "Watu Pecah"
	default:
		// Fuzzy matching or corrections
		if strings.Contains(s, "Glagah") { return "Glagahombo" }
		if strings.Contains(s, "Watu") { return "Watu Pecah" }
		if strings.Contains(s, "Banjar") { return "Banjarharjo" }
		if strings.Contains(s, "Jlapan") { return "Jlapan" }
		if strings.Contains(s, "Karang") { return "Karanglo" } // Be careful if other Karang exists
		if strings.Contains(s, "Ploteng") { return "Plotengan" }
		if strings.Contains(s, "Ngentak") { return "Ngentak" }
		if strings.Contains(s, "Jlopo") { return "Jlopo" }
		if strings.Contains(s, "Dukuh") { return "Dukuh" }
	}

	return "Padukuhan Lainnya"
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

	// Column Indices (0-based)
	// Based on previous analysis
	const (
		ColNoKK       = 2
		ColHeadName   = 14
		ColAddress    = 12
		ColDusun      = 9 // Using ID Dus (which seems to be numeric code)
		ColCoordinate = 44
		ColIDDesil    = 38
		ColName       = 232
		ColNik        = 231
		ColGender     = 237 // ID Kelamin
		
		// Aid Columns
		ColISBpnt   = 105
		ColISPkh    = 108
		ColISBlt    = 111
		ColISBanpem = 117
		ColSosKur   = 291
		ColSosMikro = 292
		ColSosPip   = 293
		ColSosJamket= 294
		
		// New Columns
		ColUshDetail   = 268 // Ush Detail
		ColKerjaDetail = 255 // Kerja Detail
		ColPkhThn      = 110 // Pkh Thn
		ColBpntThn     = 107 // Bpnt Thn
		ColLantaiLuas  = 55  // Lantai Luas
		ColKeterangan  = 19  // Keterangan
		ColExpenditure = 183
		ColFloorType   = 56
		ColWallType    = 57
		ColRoofType    = 58
		ColWaterSource = 59
		ColSanitation  = 71
		ColAge         = 260
		ColEducation   = 247
		ColMemIncome   = 256
		ColPregnant    = 242
		ColDisability  = 249
	)

	// Load Boundary
	boundary, err := LoadBoundary("PONDOKREJO.geojson")
	if err != nil {
		// Just log warning, don't fail, maybe just skip correction
		fmt.Println("Warning: Could not load boundary for correction:", err)
	}

	householdsMap := make(map[string]*Household)

	// Start from row index 3 (4th row) which is the first data row
	// Row 0: Laporan...
	// Row 1: Codes...
	// Row 2: Headers...
	// Row 3: Data...
	for i := 3; i < len(rows); i++ {
		row := rows[i]
		
		// Safety check for row length
		if len(row) <= ColName {
			continue
		}

		noKK := row[ColNoKK]
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

		// If no coordinates, we might skip or put at 0,0. 
		// For mapping, skipping is better, or user might want to see them in a list.
		// For now, we only map those with coordinates.
		if !hasCoord {
			// Try to find if household already exists and has coordinates
			if h, exists := householdsMap[noKK]; exists && h.Latitude != 0 {
				lat = h.Latitude
				lng = h.Longitude
			} else {
				// Skip if absolutely no coordinates for this household yet
				// But maybe a later member has coordinates?
				// We'll create the household with 0,0 and update if we find coords later.
			}
		}

		// Get or Create Household
		hh, exists := householdsMap[noKK]
		if !exists {
			headName := ""
			if len(row) > ColHeadName {
				headName = row[ColHeadName]
			}
			address := ""
			if len(row) > ColAddress {
				address = row[ColAddress]
			}
			welfare := ""
			if len(row) > ColIDDesil {
				welfare = row[ColIDDesil]
			}

			// Map Dusun Code to Name using Extraction
			dusunName := ExtractDusun(address)
			
			// Get Household specific fields (Assuming these are consistent for the head or we take the first non-empty)
			pkhThn := ""
			if len(row) > ColPkhThn {
				pkhThn = row[ColPkhThn]
			}
			bpntThn := ""
			if len(row) > ColBpntThn {
				bpntThn = row[ColBpntThn]
			}
			lantaiLuas := ""
			if len(row) > ColLantaiLuas {
				lantaiLuas = row[ColLantaiLuas]
			}
			keterangan := ""
			if len(row) > ColKeterangan {
				keterangan = row[ColKeterangan]
			}

			// Helper for extra fields
			getVal := func(idx int) string {
				if len(row) > idx {
					return row[idx]
				}
				return ""
			}

			hh = &Household{
				NoKK:         noKK,
				HeadName:     headName,
				Address:      address,
				Dusun:        dusunName, // Use Extracted Name
				Latitude:     lat,
				Longitude:    lng,
				WelfareLevel: welfare,
				PkhThn:       pkhThn,
				BpntThn:      bpntThn,
				LantaiLuas:   lantaiLuas,
				Keterangan:   keterangan,
				Expenditure:  getVal(ColExpenditure),
				FloorType:    getVal(ColFloorType),
				WallType:     getVal(ColWallType),
				RoofType:     getVal(ColRoofType),
				WaterSource:  getVal(ColWaterSource),
				Sanitation:   getVal(ColSanitation),
				Members:      []Resident{},
			}
			householdsMap[noKK] = hh
		} else {
			// Update coordinates if they were missing and now found
			if hh.Latitude == 0 && hh.Longitude == 0 && hasCoord {
				hh.Latitude = lat
				hh.Longitude = lng
			}
			
			// Update Household fields if empty and we found data now
			if hh.PkhThn == "" && len(row) > ColPkhThn { hh.PkhThn = row[ColPkhThn] }
			if hh.BpntThn == "" && len(row) > ColBpntThn { hh.BpntThn = row[ColBpntThn] }
			if hh.LantaiLuas == "" && len(row) > ColLantaiLuas { hh.LantaiLuas = row[ColLantaiLuas] }
			if hh.Keterangan == "" && len(row) > ColKeterangan { hh.Keterangan = row[ColKeterangan] }
		}

		// Add Member
		name := row[ColName]
		nik := row[ColNik]
		
		ushDetail := ""
		if len(row) > ColUshDetail {
			ushDetail = row[ColUshDetail]
		}
		kerjaDetail := ""
		if len(row) > ColKerjaDetail {
			kerjaDetail = row[ColKerjaDetail]
		}
		
		// Collect Aid Info
		var aids []string
		
		checkAid := func(colIdx int, aidName string) {
			if len(row) > colIdx {
				val := strings.TrimSpace(row[colIdx])
				// Assuming '1' means yes, or specific code. 
				// The user data shows '1' or '2' or empty.
				// Let's assume '1' is Yes.
				if val == "1" {
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
		checkAid(ColSosJamket, "Jamkes")

		getMVal := func(idx int) string {
			if len(row) > idx {
				return row[idx]
			}
			return ""
		}

		member := Resident{
			Name:        name,
			Nik:         nik,
			AidList:     aids,
			UshDetail:   ushDetail,
			KerjaDetail: kerjaDetail,
			Age:         getMVal(ColAge),
			Education:   getMVal(ColEducation),
			Income:      getMVal(ColMemIncome),
			Pregnant:    getMVal(ColPregnant),
			Disability:  getMVal(ColDisability),
		}
		
		hh.Members = append(hh.Members, member)
	}

	// Convert map to slice
	var result []Household
	
	// Pre-process for Coordinate Correction
	if boundary != nil {
		// Group valid households by Dusun
		dusunCentroids := make(map[string]struct{sumLat, sumLng float64; count int})
		
		for _, hh := range householdsMap {
			if hh.Latitude != 0 && hh.Longitude != 0 && hh.Dusun != "" {
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
						
						// Append note
						if hh.Keterangan != "" {
							hh.Keterangan += "; "
						}
						hh.Keterangan += "Koordinat perlu revisi (Digeser otomatis)"
					} else {
						// If no valid neighbors in same Dusun (unlikely if Dusun code is valid),
						// we might want to just mark it.
						if hh.Keterangan != "" {
							hh.Keterangan += "; "
						}
						hh.Keterangan += "Koordinat diluar wilayah (Perlu revisi)"
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

	return result, nil
}
