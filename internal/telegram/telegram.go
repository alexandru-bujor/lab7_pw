package telegram

import (
	"fmt"
	"log"
	"os"

	"MegaMobileBack/internal/stats"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI

const (
	ChatID       int64 = -1003415263352
	TopicService int   = 6
	TopicLombard int   = 4
	TopicOrders  int   = 8
	TopicStats   int   = 10
)

func StartBot() error {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		token = "8531734711:AAGlQ41cvou0sIWKrPIUde8CJG4ttrxHhbI"
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	Bot = bot
	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	go func() {
		for update := range updates {
			if update.Message != nil {
				handleMessage(update.Message)
			}
		}
	}()

	return nil
}

func handleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	text := message.Text
	chatType := message.Chat.Type
	chatTitle := message.Chat.Title
	messageID := message.MessageID

	if text == "/start" {
		handleStartCommand(message, chatID, chatTitle, chatType, messageID)
	} else if text == "statistica" || text == "/statistica" {
		if message.ReplyToMessage != nil && message.ReplyToMessage.MessageID == TopicStats {
			handleStatisticsCommand(message, chatID)
		} else if chatID == ChatID {
			handleStatisticsCommand(message, chatID)
		}
	}
}

func handleStartCommand(message *tgbotapi.Message, chatID int64, chatTitle, chatType string, messageID int) {
	var responseText string
	var topicID int

	if message.ReplyToMessage != nil {
		topicID = message.ReplyToMessage.MessageID
		responseText = fmt.Sprintf("📊 Chat Information:\n\n"+
			"Chat Name: %s\n"+
			"Chat ID: %d\n"+
			"Chat Type: %s\n\n"+
			"📌 Topic Information:\n"+
			"Topic ID: %d\n"+
			"Current Message ID: %d",
			chatTitle, chatID, chatType, topicID, messageID)

		if message.ReplyToMessage.Text != "" {
			responseText += fmt.Sprintf("\n\nTopic Starter Message:\n%s",
				truncateText(message.ReplyToMessage.Text, 100))
		} else if message.ReplyToMessage.Caption != "" {
			responseText += fmt.Sprintf("\n\nTopic Starter Caption:\n%s",
				truncateText(message.ReplyToMessage.Caption, 100))
		}
	} else {
		responseText = fmt.Sprintf("📊 Chat Information:\n\n"+
			"Chat Name: %s\n"+
			"Chat ID: %d\n"+
			"Chat Type: %s\n"+
			"Message ID: %d\n\n"+
			"ℹ️ This message is not in a topic thread.\n"+
			"To identify a topic, send /start from within a topic thread.",
			chatTitle, chatID, chatType, messageID)
	}

	msg := tgbotapi.NewMessage(chatID, responseText)
	if message.ReplyToMessage != nil {
		msg.ReplyToMessageID = messageID
	}

	if _, err := Bot.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func handleStatisticsCommand(message *tgbotapi.Message, chatID int64) {
	if Bot == nil {
		return
	}

	statsText := getDailyStatistics()

	msg := tgbotapi.NewMessage(chatID, statsText)
	if message.ReplyToMessage != nil {
		msg.ReplyToMessageID = message.MessageID
	}

	if _, err := Bot.Send(msg); err != nil {
		log.Printf("Error sending statistics: %v", err)
	}
}

func getDailyStatistics() string {
	statsSvc := &stats.Service{}
	today, yesterday, err := statsSvc.GetDailyStats()
	if err != nil {
		return fmt.Sprintf("❌ Eroare la încărcarea statisticilor: %v", err)
	}

	message := fmt.Sprintf("📊 Statistici Zilnice\n\n"+
		"📅 Astăzi (%s):\n"+
		"  • Cereri Service: %d\n"+
		"  • Cereri Lombard: %d\n"+
		"  • Comenzi: %d\n"+
		"  • Venituri: %.2f MDL\n"+
		"  • Clienți noi: %d\n\n"+
		"📅 Ieri (%s):\n"+
		"  • Cereri Service: %d\n"+
		"  • Cereri Lombard: %d\n"+
		"  • Comenzi: %d\n"+
		"  • Venituri: %.2f MDL\n"+
		"  • Clienți noi: %d",
		today.Date, today.ServiceRequests, today.LombardRequests, today.Orders, today.Revenue, today.NewCustomers,
		yesterday.Date, yesterday.ServiceRequests, yesterday.LombardRequests, yesterday.Orders, yesterday.Revenue, yesterday.NewCustomers)

	return message
}

func SendServiceRequestNotification(requestID int, customerName, customerPhone, deviceInfo, problemDesc string) {
	if Bot == nil {
		return
	}

	message := fmt.Sprintf("🔧 Nouă cerere de service!\n\n"+
		"ID: #%d\n"+
		"Client: %s\n"+
		"Telefon: %s\n"+
		"Dispozitiv: %s\n"+
		"Problemă: %s",
		requestID, customerName, customerPhone, deviceInfo, truncateText(problemDesc, 50))

	msg := tgbotapi.NewMessage(ChatID, message)
	msg.ReplyToMessageID = TopicService

	if _, err := Bot.Send(msg); err != nil {
		log.Printf("Error sending service request notification: %v", err)
	}
}

func SendLombardRequestNotification(requestID int, customerName, customerPhone, requestType string, itemsCount int) {
	if Bot == nil {
		return
	}

	requestTypeText := "Electronice"
	if requestType == "gold" {
		requestTypeText = "Aur"
	}

	message := fmt.Sprintf("💰 Nouă cerere lombard!\n\n"+
		"ID: #%d\n"+
		"Client: %s\n"+
		"Telefon: %s\n"+
		"Tip: %s\n"+
		"Produse: %d",
		requestID, customerName, customerPhone, requestTypeText, itemsCount)

	msg := tgbotapi.NewMessage(ChatID, message)
	msg.ReplyToMessageID = TopicLombard

	if _, err := Bot.Send(msg); err != nil {
		log.Printf("Error sending lombard request notification: %v", err)
	}
}

func SendOrderNotification(orderID int, customerName, customerPhone, productTitle string, totalAmount float64) {
	if Bot == nil {
		return
	}

	message := fmt.Sprintf("📦 Comandă nouă!\n\n"+
		"ID: #%d\n"+
		"Client: %s\n"+
		"Telefon: %s\n"+
		"Produs: %s\n"+
		"Total: %.2f MDL",
		orderID, customerName, customerPhone, truncateText(productTitle, 50), totalAmount)

	msg := tgbotapi.NewMessage(ChatID, message)
	msg.ReplyToMessageID = TopicOrders

	if _, err := Bot.Send(msg); err != nil {
		log.Printf("Error sending order notification: %v", err)
	}
}

func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}
