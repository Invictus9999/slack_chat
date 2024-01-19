package membership

import "net/http"

type SubscribeRequest struct {
	SubscriberId   string `json:"subscriberId"`
	SubscribedToId string `json:"subscribedToId"`
}

func (s *SubscribeRequest) Bind(r *http.Request) error {
	return nil
}

type SubscribeResponse struct {
	SubscriptionId string `json:"subscriptionId"`
}

func (s *SubscribeResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
