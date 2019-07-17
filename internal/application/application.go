package application

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Bot      *tgbotapi.BotAPI
}

func NewApplication() *Application {
	dir, _ := os.Getwd()
	logsPath := filepath.Join(dir, "logs")
	if err := os.Mkdir(logsPath, 0700); err != nil && !os.IsExist(err) {
		log.Fatalln("error creating logs directory")
	}
	// init info logger
	infoLog := newLogger("info", logsPath)
	// init error logger
	errorLog := newLogger("error", logsPath)

	app := &Application{
		InfoLog:  infoLog,
		ErrorLog: errorLog,
	}
	// setup proxy for telegram bot connection
	client := setupProxyClient()
	bot, err := tgbotapi.NewBotAPIWithClient(os.Getenv("BOT_ID"), client)
	if err != nil {
		app.ErrorLog.Fatalf("Error while connecting to telegram bot: %s", err)
	}
	app.Bot = bot
	app.InfoLog.Printf("Authorized on account %s", bot.Self.UserName)

	return app
}

func newLogger(prefix, logsPath string) *log.Logger {
	loggerFile, err := os.OpenFile(filepath.Join(logsPath, fmt.Sprintf("%s.logs.txt", prefix)),
		os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		log.Fatalf("error while creating file for logging %s: %v", prefix, err)
	}

	logPrefix := fmt.Sprintf("%s\t", strings.ToUpper(prefix))
	return log.New(loggerFile, logPrefix, log.Ldate|log.Ltime)
}

// setup proxy from package to avoid restrictions
func setupProxyClient() *http.Client {
	proxyStr := os.Getenv("PROXY_URL")
	proxyUrl, err := url.Parse(proxyStr)
	if err != nil {
		log.Fatalf("error parsing proxy %s: %v", proxyStr, err)
	}
	transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

	return &http.Client{
		Transport: transport,
	}
}

func (app *Application) InfoPrintF(format string, v ...interface{}) {
	app.InfoLog.Printf(format, v...)
}

func (app *Application) ErrorPrintF(format string, v ...interface{}) {
	app.ErrorLog.Printf(format, v...)
}

func (app *Application) BotSend(c tgbotapi.Chattable) (msg tgbotapi.Message, err error) {
	msg, err = app.Bot.Send(c)
	return
}
func (app *Application) BotGetFile(fileID string) (file tgbotapi.File, err error) {
	file, err = app.Bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	return
}

func (app *Application) BotClientGet(link string) (response *http.Response, err error) {
	response, err = app.Bot.Client.Get(link)
	return
}
