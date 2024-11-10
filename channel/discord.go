package channel

import (
	"bytes"
	"encoding/json"
	"errors"
	"message-pusher/model"
	"net/http"
	"strings"
)

type discordMessageRequest struct {
	Content string `json:"content"`
}

type discordMessageResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SendDiscordMessage(message *model.Message, user *model.User, channel_ *model.Channel) error {
	if message.Content == "" {
		message.Content = message.Description
	}
	messageRequest := discordMessageRequest{
		Content: message.Content,
	}
	// https://discord.com/developers/docs/reference#message-formatting
	if message.To != "" {
		messageRequest.Content = ""
		ids := strings.Split(message.To, "|")
		for _, id := range ids {
			messageRequest.Content = "<@" + id + "> " + messageRequest.Content
		}
		messageRequest.Content = messageRequest.Content + message.Content
	}

	jsonData, err := json.Marshal(messageRequest)
	if err != nil {
		return err
	}
	resp, err := http.Post(channel_.URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	var res discordMessageResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return err
	}
	if res.Code != 0 {
		return errors.New(res.Message)
	}
	if resp.StatusCode == http.StatusBadRequest {
		return errors.New(resp.Status)
	}
	return nil
}
