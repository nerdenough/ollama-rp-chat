package chat

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jmorganca/ollama/api"
	log "github.com/sirupsen/logrus"
)

type Character struct {
	Name    string
	Persona string
}

type Chat struct {
	BotCharacter   Character
	UserCharacters []Character
	Scenario       string
}

func New() *Chat {
	return &Chat{}
}

func (c *Chat) GetCompletion(ctx context.Context, ollama *api.Client, request *api.GenerateRequest) (*api.GenerateResponse, error) {
	var res *api.GenerateResponse

	shouldRetry := true
	for retries := 0; shouldRetry && retries < 5; retries++ {
		str := ""
		err := ollama.Generate(ctx, request, func(response api.GenerateResponse) error {
			if response.Done {
				str += "\n"
				res = &response
				return nil
			}
			str += response.Response
			return nil
		})
		if err != nil {
			if retries >= 5 {
				return nil, err
			}
			log.Error(err)
			log.Info("retrying")
			continue
		}

		color.Yellow(fmt.Sprintf("%s: %s", c.BotCharacter.Name, strings.Trim(str, " \t\n")))
		fmt.Println()
		shouldRetry = false
	}

	return res, nil
}

func (c *Chat) GetCharacterInput(sc *bufio.Scanner, uc Character) string {
	characterInput := fmt.Sprintf("%s: ", uc.Name)
	fmt.Printf(">>> %s", color.GreenString(characterInput))
	sc.Scan()

	userInput := sc.Text()
	if userInput == "" {
		return userInput
	}
	return fmt.Sprintf("%s%s\n", characterInput, userInput)
}

func (c *Chat) GetCharacterInputs(sc *bufio.Scanner) string {
	prompt := ""
	for _, uc := range c.UserCharacters {
		prompt += c.GetCharacterInput(sc, uc)
	}
	fmt.Println()
	return prompt
}

func (c *Chat) StopTokens() []string {
	tokens := []string{"\n", "\nUser:", "\n###"}
	for _, uc := range c.UserCharacters {
		tokens = append(tokens, fmt.Sprintf("\n%s:", uc.Name))
	}
	return tokens
}

func (c *Chat) SystemPrompt() string {
	charactersPrompt := ""
	characterNames := ""
	for i, uc := range c.UserCharacters {
		characterNames += uc.Name
		if i < len(c.UserCharacters)-1 {
			characterNames += ", "
		}
		charactersPrompt += fmt.Sprintf("- %s: %s\n", uc.Name, uc.Persona)
	}

	return fmt.Sprintf(`
		Characters:
		- %s: %s
		%s

		Scenario: %s.

		### Instruction:
		You must play the role of %s.

		You must engage in a roleplaying chat with %s below this line. Do not write dialogue or narration for %s.
	`, c.BotCharacter.Name, c.BotCharacter.Persona, charactersPrompt, c.Scenario, c.BotCharacter.Name, characterNames, characterNames)
}

func (c *Chat) Template() string {
	return fmt.Sprintf(`
		{{- if .First }}
		### System:
		{{ .System }}
		{{- end }}

		### Input:
		{{ .Prompt }}

		### Response: (length = short)
		%s:
	`, c.BotCharacter.Name)
}
