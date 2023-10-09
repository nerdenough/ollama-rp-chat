# Ollama RP Chat

This is an attempt to create a tool for interacting with RP models intuitively.
It allows you to define multiple human played characters, and one character
played by the AI.

Currently as a CLI tool, to be extended soon with support for Discord bots,
where each user can opt in and define a character to play.

## Requirements

- Ollama
- For now, a model that follows the LimaRP-Alpaca prompt format, e.g. [Llama 2
  LimaRP chat v2][0]

## Usage

```sh
OLLAMA_MODEL=llama2 OLLAMA_HOST=http://localhost:11434 go run main.go
```

## Disclaimer

This is a tool that uses user-specified external language models - it does not
produce content by itself. External language models may show biases, generate
false information, or produce offensive material. This tool is only intended for
exploring model capabilities and their ability to roleplay a fictional
character.

[0]: https://huggingface.co/TheBloke/llama-2-13B-chat-limarp-v2-merged-GGUF
