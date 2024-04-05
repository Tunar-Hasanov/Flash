package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

const (
	telegramToken     = "TELEGRAM_BOT_TOKEN'ınız"
	stellarHorizonURL = "https://horizon-testnet.stellar.org" // Stellar testnet için kullanılıyor
)

func main() {
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Bot çalışıyor. Çıkmak için Ctrl+C'ye basın.")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

			switch update.Message.Command() {
			case "start":
				msg.Text = "Stellar Transfer Bot'una hoş geldiniz!"
			case "help":
				msg.Text = "Komutlar:\n/sendxlm <adres> <miktar> - Belirtilen adrese XLM transferi yap"
			case "sendxlm":
				args := update.Message.CommandArguments()
				parts := strings.Fields(args)
				if len(parts) != 2 {
					msg.Text = "Kullanım: /sendxlm <adres> <miktar>"
				} else {
					adres := parts[0]
					miktarStr := parts[1]
					miktar, err := strconv.ParseFloat(miktarStr, 64)
					if err != nil {
						msg.Text = "Geçersiz miktar. Lütfen geçerli bir sayı girin."
					} else {
						err := sendXLM(adres, miktar)
						if err != nil {
							msg.Text = "XLM transferi sırasında hata oluştu: " + err.Error()
						} else {
							msg.Text = "XLM başarıyla transfer edildi!"
						}
					}
				}
			default:
				msg.Text = "Bilinmeyen komut. Kullanılabilir komutlar için /help yazın."
			}

			_, err := bot.Send(msg)
			if err != nil {
				log.Println("Mesaj gönderirken hata oluştu:", err)
			}
		}
	}
}

func sendXLM(adres string, miktar float64) error {
	// Stellar Horizon client oluştur
	client := horizon.DefaultPublicNetClient

	// Gönderen Stellar hesap oluştur
	kp := keypair.MustRandom()

	// Alıcı adresi ve miktarı belirt
	destination := adres
	amount := strconv.FormatFloat(miktar, 'f', -1, 64)

	// Transaction oluştur
	tx, err := horizon.BuildTransaction(client, kp.Address(), destination, amount)
	if err != nil {
		return err
	}

	// Transaction submit et
	resp, err := horizon.SubmitTransaction(client, tx, kp)
	if err != nil {
		return err
	}

	log.Printf("XLM başarıyla transfer edildi, TX ID: %s", resp.ID)

	return nil
}
