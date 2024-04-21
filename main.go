package main

import (
	"TG_bot_new/clients/telgram"
	"flag"
	"fmt"
	"log"
)

const tgBotHost = "api.telegram.org"

func main() {
	t := mustToken()
	tgClient := telgram.New(tgBotHost, t)

	fmt.Println(tgClient)
}

func mustToken() string {
	token := flag.String(
		"t",
		"",
		"token for my test bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is required")
	}

	return *token
}
