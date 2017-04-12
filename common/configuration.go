package common

import (
    "encoding/json"
    "os"
    "log"
)

type Configuration struct {
    TelegramAPI string
    FolderName string
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
    // if BOT_TOKEN env variable is setted it will overwrite the one defined in the configurations
    token := os.Getenv("BOT_TOKEN")
    if token != "" {
        configuration.TelegramAPI = token
    }
    return &configuration
}
