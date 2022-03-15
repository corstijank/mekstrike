package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
)

var exactSprites = make(map[string]string)
var chassisSprites = make(map[string]string)

func main() {
	log.Printf("### Starting Mekstrike Media Proxy")

	response, err := http.Get("https://raw.githubusercontent.com/MegaMek/megamek/master/megamek/data/images/units/mechset.txt") //use package "net/http"
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()
	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "chassis") {
			elements := regexp.MustCompile(`"(.*?)\"`).FindAll([]byte(line), -1)
			if len(elements) == 2 {
				key := strings.ReplaceAll(strings.ToLower(string(elements[0])), "\"", "")
				value := strings.ReplaceAll(string(elements[1]), "\"", "")
				chassisSprites[key] = value
			}
		} else if strings.HasPrefix(line, "exact") {
			elements := regexp.MustCompile(`"(.*?)\"`).FindAll([]byte(line), -1)
			if len(elements) == 2 {
				key := strings.ReplaceAll(strings.ToLower(string(elements[0])), "\"", "")
				value := strings.ReplaceAll(string(elements[1]), "\"", "")
				exactSprites[key] = value
			}
		}
	}
	r := mux.NewRouter()
	r.HandleFunc("/sprites/{unit}", getSprite).Methods("GET")

	log.Fatal(http.ListenAndServe(":9010", r))
}

func getSprite(rw http.ResponseWriter, r *http.Request) {
	u := strings.ToLower(mux.Vars(r)["unit"])
	log.Printf("MediaProxy::getSprite %+v", u)

	sprite, found := exactSprites[u]
	if !found {
		sprite, found = chassisSprites[u]
	}
	if !found {
		sprite, found = chassisSprites[strings.Split(u, " ")[0]]
	}

	url := "https://raw.githubusercontent.com/MegaMek/megamek/master/megamek/data/images/units/" + sprite
	if !found {
		url = "https://raw.githubusercontent.com/MegaMek/megamek/master/megamek/data/images/units/defaults/default_heavy.png"
	}
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()
	io.Copy(rw, response.Body)

}
