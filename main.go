package main

import (
	"html"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	var botToken = os.Getenv("BOTTOKEN")

	b, err := tb.NewBot(tb.Settings{
		Token:  botToken,
		Poller: &tb.LongPoller{Timeout: 5 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}
	b.Handle(tb.OnText, makeExecHandler(b))
	b.Start()
}

func makeExecHandler(b *tb.Bot) func(m *tb.Message) {
	return func(m *tb.Message) {
		arr := strings.Split(strings.TrimSpace(m.Text), " ")
		if len(arr) == 0 {
			b.Send(m.Chat, "No commands")
			return
		}

		name := arr[0]

		var args []string

		if len(arr) > 1 {
			for _, a := range arr[1:] {
				a = strings.TrimSpace(a)
				if a != "" {
					args = append(args, a)
				}
			}
		}

		out, err := exec.Command(name, args...).Output()
		if err != nil {
			b.Send(m.Chat, "Erorr when exec command: "+name+": "+err.Error())
			return
		}
		if len(out) == 0 {
			out = []byte("ok")
		}
		b.Send(m.Chat, "<pre>"+html.EscapeString(string(out))+"</pre>", tb.ModeHTML)
	}
}
