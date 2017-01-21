package common

import (
    "encoding/json"
    "os"
    "log"
)

type Configuration struct {
    TelegramAPI string
}

func LoadConfiguration() *Configuration {
    file, err := os.Open("config/config.json")
    if err != nil {
        log.Printf("error, invalid config file");
    }
    decoder := json.NewDecoder(file)
    configuration := Configuration{}
    err = decoder.Decode(&configuration)
    if err != nil {
        log.Printf("error:", err)
    }
    token := os.Getenv("BOT_TOKEN")
    if token != "" {
        configuration.TelegramAPI = token
    }
    log.Printf("the key is %s", configuration.TelegramAPI)
    return &configuration
}
