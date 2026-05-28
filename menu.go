func normalizeMenuText(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))

	if value == "" {
		return ""
	}

	var b strings.Builder

	space := false

	for _, r := range value {

		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			space = false
			continue
		}

		if !space {
			b.WriteByte(' ')
			space = true
		}
	}

	return strings.TrimSpace(b.String())
}

func isValidateRequest(text string) bool {
	switch normalizeMenuText(text) {

	case "validar", "validar cpf":
		return true

	default:
		return false
	}
}

func isHelpRequest(text string) bool {
	switch normalizeMenuText(text) {

	case "ajuda", "help":
		return true

	default:
		return false
	}
}

func isAboutRequest(text string) bool {
	switch normalizeMenuText(text) {

	case "sobre", "about":
		return true

	default:
		return false
	}
}

func commandName(msg *tgbotapi.Message) string {

	if msg == nil {
		return ""
	}

	if msg.IsCommand() {
		return strings.ToLower(strings.TrimSpace(msg.Command()))
	}

	text := strings.TrimSpace(msg.Text)

	if !strings.HasPrefix(text, "/") {
		return ""
	}

	part := strings.Fields(text)[0]

	part = strings.TrimPrefix(part, "/")

	if idx := strings.IndexByte(part, '@'); idx >= 0 {
		part = part[:idx]
	}

	return strings.ToLower(part)
}

func registerBotCommands(bot *tgbotapi.BotAPI) error {

	cfg := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "start",
			Description: "Iniciar o bot",
		},
		tgbotapi.BotCommand{
			Command:     "validar",
			Description: "Validar CPF",
		},
		tgbotapi.BotCommand{
			Command:     "ajuda",
			Description: "Ajuda",
		},
		tgbotapi.BotCommand{
			Command:     "sobre",
			Description: "Sobre",
		},
	)

	_, err := bot.Request(cfg)

	return err
}

func sendReply(
	bot *tgbotapi.BotAPI,
	chatID int64,
	text string,
	markup interface{},
) {

	msg := tgbotapi.NewMessage(chatID, text)

	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = markup

	_, err := bot.Send(msg)

	if err != nil {
		log.Printf("erro ao enviar mensagem: %v", err)
	}
}

func buildMainMenu() tgbotapi.ReplyKeyboardMarkup {

	keyboard := tgbotapi.NewReplyKeyboard(

		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("✅ Validar CPF"),
		),

		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("❓ Ajuda"),
			tgbotapi.NewKeyboardButton("ℹ️ Sobre"),
		),
	)

	keyboard.ResizeKeyboard = true

	return keyboard
}

func buildInlineMenu() tgbotapi.InlineKeyboardMarkup {

	return tgbotapi.NewInlineKeyboardMarkup(

		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"📋 Como validar",
				"help_validate",
			),
		),

		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🤖 Sobre",
				"about_bot",
			),
		),
	)
}

func handleCallback(
	bot *tgbotapi.BotAPI,
	callback *tgbotapi.CallbackQuery,
) {

	var text string

	switch callback.Data {

	case "help_validate":
		text = "Envie um CPF com ou sem pontuação."

	case "about_bot":
		text = "Bot validador de CPF."

	default:
		text = "Opcao invalida."
	}

	sendReply(
		bot,
		callback.Message.Chat.ID,
		text,
		buildMainMenu(),
	)
}