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
	_, err := page.Reload(playwright.PageReloadOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})

	if err != nil {
		log.Fatal("(Login) Error:", err)
	}

	username, err := settings.Get("USERNAME")
	if err != nil {
		log.Fatal("(Login) Error:", err)
	}

	password, err := settings.Get("PASSWORD")
	if err != nil {
		log.Fatal("(Login) Error:", err)
	}

	loginUserNameEntry, err := page.Locator("[name='ssoCredentials.username']").All()
	if err != nil {
		log.Fatal("(Login) Error:", err)
	}

	if len(loginUserNameEntry) < 2 {
		return
	}

	loginPasswordEntry, err := page.Locator("[name='ssoCredentials.value']").All()
	if err != nil {
		log.Fatal("(Login) Error:", err)
	}

	loginUserNameEntry[1].Fill(username.(string))
	loginPasswordEntry[1].Fill(password.(string))

	loginButton, err := page.Locator(".loginBtn").All()
	if err != nil {
		log.Fatal("(Login) Error:", err)
	}

	loginButton[0].Click()
}

func GetReturnedEnergy() float64 {
	url, err := settings.Get("FUSION_SOLAR_URL")
	smartPvnName, err := settings.Get("SMART_PVM_NAME")
	if err != nil {
		log.Fatal("(GetReturnedEnergy) Error:", err)
	}

	log.Print("Get ReturnedEnergy")
	_, err = page.Goto(url.(string), playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	if err != nil {
		log.Fatal("(GetReturnedEnergy) Error on load page :", err)
	}

	login()

	page.WaitForSelector(".nco-monitor-station-overview-management")
	err = page.Locator(".flex-node-line-expand-part").Last().Click()
	if err != nil {
		log.Fatal("(GetReturnedEnergy) Error:", err)
	}

	page.WaitForSelector("[title='" + smartPvnName.(string) + "']")
	err = page.Locator(".flex-node-line-expand-part").Last().Click()
	if err != nil {
		log.Fatal("(GetReturnedEnergy) Error:", err)
	}

	page.WaitForSelector("[title='Meter-1']")
	err = page.GetByTitle("Meter-1").First().Click()
	if err != nil {
		log.Fatal("(GetReturnedEnergy) Error:", err)
	}

	page.WaitForSelector(".realtime-content")
	page.WaitForTimeout(5000)

	entry, err := page.Locator(".realtime-content > .ant-row > .even-line > .has-padding > div").All()
	if err != nil {
		log.Fatal("(GetReturnedEnergy) Error:", err)
	}

	returnedEnergy, _ := entry[1].TextContent()
	returnedEnergy = strings.ReplaceAll(returnedEnergy, " W", "")
	returnedEnergy = strings.TrimSpace(returnedEnergy)
	returnedEnergy = strings.ReplaceAll(returnedEnergy, " ", "")
	returnedEnergy = strings.ReplaceAll(returnedEnergy, "\u202f", "")
	returnedEnergy = strings.ReplaceAll(returnedEnergy, ",", ".")

	returnedEnergyFloat, err := strconv.ParseFloat(returnedEnergy, 64)
	if err != nil {
		log.Fatal("(GetReturnedEnergy) Error: convert string to float", returnedEnergy, err)
	}

	return returnedEnergyFloat
}

func GetSolarData() float64 {
	url, err := settings.Get("FUSION_SOLAR_URL")
	smartPvnName, err := settings.Get("SMART_PVM_NAME")
	if err != nil {
		log.Fatal("Error:", err)
	}

	log.Print("Get SolarData")
	_, err = page.Goto(url.(string), playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	if err != nil {
		log.Fatal("(GetSolarData) Error on load page :", err)
	}

	login()
	page.WaitForSelector(".nco-monitor-station-overview-management")

	err = page.Locator(".flex-node-line-expand-part").Last().Click()
	if err != nil {
		log.Fatal("(GetSolarData) Error:", err)
	}

	page.WaitForSelector("[title='" + smartPvnName.(string) + "']")

	err = page.GetByTitle(smartPvnName.(string)).First().Click()
	if err != nil {
		log.Fatal("(GetSolarData) Error:", err)
	}

	page.WaitForSelector(".realtime-content")
	page.WaitForTimeout(5000)
	entry, err := page.Locator(".realtime-content > .ant-row > .even-line > .has-padding > div").All()
	if err != nil {
		log.Fatal("(GetSolarData) Error:", err)
	}

	production, _ := entry[0].TextContent()
	production = strings.ReplaceAll(production, " kW", "")
	production = strings.ReplaceAll(production, " ", "")
	production = strings.ReplaceAll(production, "\u202f", "")
	production = strings.TrimSpace(production)
	production = strings.ReplaceAll(production, ",", ".")
	productionFloat, err := strconv.ParseFloat(production, 64)
	if err != nil {
		log.Fatal("(GetSolarData) Error: convert string to float", production, err)
	}

	return productionFloat
}

func Init() {
	debug, err := settings.Get("DEBUG_MODE")
	if err != nil {
		log.Fatal("Error:", err)
	}

	err = playwright.Install(&playwright.RunOptions{
		Browsers: []string{"firefox"},
	})
	if err != nil {
		log.Fatal("Error:", err)
	}

	pw, err := playwright.Run()
	if err != nil {
		log.Fatal("Error:", err)
	}

	browserType, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(!debug.(bool)),
		SlowMo:   playwright.Float(0.0),
		Args: []string{
			"--no-sandbox",
			"--disable-setuid-sandbox",
			"--window-size=1920,1080",
		},
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

	page.SetExtraHTTPHeaders(map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:88.0) Gecko/20100101 Firefox/88.0",
		"Accept-Language": "fr,fr-FR;q=0.8,en-US;q=0.5,en;q=0.3",
	})

	page.On("error", func(err error) {
		log.Printf("(Headless browser) Error on: %s\n", err.Error())
	})

	if debug.(bool) {
		page.On("console", func(msg playwright.ConsoleMessage) {
			log.Printf("(Headless browser) Console log: %s\n", msg.Text())
		})

		page.On("request", func(request playwright.Request) {
			log.Printf("(Headless browser) Request: %s %s\n", request.Method(), request.URL())
		})
		page.On("response", func(response playwright.Response) {
			log.Printf("(Headless browser) Response: %d %s\n", response.Status(), response.URL())
		})
	}
}
