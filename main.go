package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Parse CLI flags
	genStatic := flag.Bool("gen", false, "Generate static JSON files for deployment")
	clipGeoJSON := flag.String("clip-geojson", "", "Clip input GeoJSON to Pondokrejo boundary and write output")
	clipOut := flag.String("clip-out", "", "Output GeoJSON file path for -clip-geojson")
	clipBoundary := flag.String("clip-boundary", "DEFAULT_VIEW", "Boundary GeoJSON file for -clip-geojson")
	flag.Parse()

	if *clipGeoJSON != "" {
		if err := ClipGeoJSONToPondokrejoBoundary(*clipGeoJSON, *clipBoundary, *clipOut); err != nil {
			log.Fatalf("Error clipping geojson: %v", err)
		}
		fmt.Println("Clipped GeoJSON written successfully")
		return
	}

	// Parse Excel Data
	households, err := ParseExcel("1 KK_ART Pondokrejo.xlsx")
	if err != nil {
		log.Fatalf("Error parsing excel: %v", err)
	}

	// Generate Static Files Mode
	if *genStatic {
		fmt.Println("Generating static files...")

		// 1. residents.json
		jsonData, err := json.MarshalIndent(households, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile("residents.json", jsonData, 0644); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Created residents.json")

		// 2. boundary.geojson (Copy)
		boundaryData, err := os.ReadFile("PONDOKREJO.geojson")
		if err != nil {
			log.Printf("Warning: Could not read PONDOKREJO.geojson: %v", err)
		} else {
			if err := os.WriteFile("boundary.geojson", boundaryData, 0644); err != nil {
				log.Fatal(err)
			}
			fmt.Println("Created boundary.geojson")
		}

		fmt.Println("Static build complete. You can now upload 'index.html', 'residents.json', and 'boundary.geojson' to a static host.")
		return
	}

	// Server Mode
	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	// Serve residents.json dynamically (matching static filename)
	r.GET("/residents.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, households)
	})

	// Serve boundary.geojson (matching static filename)
	r.StaticFile("/boundary.geojson", "PONDOKREJO.geojson")

	// Serve clipped geographic layers (generated via -clip-geojson)
	r.StaticFile("/layers/pemukiman-area.geojson", "pemukiman-area-pondokrejo.geojson")
	r.StaticFile("/layers/pemukiman-area-pondokrejo.geojson", "pemukiman-area-pondokrejo.geojson")
	r.StaticFile("/layers/sungai-line-pondokrejo.geojson", "sungai-line-pondokrejo.geojson")
	r.StaticFile("/layers/sawah-area.geojson", "sawah-area.json")
	r.StaticFile("/layers/bangunan-point-pondokrejo.geojson", "bangunan-point-pondokrejo.json")
	r.StaticFile("/layers/jalan-line-pondokrejo.geojson", "jalan-line-pondokrejo.json")
	r.StaticFile("/layers/kontur-line-pondokrejo.geojson", "kontur-line-pondokrejo.json")
	r.StaticFile("/layers/irigasi-line-pondokrejo.geojson", "irigasi-line-pondokrejo.geojson")
	r.StaticFile("/layers/pendidikan-point-pondokrejo.geojson", "pendidikan-point-pondokrejo.geojson")
	r.StaticFile("/layers/toponimi-point-pondokrejo.geojson", "toponimi-point-pondokrejo.geojson")
	r.StaticFile("/layers/tonggak-km-point-pondokrejo.geojson", "tonggak-km-point-pondokrejo.geojson")
	r.StaticFile("/layers/pertambangan-point-pondokrejo.geojson", "pertambangan-point-pondokrejo.geojson")

	// New layers extracted from ArcGIS
	r.StaticFile("/layers/waste_banks.geojson", "waste_banks.geojson")
	r.StaticFile("/layers/tps_locations.geojson", "tps_locations.geojson")
	r.StaticFile("/layers/waste_routes.geojson", "waste_routes.geojson")
	r.StaticFile("/layers/district_boundary.geojson", "district_boundary.geojson")

	// New layers extracted from ArcGIS (Disaster Mitigation)
	r.StaticFile("/layers/krb_1.geojson", "krb_1.geojson")
	r.StaticFile("/layers/krb_2.geojson", "krb_2.geojson")
	r.StaticFile("/layers/krb_3.geojson", "krb_3.geojson")
	r.StaticFile("/layers/merapi_distance.geojson", "merapi_distance.geojson")
	r.StaticFile("/layers/evacuation_route.geojson", "evacuation_route.geojson")
	r.StaticFile("/layers/ews_location.geojson", "ews_location.geojson")
	r.StaticFile("/layers/barracks.geojson", "barracks.geojson")
	r.StaticFile("/layers/assembly_points.geojson", "assembly_points.geojson")

	// Serve Images
	r.StaticFile("/logo.png", "./logo.png")
	r.StaticFile("/veda-logo.png", "./veda-logo.png")
	r.StaticFile("/clasnet-logo.png", "./clasnet-logo.png")

	// Start server
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
