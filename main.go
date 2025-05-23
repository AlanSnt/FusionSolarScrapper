package main

import (
	"AlanSnt/FusionSolarScrapper/mqtt"
	"AlanSnt/FusionSolarScrapper/scrapper"
	"AlanSnt/FusionSolarScrapper/settings"
	"log"
	"net"
	"os/exec"
	"time"

	"math"

	"github.com/joho/godotenv"
)

func scrapp() {
	solarProduction := 0.0
	returnedEnergy := 0.0
	consumedEnergy := 0.0

	solarProduction = scrapper.GetSolarData() * 1000
	returnedEnergy = scrapper.GetReturnedEnergy()

	if returnedEnergy < 0 {
		consumedEnergy = solarProduction + math.Abs(returnedEnergy)
		returnedEnergy = 0
	} else {
		consumedEnergy = solarProduction - returnedEnergy
	}

	mqtt.SendMessage("production", solarProduction/1000)
	mqtt.SendMessage("returned", returnedEnergy/1000)
	mqtt.SendMessage("consumed", consumedEnergy/1000)
}

func tryNetwork() {
	host := "www.google.com:80"
	timeout := 5 * time.Second

	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		panic("Network error, please check your network connection: " + err.Error())
	}
	defer conn.Close()
}

func clearRootCache() {
	cmd := "rm -rf /root/.cache/ms-playwright/*"
	_, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Fatal("Error:", err)
	}

	cmd = "rm -rf /root/.cache/playwright/*"
	_, err = exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func main() {
	clearRootCache()
	godotenv.Load()

	settings.Init()
	mqtt.Init()
	scrapper.Init()

	tryNetwork()

	log.Print("Fusion Solar scrapper is ready")

	timeDelta, err := settings.Get("TIME_DELTA")
	if err != nil {
		log.Fatal("Error:", err)
	}

	for {
		scrapp()
		time.Sleep(time.Duration(timeDelta.(int)) * time.Second)
	}
}
