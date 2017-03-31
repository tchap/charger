package charger_test

import (
	"crypto/rand"
	"encoding/hex"
	"reflect"
	"strings"
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

	// Configure the main application charger.
	main := charger.New()

	main.SetMapFunc(func(key string) (string, error) {
		switch key {
		case "SERVICE_NAME":
			return "charger", nil
		case "TASK_SLOT":
			return "2", nil
		case "LOG_LEVEL":
			return "debug", nil
		default:
			return "", charger.ErrNotFound
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

	// Configure the MQTT subcharger.
	mqtt := main.WithPrefix("MQTT_")

	mqtt.SetMapFunc(func(key string) (string, error) {
		switch key {
		case "PASSWORD":
			return "secret", nil
		default:
			return "", charger.ErrNotFound
		}
	})

	var rnd string
	mqtt.AddTemplateFunc("rand", func(length uint) string {
		raw := make([]byte, (8*length)/16)
		if err := rand.Read(raw); err != nil {
			mqtt.Error(errors.Wrap(err), "rand template function failed")
			return
		}
		rnd = hex.EncodeToString(raw)
		return rnd
	})

	mqtt.Add(charger.String{
		Name:     "CLIENT_ID",
		Default:  `{{ get "SERVICE_NAME" }}.{{ get "TASK_SLOT" }}-{{ rand 6 }}`,
		Required: true,
	})

	mqtt.Add(charger.String{
		Name:     "USERNAME",
		Default:  `{{ get "SERVICE_NAME" }}`,
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
			ClientID: "charger.2-" + rnd,
			Username: "charger",
			Password: "secret",
		},
	}

	if !reflect.DeepEqual(&config, &expected) {
		t.Errorf("mismatch; expected = %+v got = %+v", expected, config)
	}
}
