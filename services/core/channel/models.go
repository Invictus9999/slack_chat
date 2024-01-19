package channel

import "net/http"

type CreateChannelRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (s *CreateChannelRequest) Bind(r *http.Request) error {
	return nil
}

type CreateChannelResponse struct {
	Id string `json:"id"`
}

func (s *CreateChannelResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
