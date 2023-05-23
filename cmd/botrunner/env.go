package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func getTokenFromEnvFile(filenames []string, keywords ...string) string {
	if len(keywords) == 0 {
		keywords = []string{"token"}
	}
	for i := range keywords {
		keywords[i] = strings.ToLower(keywords[i])
	}
	for _, filename := range filenames {
		filename = strings.TrimSpace(filename)
		if filename == "" {
			continue
		}
		_, err := os.Stat(filename)
		if err == os.ErrNotExist {
			continue
		}
		if err != nil {
			log.Println("error stat env file:", err)
			continue
		}
		f, err := os.Open(filename)
		if err != nil {
			continue
		}
		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "#") {
				continue
			}
			// split after =
			ss := strings.Split(line, "=")
			if len(ss) != 2 {
				// maybe even empty
				continue
			}
			key := strings.TrimSpace(ss[0])
			found := false
			for _, keyword := range keywords {
				if !strings.Contains(strings.ToLower(key), keyword) {
					// not token
					continue
				}
				found = true
			}
			if !found {
				// not in keyword
				continue
			}

			// check if token is tg format: TOKEN=1234567:AbcDef
			val := strings.TrimSpace(ss[1])
			splitVal := strings.Split(val, ":")
			if len(splitVal) != 2 {
				// not the tg token
				continue
			}
			n, _ := strconv.ParseUint(splitVal[0], 10, 64)
			if n == 0 {
				continue
			}
			if len(splitVal[1]) == 0 {
				continue
			}
			// got token
			return val
		}
	}
	return os.Getenv(strings.ToUpper(keywords[0])) // token -> $TOKEN
}
