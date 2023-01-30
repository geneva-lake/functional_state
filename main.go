package main

import (
	"net/http"

	"github.com/geneva-lake/functional_state/logger"
	"github.com/geneva-lake/functional_state/service"
)

func main() {
	config, err := service.NewConfig().FromFile("config.yaml").Yaml()
	if err != nil {
		logger.Log(logger.Error, "main", 0, err, nil, nil, nil)
		return
	}
	handlers := CreateHandlers(config)
	if err := http.ListenAndServe(":"+config.Port, handlers); err != nil {
		logger.Log(logger.Error, "main", 0, err, nil, nil, nil)
	}
}
