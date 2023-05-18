// Copyright © 2023 aerth
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package botrunner

import (
	"fmt"
	"log"

	"gitlab.com/aerth/tgbot"
	"go.etcd.io/bbolt"
)

// InitFunc must be named "Load"
//
// check type:
//
//	var _ botrunner.InitFunc = Load
//	var _ botrunner.PreloadPlugin = MyPlugin{}
type InitFunc = func(tgbot *tgbot.Bot, db *bbolt.DB) (Plugin, error)

type Plugin interface {
	Name() string
	Commands() []string
	HandleUpdate(b *tgbot.Bot, u tgbot.Update) error
}
type PreloadPlugin interface {
	Plugin
	Preload(b *tgbot.Bot, u tgbot.Update) error
}

// RegisterBotCommands adds a plugin's available commands (and preload function if exists)
//
// Each plugin can replace previously loaded commands. (If two plugins provide the same command, the last-loaded takes precedence)
//
// If a Preload() function exists, it is stacked.
func RegisterBotCommands(b *tgbot.Bot, p Plugin) error {
	switch p := p.(type) {
	case PreloadPlugin:
		log.Println("registering preload:", p.Name())
		b.Preload(p.Preload)
	case Plugin:
	default:
		return fmt.Errorf("invalid plugin type: %T", p)
	}
	for _, cmd := range p.Commands() {
		b.Handle(cmd, p.HandleUpdate) // register bot commands
	}
	return nil
}
