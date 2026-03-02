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
	Gender      string   `json:"gender"` // ID Kelamin
	AidList     []string `json:"aid_list"`
	KerjaDetail string   `json:"kerja_detail"` // Kerja Detail
	UshDetail   string   `json:"ush_detail"`   // Ush Detail
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

			dusunCode := ""
			if len(row) > ColDusun {
				dusunCode = row[ColDusun]
			}
			
			// Map Dusun Code to Name (This is a guess, but better than nothing. User can correct later)
			// Based on sample: 1 might be Ngentak. 
			// We can also try to parse from Address if needed, but let's use the Code for grouping.
			// Or we can just pass the code. Let's pass the code for now and label it 'Padukuhan X'.
			// Ideally we would have a mapping table.
			
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

			hh = &Household{
				NoKK:         noKK,
				HeadName:     headName,
				Address:      address,
				Dusun:        dusunCode,
				Latitude:     lat,
				Longitude:    lng,
				WelfareLevel: welfare,
				PkhThn:       pkhThn,
				BpntThn:      bpntThn,
				LantaiLuas:   lantaiLuas,
				Keterangan:   keterangan,
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

		member := Resident{
			Name:        name,
			Nik:         nik,
			AidList:     aids,
			UshDetail:   ushDetail,
			KerjaDetail: kerjaDetail,
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
