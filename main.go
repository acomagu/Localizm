package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/acomagu/chatroom-go-v2/chatroom"
	"github.com/line/line-bot-sdk-go/linebot"
)

var (
	lineChannelSecret      = os.Getenv("LINE_CHANNEL_SECRET")
	lineChannelAccessToken = os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	port                   = os.Getenv("PORT")
)

type chatrooms map[string]chatroom.Chatroom

var crs = make(chatrooms)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(lineChannelSecret, lineChannelAccessToken)
	if err != nil {
		fmt.Println(err)
	}
	http.HandleFunc("/", handleRequest)
	http.Handle("/resource/", http.StripPrefix("/resource/", http.FileServer(http.Dir("resource"))))
	fmt.Println(http.ListenAndServe(":"+port, nil))
}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	events, err := bot.ParseRequest(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, event := range events {
		if event == nil {
			continue
		}
		handleMsg(event)
	}
}

func handleMsg(event *linebot.Event) {
	userID := event.Source.UserID
	cr, ok := crs[userID]
	if !ok {
		cr = chatroom.New(topics(userID))
		crs[userID] = cr
		go sender(userID, cr)
	}
	cr.In <- event
}

func sender(userID string, cr chatroom.Chatroom) {
	for {
		msg := <-cr.Out
		if _msg, ok := msg.(linebot.Message); ok {
			_, err := bot.PushMessage(userID, _msg).Do()
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("message type must be `linebot.Message`")
		}
	}
}
