package scrapper

import (
	"log"
	"strings"

	"github.com/playwright-community/playwright-go"

	"strconv"

	"AlanSnt/FusionSolarScrapper/settings"
)

var page playwright.Page

func login() {
	page.Reload()

	username, err := settings.Get("USERNAME")
	if err != nil {
		log.Fatal("Error:", err)
	}

	password, err := settings.Get("PASSWORD")
	if err != nil {
		log.Fatal("Error:", err)
	}

	loginUserNameEntry, err := page.Locator("[name='ssoCredentials.username']").All()
	if err != nil {
		log.Fatal("Error:", err)
	}

	if len(loginUserNameEntry) < 2 {
		return
	}

	loginPasswordEntry, err := page.Locator("[name='ssoCredentials.value']").All()
	if err != nil {
		log.Fatal("Error:", err)
	}

	loginUserNameEntry[1].Fill(username.(string))
	loginPasswordEntry[1].Fill(password.(string))

	loginButton, err := page.Locator(".loginBtn").All()
	if err != nil {
		log.Fatal("Error:", err)
	}

	loginButton[0].Click()
}

func scrappeData() (string, string) {
	log.Println("Scraping data...")
	page.WaitForSelector(".nco-single-energy-total-content > span[title]")

	entries, err := page.Locator(".nco-single-energy-total-content > span[title]").All()
	if err != nil {
		log.Fatal("Error:", err)
	}

	productionEntry, err := entries[0].TextContent()
	if err != nil {
		log.Fatal("Error:", err)
	}

	consumptionEntry, err := entries[1].TextContent()
	if err != nil {
		log.Fatal("Error:", err)
	}

	production := strings.ReplaceAll(productionEntry, "\u00A0", "")
	production = strings.TrimSpace(production)
	production = strings.ReplaceAll(production, ",", ".")

	consumption := strings.ReplaceAll(consumptionEntry, "\u00A0", "")
	consumption = strings.TrimSpace(consumption)
	consumption = strings.ReplaceAll(consumption, ",", ".")

	return production, consumption
}

func GetReturnedEnergy() float64 {
	url, err := settings.Get("FUSION_SOLAR_URL")
	if err != nil {
		log.Fatal("Error:", err)
	}

	page.Goto(url.(string))
	login()
	page.WaitForSelector(".nco-monitor-station-overview-management")
	err = page.Locator(".flex-node-line-expand-part").Last().Click()
	if err != nil {
		log.Fatal("Error:", err)
	}

	page.WaitForSelector("[title='INV-HV2250159657']")
	err = page.Locator(".flex-node-line-expand-part").Last().Click()
	if err != nil {
		log.Fatal("Error:", err)
	}

	page.WaitForSelector("[title='Meter-1']")
	err = page.GetByTitle("Meter-1").First().Click()
	if err != nil {
		log.Fatal("Error:", err)
	}

	page.WaitForSelector(".realtime-content")
	page.WaitForTimeout(5000)

	entry, err := page.Locator(".realtime-content > .ant-row > .even-line > .has-padding > div").All()
	if err != nil {
		log.Fatal("Error:", err)
	}

	returnedEnergy, err := entry[1].TextContent()
	returnedEnergy = strings.ReplaceAll(returnedEnergy, " W", "")
	returnedEnergy = strings.TrimSpace(returnedEnergy)
	returnedEnergy = strings.ReplaceAll(returnedEnergy, " ", "")
	returnedEnergy = strings.ReplaceAll(returnedEnergy, "\u202f", "")
	returnedEnergy = strings.ReplaceAll(returnedEnergy, ",", ".")

	returnedEnergyFloat, err := strconv.ParseFloat(returnedEnergy, 64)
	if err != nil {
		log.Fatal("Error: convert string to float", returnedEnergy, err)
	}

	return returnedEnergyFloat
}

func GetSolarData() float64 {
	url, err := settings.Get("FUSION_SOLAR_URL")
	if err != nil {
		log.Fatal("Error:", err)
	}

	page.Goto(url.(string))
	login()
	page.WaitForSelector(".nco-monitor-station-overview-management")

	err = page.Locator(".flex-node-line-expand-part").Last().Click()
	if err != nil {
		log.Fatal("Error:", err)
	}

	page.WaitForSelector("[title='INV-HV2250159657']")

	err = page.GetByTitle("INV-HV2250159657").First().Click()
	if err != nil {
		log.Fatal("Error:", err)
	}

	page.WaitForSelector(".realtime-content")
	page.WaitForTimeout(5000)
	entry, err := page.Locator(".realtime-content > .ant-row > .even-line > .has-padding > div").All()
	if err != nil {
		log.Fatal("Error:", err)
	}

	production, err := entry[0].TextContent()
	production = strings.ReplaceAll(production, " kW", "")
	production = strings.ReplaceAll(production, " ", "")
	production = strings.ReplaceAll(production, "\u202f", "")
	production = strings.TrimSpace(production)
	production = strings.ReplaceAll(production, ",", ".")
	productionFloat, err := strconv.ParseFloat(production, 64)
	if err != nil {
		log.Fatal("Error: convert string to float", production, err)
	}

	return productionFloat
}

func Init() {
	debug, err := settings.Get("DEBUG_MODE")
	if err != nil {
		log.Fatal("Error:", err)
	}
	err = playwright.Install()

	if err != nil {
		log.Fatal("Error:", err)
	}

	pw, err := playwright.Run()

	if err != nil {
		log.Fatal("Error:", err)
	}

	browserType, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(!debug.(bool)),
	})

	if err != nil {
		log.Fatal("Error:", err)
	}

	context, err := browserType.NewContext()

	if err != nil {
		log.Fatal("Error:", err)
	}

	page, err = context.NewPage()
	if err != nil {
		log.Fatal("Error:", err)
	}
}
