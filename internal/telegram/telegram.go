package telegram

import (
	"fmt"
	"log"
	"os"

	"MegaMobileBack/internal/stats"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Bot *tgbotapi.BotAPI

// Topic IDs for the MegaMobile group
const (
	ChatID       int64 = -1003415263352
	TopicService int   = 6  // service topic
	TopicLombard int   = 4  // Lombard topic
	TopicOrders  int   = 8  // comenzi topic
	TopicStats   int   = 10 // statistica topic
)

// StartBot initializes and starts the Telegram bot
func StartBot() error {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		// Use the provided token as default
		token = "8531734711:AAGlQ41cvou0sIWKrPIUde8CJG4ttrxHhbI"
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	Bot = bot
	bot.Debug = false

	log.Printf("🤖 Telegram bot authorized on account %s", bot.Self.UserName)

	// Set up update configuration
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Handle incoming updates
	go func() {
		for update := range updates {
			if update.Message != nil {
				handleMessage(update.Message)
			}
		}
	}()

	log.Println("✅ Telegram bot is running and listening for messages")
	return nil
}

// handleMessage processes incoming messages
func handleMessage(message *tgbotapi.Message) {
	chatID := message.Chat.ID
	text := message.Text
	chatType := message.Chat.Type
	chatTitle := message.Chat.Title
	messageID := message.MessageID

	// Log message details
	log.Printf("📨 Received message from chat: %s (ID: %d), chat type: %s, message ID: %d, text: %s",
		chatTitle, chatID, chatType, messageID, text)

	// Check if this is a reply (which can help identify topics)
	if message.ReplyToMessage != nil {
		log.Printf("   ↳ This is a reply to message ID: %d", message.ReplyToMessage.MessageID)
	}

	// Handle /start command
	if text == "/start" {
		handleStartCommand(message, chatID, chatTitle, chatType, messageID)
	} else if text == "statistica" || text == "/statistica" {
		// Handle statistics command in statistica topic
		if message.ReplyToMessage != nil && message.ReplyToMessage.MessageID == TopicStats {
			handleStatisticsCommand(message, chatID)
		} else if chatID == ChatID {
			// Check if we're in the stats topic by checking if it's a reply to the stats topic
			handleStatisticsCommand(message, chatID)
		}
	}
}

// handleStartCommand handles the /start command
func handleStartCommand(message *tgbotapi.Message, chatID int64, chatTitle, chatType string, messageID int) {
	var responseText string
	var topicID int

	// Check if this is in a topic (reply to a topic starter message)
	if message.ReplyToMessage != nil {
		// The Reply To Message ID is the topic ID
		topicID = message.ReplyToMessage.MessageID

		responseText = fmt.Sprintf("📊 Chat Information:\n\n"+
			"Chat Name: %s\n"+
			"Chat ID: %d\n"+
			"Chat Type: %s\n\n"+
			"📌 Topic Information:\n"+
			"Topic ID: %d\n"+
			"Current Message ID: %d",
			chatTitle, chatID, chatType, topicID, messageID)

		// Try to get more info about the topic starter message
		if message.ReplyToMessage.Text != "" {
			responseText += fmt.Sprintf("\n\nTopic Starter Message:\n%s",
				truncateText(message.ReplyToMessage.Text, 100))
		} else if message.ReplyToMessage.Caption != "" {
			responseText += fmt.Sprintf("\n\nTopic Starter Caption:\n%s",
				truncateText(message.ReplyToMessage.Caption, 100))
		}

		log.Printf("📌 Topic identified: Chat=%d, TopicID=%d", chatID, topicID)
	} else {
		// Not in a topic
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

	// If this was a reply, reply in the same thread
	if message.ReplyToMessage != nil {
		msg.ReplyToMessageID = messageID
	}

	if _, err := Bot.Send(msg); err != nil {
		log.Printf("❌ Error sending message: %v", err)
	} else {
		if topicID != 0 {
			log.Printf("✅ Sent topic info: chat=%d, topic=%d", chatID, topicID)
		} else {
			log.Printf("✅ Sent chat info: chat=%d", chatID)
		}
	}
}

// handleStatisticsCommand handles the statistica command
func handleStatisticsCommand(message *tgbotapi.Message, chatID int64) {
	if Bot == nil {
		return
	}

	statsText := getDailyStatistics()

	msg := tgbotapi.NewMessage(chatID, statsText)

	// Reply in the same topic if it's a reply
	if message.ReplyToMessage != nil {
		msg.ReplyToMessageID = message.MessageID
	}

	if _, err := Bot.Send(msg); err != nil {
		log.Printf("❌ Error sending statistics: %v", err)
	} else {
		log.Printf("✅ Sent statistics to chat: %d", chatID)
	}
}

// getDailyStatistics fetches and formats daily statistics
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

// SendServiceRequestNotification sends a notification about a new service request
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
		requestID, customerName, customerPhone, deviceInfo, truncateText(problemDesc, 100))

	msg := tgbotapi.NewMessage(ChatID, message)
	// Reply to the service topic starter message to post in the correct topic
	msg.ReplyToMessageID = TopicService

	if _, err := Bot.Send(msg); err != nil {
		log.Printf("❌ Error sending service request notification: %v", err)
	} else {
		log.Printf("✅ Sent service request notification: #%d to topic %d", requestID, TopicService)
	}
}

// SendLombardRequestNotification sends a notification about a new lombard request
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
	// Reply to the lombard topic starter message to post in the correct topic
	msg.ReplyToMessageID = TopicLombard

	if _, err := Bot.Send(msg); err != nil {
		log.Printf("❌ Error sending lombard request notification: %v", err)
	} else {
		log.Printf("✅ Sent lombard request notification: #%d to topic %d", requestID, TopicLombard)
	}
}

// SendOrderNotification sends a notification about a new order
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
		orderID, customerName, customerPhone, truncateText(productTitle, 100), totalAmount)

	msg := tgbotapi.NewMessage(ChatID, message)
	// Reply to the orders topic starter message to post in the correct topic
	msg.ReplyToMessageID = TopicOrders

	if _, err := Bot.Send(msg); err != nil {
		log.Printf("❌ Error sending order notification: %v", err)
	} else {
		log.Printf("✅ Sent order notification: #%d to topic %d", orderID, TopicOrders)
	}
}

// truncateText truncates text to a maximum length
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}
