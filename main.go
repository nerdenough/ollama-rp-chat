package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/jmorganca/ollama/api"
	_ "github.com/joho/godotenv/autoload"
	"github.com/nerdenough/ollama-chat/chat"
	log "github.com/sirupsen/logrus"
)

func getUserInput(sc *bufio.Scanner) string {
	fmt.Print(">>> ")
	sc.Scan()
	fmt.Println()
	return sc.Text()
}

func newChat(sc *bufio.Scanner) *chat.Chat {
	c := chat.New()
	color.Green("How many characters will you play? (1-4)")
	numUserCharacters, err := strconv.Atoi(getUserInput(sc))
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < numUserCharacters; i++ {
		uc := chat.Character{}
		color.Green(fmt.Sprintf("Name of your character: (%d/%d)", i+1, numUserCharacters))
		uc.Name = getUserInput(sc)
		color.Green(fmt.Sprintf("Short description of your character's persona: (%d/%d)", i+1, numUserCharacters))
		uc.Persona = getUserInput(sc)
		c.UserCharacters = append(c.UserCharacters, uc)
	}

	color.Green("Name of the AI's character:")
	c.BotCharacter.Name = getUserInput(sc)
	color.Green("Short description of the AI character's persona:")
	c.BotCharacter.Persona = getUserInput(sc)

	color.Green("Scenario:")
	c.Scenario = getUserInput(sc)
	return c
}

func main() {
	fmt.Println("Ollama Chat Roleplay CLI")
	fmt.Println()

	ollama, err := api.FromEnv()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	sc := bufio.NewScanner(os.Stdin)
	var c *chat.Chat

	if err := ollama.Heartbeat(ctx); err != nil {
		color.Red(fmt.Sprintf("Unable to access ollama at %s", api.Host()))
		os.Exit(1)
	}

	list, err := ollama.List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	modelName := os.Getenv("OLLAMA_MODEL")
	if len(strings.Split(modelName, ":")) == 1 {
		modelName += ":latest"
	}

	modelFound := false
	for _, model := range list.Models {
		if model.Name == modelName {
			modelFound = true
		}
	}
	if !modelFound {
		color.Red(fmt.Sprintf("Unable to find specified model: %s", modelName))
		os.Exit(1)
	}

	color.Cyan(fmt.Sprintf("Using model: %s\n", modelName))
	fmt.Println()
	color.Cyan("(1) Start a new roleplay")
	color.Cyan("(2) Load an existing roleplay")
	fmt.Println()

	opt := getUserInput(sc)
	if opt == "1" {
		c = newChat(sc)
	} else {
		log.Fatal("not implemented")
	}

	request := &api.GenerateRequest{
		Model: modelName,
		Options: map[string]interface{}{
			"stop":            c.StopTokens(),
			"num_predict":     128,
			"repeat_pentalty": 1.1,
			"temperature":     0.85,
			"tfs_z":           0.95,
			"top_k":           0,
			"top_p":           1,
		},
		Template: c.Template(),
		System:   c.SystemPrompt(),
	}

	prompt := c.GetCharacterInputs(sc)
	for prompt != "!exit" {
		request.Prompt = prompt
		res, err := c.GetCompletion(ctx, ollama, request)
		if err != nil {
			log.Fatal(err)
		}
		request.Context = res.Context
		prompt = c.GetCharacterInputs(sc)
	}
}
