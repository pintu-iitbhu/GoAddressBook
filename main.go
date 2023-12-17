package main

import (
	"GoAddressBook/addressbook"
	"GoAddressBook/cli"
	"GoAddressBook/configs"
	"github.com/sagikazarmark/slog-shim"
)

func main() {
	configs.NewConfig()
	bookInstance := addressbook.NewAddressBook()
	err := bookInstance.LoadFromFile()
	if err != nil {
		slog.Info("failed to load data from json file : ", err)
		return
	}

	cliInstance, err := cli.NewCliInstance(bookInstance)
	if err != nil {
		slog.Info("Error while instancing command-line interface :", err)
		return
	}
	cliInstance.Menu()
	return
}
