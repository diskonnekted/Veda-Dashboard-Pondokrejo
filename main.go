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

	editorStore := NewEditorStore("editor_state.json")
	if err := editorStore.Load(); err != nil {
		log.Fatalf("Error loading editor state: %v", err)
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

	r.GET("/editor", func(c *gin.Context) {
		c.HTML(http.StatusOK, "editor.html", nil)
	})

	r.GET("/editor/layers", func(c *gin.Context) {
		c.JSON(http.StatusOK, editorStore.GeoLayersSnapshot())
	})

	r.GET("/editor/layers/:name", func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			c.Data(http.StatusOK, "application/json", []byte(`{"type":"FeatureCollection","features":[]}`))
			return
		}
		if raw, ok := editorStore.GeoLayer(name); ok {
			c.Data(http.StatusOK, "application/json", raw)
			return
		}
		c.Data(http.StatusOK, "application/json", []byte(`{"type":"FeatureCollection","features":[]}`))
	})

	// Serve residents.json dynamically (matching static filename)
	r.GET("/residents.json", func(c *gin.Context) {
		c.JSON(http.StatusOK, editorStore.Apply(households))
	})

	r.POST("/editor/save", func(c *gin.Context) {
		var req EditorSaveRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": err.Error()})
			return
		}
		if err := editorStore.UpdateAndSave(req); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// Serve boundary.geojson (matching static filename)
	r.StaticFile("/boundary.geojson", "PONDOKREJO.geojson")

	// Serve Layers (GeoJSON) from 'layers' directory
	// This replaces all individual r.StaticFile calls for layers
	// Ensure the 'layers' directory exists in the root
	r.Static("/layers", "./layers")

	// Serve Images
	r.StaticFile("/logo.png", "./logo.png")
	r.StaticFile("/veda-logo.png", "./veda-logo.png")
	r.StaticFile("/clasnet-logo.png", "./clasnet-logo.png")
	r.StaticFile("/login.jpg", "./login.jpg")
	r.StaticFile("/background.jpg", "./background.jpg")

	// Recommendations Page
	r.GET("/recommendations", func(c *gin.Context) {
		rek, err1 := os.ReadFile("rekomendasi-hasil-analitik.md")
		if err1 != nil {
			rek = []byte("# Error loading recommendation file")
		}

		act, err2 := os.ReadFile("action-plan.md")
		if err2 != nil {
			act = []byte("# Error loading action plan file")
		}

		gis, err3 := os.ReadFile("gis-desa.md")
		if err3 != nil {
			gis = []byte("# Error loading GIS file")
		}

		c.HTML(http.StatusOK, "recommendations.html", gin.H{
			"Rekomendasi": string(rek),
			"ActionPlan":  string(act),
			"GISDesa":     string(gis),
		})
	})

	// Analytics Dashboard Page
	r.GET("/analytics", func(c *gin.Context) {
		// Calculate on the fly (or cache it)
		data := CalculateAnalytics(households)

		// Convert to JSON for injection
		jsonData, err := json.Marshal(data)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error processing data")
			return
		}

		c.HTML(http.StatusOK, "analytics.html", gin.H{
			"AnalyticsData": string(jsonData),
		})
	})

	// Bansos Dashboard Page
	r.GET("/bansos", func(c *gin.Context) {
		data := CalculateBansos(households)
		jsonData, err := json.Marshal(data)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error processing bansos data")
			return
		}
		c.HTML(http.StatusOK, "bansos.html", gin.H{
			"BansosData": string(jsonData),
		})
	})

	// Bansos JSON API
	r.GET("/bansos/data", func(c *gin.Context) {
		data := CalculateBansos(households)
		c.JSON(http.StatusOK, data)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server starting on :" + port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
