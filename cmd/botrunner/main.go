// Copyright © 2023 aerth
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/aerth/botrunner"
	tgbotapi "gitlab.com/aerth/telegrambotapi"
	"gitlab.com/aerth/tgbot"
	"go.etcd.io/bbolt"
)

type InitConfig struct {
	DatabasePath string
	PluginPath   string
	BotToken     string
}

func or(s string, def string) string {
	got, ok := os.LookupEnv(s)
	if !ok {
		return def
	}
	return got
}
func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	var cfg = InitConfig{
		BotToken:     or("TOKEN", ""),
		DatabasePath: or("DATABASE", "database.db"),
		PluginPath:   or("PLUGIN", ""),
	}
	flag.StringVar(&cfg.BotToken, "tg", cfg.BotToken, "@botfather token")
	flag.StringVar(&cfg.DatabasePath, "db", cfg.DatabasePath, "path to db file (will be created)")
	flag.StringVar(&cfg.PluginPath, "plugin", cfg.PluginPath, "botrunner plugin path (can be comma seperated for multiple, or directory for all)")
	flag.Parse()
	if cfg.BotToken == "" {
		cfg.BotToken = getTokenFromEnvFile([]string{".env", os.ExpandEnv("$HOME/.env"), os.Getenv("ENV_FILE")}, "TOKEN")
	}
	if cfg.BotToken == "" {
		log.Println("need bot token from @botfather")
	}
	db, bot, _, err := initialize(cfg)
	if err != nil {
		log.Fatalln("fatal init:", err)
	}
	defer db.Close()
	log.Printf("telegram: @%s", bot.Bot().Self.UserName)
	if err := bot.Start(); err != nil {
		log.Fatalln("fatal tg:", err)
	}
}

func initialize(cfg InitConfig) (*bbolt.DB, *tgbot.Bot, []botrunner.Plugin, error) {
	db, err := bbolt.Open(cfg.DatabasePath, 0600, &bbolt.Options{
		Timeout: 3 * time.Second,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	bot, err := tgbot.NewBot(cfg.BotToken)
	if err != nil {
		return nil, nil, nil, err
	}
	plugins, err := loadPlugins(bot, db, cfg.PluginPath)
	if err != nil {
		return nil, nil, nil, err
	}

	return db, bot, plugins, nil
}

// mkreply to reply a messagein telegram handlers
func MkReply(b *tgbot.Bot, u tgbot.Update) func(text string, i ...any) (tgbotapi.Message, error) {
	fn := func(text string, i ...any) (tgbotapi.Message, error) {
		return b.ReplyText(u.Message.Chat.ID, u.Message.MessageID, text, i...)
	}
	return fn
}
