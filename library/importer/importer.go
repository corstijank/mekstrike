package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/corstijank/mekstrike/domain/storage"
	"github.com/corstijank/mekstrike/domain/unit"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/gocolly/colly"
)

func main() {
	log.Printf("Starting import of new units")
	log.Printf("Starting with 10 sec sleep to allow daprd sidecar to come up")
	time.Sleep(10 * time.Second)
	ctx := context.Background()
	s, err := storage.New("library-store")
	if err != nil {
		log.Printf("Error connecting to store:")
		log.Panic(err)
	}
	client, err := dapr.NewClient()
	if err != nil {
		log.Printf("Error creating a DAPR client")

		log.Panic(err)
	}

	log.Printf("Checking if there are already units")
	ir, err := s.ReadMany(ctx, "_units", &unit.Stats{})
	if err != nil {
		log.Panic(err)
	}
	if len(ir) > 0 {
		log.Printf("Already found units in store, exitting")
		os.Exit(0)
	}

	log.Printf("No units exist; scraping from www.masterunitlist.info")

	units := make([]*unit.Stats, 0)

	// Table Collector
	tc := colly.NewCollector()
	// Unit Collector
	uc := colly.NewCollector()

	// Find and visit all mechs on the table collector
	tc.OnHTML("tr td:nth-of-type(1)", func(e *colly.HTMLElement) {
		name := e.Text
		// Skip omni BASE variants
		if !strings.Contains(name, "<Base>") {
			subPath := e.ChildAttr("a", "href")
			identifier := strings.Split(subPath, "/")[3]
			fmt.Printf("Importing %s (%s)\n", e.Text, identifier)
			uc.Visit(e.Request.AbsoluteURL(fmt.Sprintf("/Tools/CustomCard/%s", identifier)))
		}
	})

	uc.OnHTML("form[action*='/Tools/CardGenerator']", func(e *colly.HTMLElement) {
		units = append(units, &unit.Stats{
			Name:       readStringValueFromInput(e, "Name"),
			Model:      readStringValueFromInput(e, "Model"),
			Pointvalue: readIntValueFromInput(e, "PV"),
			Type:       readStringValueFromInput(e, "Type"),
			Size:       readIntValueFromInput(e, "Size"),
			Movement:   readStringValueFromInput(e, "Move"),
			Role:       readStringValueFromInput(e, "Role"),
			Shortdmg:   readIntValueFromInput(e, "Short"),
			Meddmg:     readIntValueFromInput(e, "Medium"),
			Longdmg:    readIntValueFromInput(e, "Long"),
			Ovhdmg:     readIntValueFromInput(e, "Overheat"),
			Armor:      readIntValueFromInput(e, "Armor"),
			Struct:     readIntValueFromInput(e, "Structure"),
			Specials:   strings.Split(e.ChildText("textarea[name='Data.Specials']"), ","),
			Image:      e.ChildText("textarea[name='Data.Image']"),
		})
	})

	// Only introductory mechs
	tc.Visit("http://www.masterunitlist.info/Unit/Filter?Name=&HasBV=false&MinTons=&MaxTons=&MinBV=&MaxBV=&MinIntro=&MaxIntro=&MinCost=&MaxCost=&HasRole=&HasBFAbility=&MinPV=&MaxPV=&Role=None+Selected&Rules=55&Types=18&BookAuto=&FactionAuto=&AvailableEras=256")

	log.Printf("Saving units to store")

	for _, u := range units {
		err := s.Persist(ctx, u)
		if err != nil {
			log.Println(err)
		}
	}

	s.Close()
	time.Sleep(5 * time.Second)
	client.Close()

	client.Shutdown(ctx)

	log.Println("All done! You can close now")
}

func readStringValueFromInput(e *colly.HTMLElement, inputName string) string {
	return e.ChildAttr("input[name='Data."+inputName+"']", "value")
}

func readIntValueFromInput(e *colly.HTMLElement, inputName string) int32 {
	if e.ChildAttr("input[name='Data."+inputName+"']", "value") != "" {
		result, err := strconv.Atoi(e.ChildAttr("input[name='Data."+inputName+"']", "value"))
		if err != nil {
			log.Printf("Error with %s: %s", inputName, err)
		}
		return int32(result)
	}
	return 0
}
