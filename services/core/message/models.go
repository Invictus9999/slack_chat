package message

import "net/http"

type SendRequest struct {
	Content    string `json:"content"`
	SenderId   string `json:"senderId"`
	ReceiverId string `json:"receiverId"`
}

func (s *SendRequest) Bind(r *http.Request) error {
	return nil
}

type SendResponse struct {
	MessageId string `json:"messageId"`
}

func (s *SendResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type PublishMessage struct {
	SenderId string `json:"from"`
	Content  string `json:"message"`
}
