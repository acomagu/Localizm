package main

import (
	"github.com/acomagu/chatroom-go-v2/chatroom"
	"github.com/line/line-bot-sdk-go/linebot"
)

func sendText(room chatroom.Room, msg string) {
	room.Out <- linebot.NewTextMessage(msg)
}

func waitTextMsg(room chatroom.Room) (string, error) {
	for {
		text, ok := pickText(waitEvent(room))
		if ok {
			return text, nil
		}
	}
}

func pickText(event *linebot.Event) (string, bool) {
	if event.Type != linebot.EventTypeMessage {
		return "", false
	}
	eventMsg, ok := event.Message.(*linebot.TextMessage)
	if !ok {
		return "", false
	}
	return eventMsg.Text, true
}

func waitEvent(room chatroom.Room) *linebot.Event {
	for {
		if ev, ok := pickEvent(<-room.In); ok {
			return ev
		}
	}
}

func pickEvent(msg interface{}) (*linebot.Event, bool) {
	if event, ok := msg.(*linebot.Event); ok {
		return event, true
	}
	return nil, false
}
