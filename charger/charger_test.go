package charger_test

import (
	"reflect"
	"testing"

	"github.com/tchap/go-charger/charger"
)

func TestCharger_Charge(t *testing.T) {
	// Specify the configuration struct to charge.
	// This is used later, but it is good to have the idea.
	type MQTTConfig struct {
		ClientID string `charger:"CLIENT_ID"`
		Username string `charger:"USERNAME"`
		Password string `charger:"PASSWORD"`
	}

	type Config struct {
		ServiceName string      `charger:"SERVICE_NAME"`
		TaskSlot    int         `charger:"TASK_SLOT"`
		LogLevel    string      `charger:"LOG_LEVEL"`
		MQTT        *MQTTConfig `charger:"MQTT"`
	}

	// Configure a charger.
	main := charger.New()

	main.SetMapFunc(func(key string) (bool, string) {
		switch key {
		case "SERVICE_NAME":
			return "charger", true
		case "TASK_SLOT":
			return "2", true
		case "LOG_LEVEL":
			return "debug", true
		case "MQTT_PASSWORD":
			return "secret", true
		default:
			return "", false
		}
	})

	main.Add(charger.String{
		Name:     "SERVICE_NAME",
		Required: true,
	})

	main.Add(charger.Int{
		Name:    "TASK_SLOT",
		Default: 1,
	})

	main.Add(charger.String{
		Name:    "LOG_LEVEL",
		Default: "info",
	})

	mqtt := main.WithPrefix("MQTT_")

	mqtt.Add(charger.String{
		Name:     "CLIENT_ID",
		Template: `{{ get "SERVICE_NAME" }}.{{ get "TASK_SLOT" }}-RAND`,
		Required: true,
	})

	mqtt.Add(charger.String{
		Name:     "USERNAME",
		Template: `{{ get "SERVICE_NAME" }}`,
		Required: true,
	})

	mqtt.Add(charger.String{
		Name:     "PASSWORD",
		Required: true,
	})

	// Charge!
	var config Config
	if err := main.Charge(&config); err != nil {
		t.Fatal(err)
	}

	// Make sure the config struct is charged correctly.
	expected := Config{
		ServiceName: "charger",
		TaskSlot:    2,
		LogLevel:    "debug",
		MQTT: {
			ClientID: "charger.2-RAND",
			Username: "charger",
			Password: "secret",
		},
	}

	if !reflect.DeepEqual(&config, &expected) {
		t.Errorf("mismatch; expected = %+v got = %+v", expected, config)
	}
}
