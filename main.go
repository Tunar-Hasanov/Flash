package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/tronprotocol/go-tron/api"
	"github.com/tronprotocol/go-tron/common"
)

const (
	telegramToken = "YOUR_TELEGRAM_BOT_TOKEN"
	usdtContractAddress = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t" // USDT contract address on Tron mainnet
)

var tronClient *api.GrpcClient

func main() {
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}

	tronClient = api.NewGrpcClient("https://api.trongrid.io")

	log.Printf("Bot is running. Press Ctrl+C to exit.")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			switch update.Message.Command() {
			case "start":
				msg.Text = "Welcome to USDT Transfer Bot!"
			case "help":
				msg.Text = "Commands:\n/sendusdt <address> <amount> - Send USDT to the specified address"
			case "sendusdt":
				args := update.Message.CommandArguments()
				parts := strings.Fields(args)
				if len(parts) != 2 {
					msg.Text = "Usage: /sendusdt <address> <amount>"
				} else {
					address := parts[0]
					amount, err := strconv.ParseFloat(parts[1], 64)
					if err != nil {
						msg.Text = "Invalid amount. Please enter a valid number."
					} else {
						err := sendUSDT(bot, update.Message.Chat.ID, address, amount)
						if err != nil {
							msg.Text = "Error sending USDT: " + err.Error()
						} else {
							msg.Text = "USDT sent successfully!"
						}
					}
				}
			default:
				msg.Text = "Unknown command. Type /help for available commands."
			}

			_, err := bot.Send(msg)
			if err != nil {
				log.Println("Error sending message:", err)
			}
		}
	}
}

func sendUSDT(bot *tgbotapi.BotAPI, chatID int64, address string, amount float64) error {
	// Prepare USDT transfer
	fromAddress := tronClient.GetDefaultAccount()
	toAddress := common.HexToAddress(address)
	assetName := common.USDT
	amountInSun := int64(amount * 1000000) // Convert amount to SUN (1 USDT = 10^6 SUN)

	// Execute USDT transfer
	txID, err := tronClient.SendAsset(fromAddress, toAddress, assetName, amountInSun)
	if err != nil {
		return err
	}

	log.Printf("USDT sent: %s", txID)

	return nil
}
