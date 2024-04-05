package main

import (
	"log"
	"net/http"
	"bytes"
	"encoding/json"
)

const (
	telegramToken = "TELEGRAM_BOT_TOKEN'ınız"
	usdtContractAddress = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t" // Tron mainnet üzerindeki USDT sözleşme adresi
	tronGridAPIKey = "TRONGRID_API_KEY'iniz"
)

type TransferRequest struct {
	To     string `json:"to"`
	Value  int64  `json:"value"`
	TokenID string `json:"tokenID"`
}

func main() {
	// Telegram botunu başlat
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
				msg.Text = "USDT Transfer Bot'una hoş geldiniz!"
			case "help":
				msg.Text = "Komutlar:\n/sendusdt <adres> <miktar> - Belirtilen adrese USDT gönder"
			case "sendusdt":
				args := update.Message.CommandArguments()
				parts := strings.Fields(args)
				if len(parts) != 2 {
					msg.Text = "Kullanım: /sendusdt <adres> <miktar>"
				} else {
					adres := parts[0]
					miktarStr := parts[1]
					miktar, err := strconv.ParseFloat(miktarStr, 64)
					if err != nil {
						msg.Text = "Geçersiz miktar. Lütfen geçerli bir sayı girin."
					} else {
						err := sendUSDT(bot, update.Message.Chat.ID, adres, miktar)
						if err != nil {
							msg.Text = "USDT gönderilirken hata oluştu: " + err.Error()
						} else {
							msg.Text = "USDT başarıyla gönderildi!"
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

func sendUSDT(bot *tgbotapi.BotAPI, chatID int64, adres string, miktar float64) error {
	// USDT transfer isteğini hazırla
	transferReq := TransferRequest{
		To:     adres,
		Value:  int64(miktar * 1000000), // Miktarı SUN cinsinden belirt (1 USDT = 10^6 SUN)
		TokenID: usdtContractAddress,
	}

	// İsteği JSON'a dönüştür
	reqBody, err := json.Marshal(transferReq)
	if err != nil {
		return err
	}

	// TronGrid API'sine transfer isteği gönder
	apiURL := "https://api.trongrid.io/wallet/transferasset"
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("TRON-PRO-API-KEY", tronGridAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("TronGrid API isteği başarısız: %s", resp.Status)
	}

	log.Printf("USDT başarıyla gönderildi, Adres: %s", adres)

	return nil
}
