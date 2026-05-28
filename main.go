package main

import (
	"log"
	"os"
	"strings"
	"unicode"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func resolveBotToken() string {
	for _, key := range []string{"token", "TOKEN", "BOT_TOKEN", "TELEGRAM_BOT_TOKEN"} {
		if v, ok := os.LookupEnv(key); ok {
			v = strings.TrimSpace(v)
			if v != "" {
				return v
			}
		}
	}

	data, err := os.ReadFile(".env")
	if err != nil {
		return ""
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		v = strings.Trim(v, `"'`)

		if k == "TOKEN" && v != "" {
			return v
		}
	}

	return ""
}

func main() {
	token := resolveBotToken()

	if token == "" {
		log.Fatal("token do bot nao encontrado")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("erro ao criar bot: %v", err)
	}

	if _, err := bot.Request(tgbotapi.DeleteWebhookConfig{
		DropPendingUpdates: true,
	}); err != nil {
		log.Printf("erro ao remover webhook: %v", err)
	}

	if err := registerBotCommands(bot); err != nil {
		log.Printf("erro ao registrar comandos: %v", err)
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	log.Printf("bot iniciado como @%s", bot.Self.UserName)

	for update := range updates {

		if update.CallbackQuery != nil {
			handleCallback(bot, update.CallbackQuery)
			continue
		}

		if update.Message == nil {
			continue
		}

		text := strings.TrimSpace(update.Message.Text)

		if text == "" {
			continue
		}

		cpfDigits := onlyDigits(text)

		if len(cpfDigits) == 11 {

			valid := isCPF(text)

			if valid {
				sendReply(
					bot,
					update.Message.Chat.ID,
					"✅ CPF <b>valido</b>.",
					buildMainMenu(),
				)
			} else {
				sendReply(
					bot,
					update.Message.Chat.ID,
					"❌ CPF <b>invalido</b>.",
					buildMainMenu(),
				)
			}

			continue
		}

		cmd := commandName(update.Message)

		var response string
		var markup interface{} = buildMainMenu()

		switch {

		case cmd == "start":
			response = "👋 Bem-vindo ao bot validador de CPF."

		case isValidateRequest(text) || cmd == "validar":
			response = "🔎 Envie um CPF para validar."

		case isHelpRequest(text) || cmd == "ajuda":
			response = "🧭 Centro de ajuda."
			markup = buildInlineMenu()

		case isAboutRequest(text) || cmd == "sobre":
			response = "📘 Bot de validação de CPF."
			markup = buildInlineMenu()

		default:
			response = "❌ Mensagem invalida."
		}

		sendReply(
			bot,
			update.Message.Chat.ID,
			response,
			markup,
		)
	}
}