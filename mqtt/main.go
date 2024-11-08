package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"

	"AlanSnt/FusionSolarScrapper/settings"
	"encoding/json"
	"log"
	"strconv"
)

type DiscoveryPayload struct {
	Name              string `json:"name"`
	StateTopic        string `json:"state_topic"`
	UnitOfMeasurement string `json:"unit_of_measurement"`
	DeviceClass       string `json:"device_class"`
	StateClass        string `json:"state_class"`
	UniqueID          string `json:"unique_id"`
}

var client mqtt.Client

func SendMessage(topicName string, value float64) {
	token := client.Publish("fusion_solar_scrapper/sensor/"+topicName, 0, false, strconv.FormatFloat(value, 'f', -1, 64))
	token.Wait()
	log.Println("Message sent:", topicName, strconv.FormatFloat(value, 'f', -1, 64))
}

func sendDiscoverMessage() {
	productionDiscovery := DiscoveryPayload{
		Name:              "Energie produite",
		StateTopic:        "fusion_solar_scrapper/sensor/production",
		UnitOfMeasurement: "kW",
		DeviceClass:       "energy",
		StateClass:        "measurement",
		UniqueID:          "energy_production_1",
	}
	prodPayload, _ := json.Marshal(productionDiscovery)
	client.Publish("homeassistant/sensor/fusion_solar_scrapper_production_energy/config", 0, true, prodPayload)

	returnedDiscovery := DiscoveryPayload{
		Name:              "Energie retournée",
		StateTopic:        "fusion_solar_scrapper/sensor/returned",
		UnitOfMeasurement: "kW",
		DeviceClass:       "energy",
		StateClass:        "measurement",
		UniqueID:          "returned_energy_1",
	}
	returnedPayload, _ := json.Marshal(returnedDiscovery)
	client.Publish("homeassistant/sensor/fusion_solar_scrapper_returned_energy/config", 0, true, returnedPayload)

	consumedDiscovery := DiscoveryPayload{
		Name:              "Energie consomée",
		StateTopic:        "fusion_solar_scrapper/sensor/consumed",
		UnitOfMeasurement: "kW",
		DeviceClass:       "energy",
		StateClass:        "measurement",
		UniqueID:          "consumed_energy_1",
	}
	consumedPayload, _ := json.Marshal(consumedDiscovery)
	client.Publish("homeassistant/sensor/fusion_solar_scrapper_consumed_energy/config", 0, true, consumedPayload)
}

func Init() {
	mqttHost, err := settings.Get("MQTT_HOST")
	if err != nil {
		log.Fatal("Error:", err)
	}

	mqttPort, err := settings.Get("MQTT_PORT")
	if err != nil {
		log.Fatal("Error:", err)
	}

	mqttUser, err := settings.Get("MQTT_USER")
	if err != nil {
		log.Fatal("Error:", err)
	}

	mqttPass, err := settings.Get("MQTT_PASS")
	if err != nil {
		log.Fatal("Error:", err)
	}

	opts := mqtt.NewClientOptions().AddBroker("tcp://" + mqttHost.(string) + ":" + strconv.Itoa(mqttPort.(int)))
	opts.SetUsername(mqttUser.(string))
	opts.SetPassword(mqttPass.(string))
	opts.SetClientID("fusion_solar_scrapper")

	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	sendDiscoverMessage()
}
