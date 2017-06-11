package main

import (
	"github.com/acomagu/chatroom-go-v2/chatroom"
	"fmt"
	"regexp"
)

var livingPlaces = map[string]string{}

func topics(userID string) []chatroom.Topic {
	welcomeTopic := WelcomeTopic{userID: userID}
	return []chatroom.Topic{welcomeTopic.responseTopic, handleResponseTopic, welcomeTopic.welcomeTopic, welcomeTopic.askTopic}
}

type WelcomeTopic struct {
	userID string
}

func (c WelcomeTopic) welcomeTopic(room chatroom.Room) chatroom.DidTalk {
	<-room.In
	_, ok := livingPlaces[c.userID]
	if ok {
		return false
	}
	sendText(room, "こんにちは!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	sendText(room, "あなたが今住んでいる地域はどこですか?")
	msg, err := waitTextMsg(room)
	if err != nil {
		fmt.Println(err)
		return true
	}
	livingPlaces[c.userID] = msg
	sendText(room, "OK! その地域についてなにか質問するときもあると思うから、そのときは協力してくれると嬉しいな!")
	sendText(room, "その代わり「郡山でおすすめの肉屋は?」みたいに使ってみてね!")
	return true
}

func (c WelcomeTopic) askTopic(room chatroom.Room) chatroom.DidTalk {
	msg, err := waitTextMsg(room)
	if err != nil {
		return false
	}

	var place, purpose string
	if matches := regexp.MustCompile(`(.*)で(.*)ない`).FindStringSubmatch(msg); len(matches) >= 3 {
		place = matches[1]
		purpose = matches[2]
	} else if matches = regexp.MustCompile(`(.*)の(.*)`).FindStringSubmatch(msg); len(matches) >= 3 {
		place = matches[1]
		purpose = matches[2]
	} else {
		return false
	}

	sendText(room, fmt.Sprintf("わかりました! %sの%sを探してみますね!", place, purpose))

	for id, cr := range crs {
		if livingPlace, ok := livingPlaces[id]; ok && livingPlace == place {
			cr.In <- newRecommendationRequest(place, purpose, c.userID)
		}
	}
	return true
}

func (c WelcomeTopic) responseTopic(room chatroom.Room) chatroom.DidTalk {
	rr, ok := (<-room.In).(RecommendationRequest)
	if !ok {
		return false
	}

	sendText(room, fmt.Sprintf("%sでおすすめの%sはありますか?", rr.place, rr.purpose))
	msg, err := waitTextMsg(room)
	if err != nil {
		fmt.Println(err)
		return false
	}

	crs[rr.from].In <- newResponse(msg)
	sendText(room, "ありがとうございます!!")
	return true
}

func handleResponseTopic(room chatroom.Room) chatroom.DidTalk {
	if msg, ok := (<-room.In).(Response); ok {
		sendText(room, msg.string())
		return true
	}
	return false
}

type Response string

func newResponse(msg string) Response {
	return Response(msg)
}

func (r Response) string() string {
	return string(r)
}

type RecommendationRequest struct {
	place string
	purpose string
	from string
}

func newRecommendationRequest(place string, purpose string, userID string) RecommendationRequest {
	return RecommendationRequest{
		place: place,
		purpose: purpose,
		from: userID,
	}
}
