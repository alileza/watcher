package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	path    string
	words   string
	webhook string
	channel string
)

func main() {
	flag.StringVar(&path, "path", "", "file path")
	flag.StringVar(&words, "words", "", "words to capture")
	flag.StringVar(&webhook, "webhook", "", "webhook url")
	flag.StringVar(&channel, "channel", "", "webhook channel")
	flag.Parse()

	cmds := []string{"tail", "-f", path}
	fmt.Println("Executing : ", cmds)

	cmd := exec.Command(cmds[0], cmds[1:]...)

	w := newWatcher()

	cmd.Stdout = w
	cmd.Stderr = w

	if errno := cmd.Run(); errno != nil {
		panic(errno)
	}
}

func notify(msg string) {
	if webhook == "" {
		return
	}

	payload := map[string]interface{}{
		"text": msg,
	}
	if channel != "" {
		payload["channel"] = channel
	}
	b, _ := json.Marshal(payload)

	resp, err := http.DefaultClient.Post(webhook, "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Printf("[watcher] error on capturing error : %s \n", err.Error())
	}

	if resp.StatusCode >= 300 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[watcher] error on capturing error : %s \n", err)
			return
		}
		log.Printf("[watcher] error on capturing error : %s \n", string(b))
	}
}

type watcher struct {
	out chan []byte
}

func newWatcher() *watcher {
	w := &watcher{
		out: make(chan []byte),
	}

	return w.Run()
}

func (c *watcher) Run() *watcher {
	go func() {
		hostname, _ := os.Hostname()
		for {
			input := strings.ToLower(
				strings.TrimSpace(string(<-c.out)),
			)
			words := strings.Split(strings.ToLower(words), ",")

			for _, word := range words {
				if match, _ := regexp.MatchString(word, input); match {
					notify(fmt.Sprintf("capture `(%s:%s)` on `%s`", input, word, hostname))
				}
			}
		}
	}()
	return c
}

func (c *watcher) Write(b []byte) (int, error) {
	c.out <- b
	return len(b), nil
}
