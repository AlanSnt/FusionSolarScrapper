package settings

import (
	"errors"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"sync"
)

type Settings struct {
	FUSION_SOLAR_URL string
	SMART_PVM_NAME   string
	TIME_DELTA       int
	USERNAME         string
	PASSWORD         string
	DEBUG_MODE       bool
	MQTT_HOST        string
	MQTT_PORT        int
	MQTT_USER        string
	MQTT_PASS        string
}

var (
	instance *Settings
	once     sync.Once
)

func getEnv(key string, defaultValue ...string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	if len(defaultValue) == 0 {
		panic("Missing environment variable: " + key)
	}

	return defaultValue[0]
}

func Init() {
	time_delta := getEnv("TIME_DELTA", "30")
	time_delta_int, err := strconv.ParseInt(time_delta, 10, 64)
	if err != nil {
		log.Fatal("Error:", err)
	}

	mqttPort := getEnv("MQTT_PORT", "1883")
	mqttPortInt, err := strconv.ParseInt(mqttPort, 10, 64)
	if err != nil {
		log.Fatal("Error:", err)
	}

	once.Do(func() {
		instance = &Settings{
			FUSION_SOLAR_URL: getEnv("FUSION_SOLAR_URL"),
			SMART_PVM_NAME:   getEnv("SMART_PVM_NAME"),
			TIME_DELTA:       int(time_delta_int),
			USERNAME:         getEnv("USERNAME"),
			PASSWORD:         getEnv("PASSWORD"),
			DEBUG_MODE:       getEnv("DEBUG_MODE", "false") == "true",
			MQTT_HOST:        getEnv("MQTT_HOST"),
			MQTT_PORT:        int(mqttPortInt),
			MQTT_USER:        getEnv("MQTT_USER"),
			MQTT_PASS:        getEnv("MQTT_PASS"),
		}

		if !instance.DEBUG_MODE {
			log.SetOutput(io.Discard)
		}

		log.Println("Settings initialized")
	})
}

func Get(key string) (interface{}, error) {
	if instance == nil {
		return nil, errors.New("Settings not initialized")
	}

	val := reflect.ValueOf(instance).Elem()
	exist := val.FieldByName(key).IsValid()

	if !exist {
		return nil, errors.New("Setting not found")
	}

	return val.FieldByName(key).Interface(), nil
}
