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
	flag.Parse()

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
