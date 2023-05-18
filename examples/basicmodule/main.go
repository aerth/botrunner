// Copyright © 2023 aerth
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import "C"
import (
	"fmt"
	"log"

	"github.com/aerth/botrunner"
	"gitlab.com/aerth/tgbot"
	"go.etcd.io/bbolt"
)

// check type
var _ botrunner.InitFunc = Load
var _ botrunner.PreloadPlugin = BasicPlugin{}

func main() { println("plugin") }

func Load(tg *tgbot.Bot, db *bbolt.DB) (botrunner.Plugin, error) {
	if db == nil {
		return nil, fmt.Errorf("need db")
	}
	err := db.Update(func(tx *bbolt.Tx) error {
		return nil
	})
	if err != nil {
		return nil, err
	}
	return BasicPlugin{
		db: db,
		tg: tg,
	}, nil
}

type BasicPlugin struct {
	db *bbolt.DB
	tg *tgbot.Bot
}

func (b BasicPlugin) Name() string {
	return "Basic Example Plugin"
}
func (b BasicPlugin) ProcessInput(s string) (string, error) {
	if s == "version" {
		return "0.0.1", nil
	}
	if s == "help" {
		return "version", nil
	}
	return "it works", nil
}
func (b BasicPlugin) Commands() []string {
	return []string{}
}
func (b BasicPlugin) HandleUpdate(bot *tgbot.Bot, u tgbot.Update) error {
	// do nothing
	return nil
}
func (b BasicPlugin) Preload(bot *tgbot.Bot, u tgbot.Update) error {
	// if preload exists, this method will be called for *every* message (even non-text).
	if u.Message == nil || u.Message.Text != "" {
		return nil
	}
	// logs all input
	log.Printf("%s: %s\n", u.Message.From.String(), u.Message.Text)
	return nil

}
