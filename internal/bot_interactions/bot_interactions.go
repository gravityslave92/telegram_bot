package bot_interactions

import (
	"bufio"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"telegram_bot/internal/application"
	"unicode/utf8"
)

const botID = "784720809:AAGQBCIdvrtzbCLW2pxwHt1j0N93bUiMlfU"

type botApplication interface {
	InfoPrintF(string, ...interface{})
	ErrorPrintF(string, ...interface{})
	BotSend(tgbotapi.Chattable) (tgbotapi.Message, error)
	BotGetFile(string) (tgbotapi.File, error)
	BotClientGet(string) (*http.Response, error)
}

// exported for main()
func StartBotChat() {
	app := application.NewApplication()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := app.Bot.GetUpdatesChan(u)
	if err != nil {
		app.ErrorLog.Fatalf("error while connecting to channel updates: %v", err)
	}
	for update := range updates {
		userMessage := update.Message
		if userMessage == nil { // ignore any non-Message Updates
			continue
		}

		msg := tgbotapi.NewMessage(userMessage.Chat.ID, "")
		if userMessage.IsCommand() {
			userCmd := userMessage.Command()
			app.InfoLog.Printf("%s command has been requested!", userCmd)

			switch userCmd {
			case "help":
				msg.Text = "type /checkUrlList or /status"
			case "status":
				msg.Text = "I'm http.Ok for 200%)"
			case "checkUrlList":
				msg.Text = "please upload your url list"
				app.BotSend(msg)
				processUserRequest(app, updates, &msg)
			default:
				msg.Text = userMessage.Text
			}
		}

		app.InfoLog.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg.ReplyToMessageID = userMessage.MessageID
		app.BotSend(msg)
	}
}

func processUserRequest(app botApplication, updates tgbotapi.UpdatesChannel, msg *tgbotapi.MessageConfig) {
	urlListMsg := <-updates
	if document := urlListMsg.Message.Document; document != nil && document.FileID != "" {
		urlsCh, err := downloadTelegramFile(app, document.FileID)
		if err != nil {
			app.ErrorPrintF("error occured while downloading %s file from telegram %v", document.FileName, err)
			msg.Text = err.Error()

			return
		}

		msg.Text = "please provide goroutines limit"
		app.InfoPrintF("Url list has been provided! Responded with message: %s", msg.Text)
		app.BotSend(msg)

		limitMsg := <-updates
		limit, err := parseLimitFromMsg(limitMsg)
		if err != nil {
			app.ErrorPrintF("error while parsing limit value from message: %v", err)
			msg.Text = err.Error()

			return
		}

		resultCh := make(chan string)
		go processUrlsWithLimit(app, urlsCh, resultCh, limit)

		msg.Text = buildResultMsg(resultCh)
		app.InfoPrintF("Message has been completed! Result is: %s ", msg.Text)

		return
	}

	msg.Text = "url list must be provided"
	app.ErrorPrintF("Failed to complete response! Replied with message: %s", msg.Text)
}

func downloadTelegramFile(app botApplication, fileID string) (<-chan string, error) {
	ch := make(chan string)
	file, err := app.BotGetFile(fileID)
	if err != nil {
		return nil, err
	}
	fileLink := file.Link(os.Getenv("BOT_ID"))

	resp, err := app.BotClientGet(fileLink)
	if err != nil {
		return nil, fmt.Errorf("failed to download file from telegram chat")
	}

	go func() {
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan(); len(scanner.Bytes()) != 0; scanner.Scan() {
			app.InfoPrintF("%s received as argument for downloading", scanner.Text())
			ch <- scanner.Text()
		}

		resp.Body.Close()
		close(ch)
	}()

	return ch, nil
}

func processUrlsWithLimit(app botApplication, urlsCh <-chan string, resultCh chan<- string, limit int) {
	var wg sync.WaitGroup
	semaphore := make(chan int, limit)

	for url := range urlsCh {
		// enter semaphore
		semaphore <- 1
		wg.Add(1)

		go func(urlString string) {
			app.InfoPrintF("Starting of gathering response from %s", url)

			defer func() {
				// release semaphore
				<-semaphore
				wg.Done()
				app.InfoPrintF("Successfully completed gathering response from %s", url)
			}()

			resp, err := http.Get(urlString)
			if err != nil {
				app.ErrorPrintF("error while downloading response from %s: %v", urlString, err)
				return
			}
			defer resp.Body.Close()

			bytesResp, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				app.ErrorPrintF("error occured while reading response from %s: %v", urlString, err)
				return
			}
			// send result string to output chan
			resultCh <- fmt.Sprintf("%s: %d\n", urlString, utf8.RuneCount(bytesResp))

			return
		}(url)
	}

	// waiting for all info to become aggregated
	wg.Wait()
	close(resultCh)
}

func parseLimitFromMsg(limitMsg tgbotapi.Update) (int, error) {
	limitStr := limitMsg.Message.Text
	if limitStr == "" {
		return 0, fmt.Errorf("invalid value provided! must be of type int")
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, err
	}

	return limit, nil
}

func buildResultMsg(ch <-chan string) string {
	var builder strings.Builder
	for result := range ch {
		builder.WriteString(result)
	}

	return builder.String()
}
