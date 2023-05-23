// Copyright © 2023 aerth
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"log"
	"os"
	"plugin"
	"strings"

	"github.com/aerth/botrunner"
	"gitlab.com/aerth/tgbot"
	"go.etcd.io/bbolt"
)

func loadPlugins(bot *tgbot.Bot, db *bbolt.DB, paths string) ([]botrunner.Plugin, error) {
	var (
		pluginPaths []string
	)

	switch {
	case strings.Contains(paths, ","):
		pluginPaths = strings.Split(paths, ",") // files not dirs
	case !strings.Contains(paths, ","):
		dir, err := os.ReadDir(paths) // try dir
		if err != nil {
			pluginPaths = append(pluginPaths, paths) // not dir, single file ?
			break
		}
		for _, e := range dir {
			name := e.Name()
			if !strings.HasPrefix(name, ".") && !e.IsDir() && e.Type().IsRegular() {
				pluginPaths = append(pluginPaths, name)
			}
		}
	}

	var (
		plugins = make([]botrunner.Plugin, len(pluginPaths))
	)

	if len(pluginPaths) == 0 || pluginPaths[0] == "" {
		log.Fatal("need at least one plugin")
	}
	for i, pluginPath := range pluginPaths {
		pluginPath = os.ExpandEnv(pluginPath)
		p, err := plugin.Open(pluginPath)
		if err != nil {
			log.Fatal(err)
		}
		sym, err := p.Lookup("Load") // Load()
		if err != nil {
			log.Fatal(err)
		}
		fn, ok := sym.(botrunner.InitFunc)
		if !ok {
			log.Fatalf("bad module: %T", sym)
		}
		pl, err := fn(bot, db)
		if err != nil {
			return nil, err
		}
		plugins[i] = pl
		if err := botrunner.RegisterBotCommands(bot, pl); err != nil {
			return nil, err
		}
	}
	return plugins, nil
}
