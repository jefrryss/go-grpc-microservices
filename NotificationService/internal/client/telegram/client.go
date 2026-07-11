package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	token      string
	httpClient *http.Client
}

func New(token string) *Client {
	return &Client{token: token, httpClient: &http.Client{Timeout: 35 * time.Second}}
}

func (c *Client) Send(ctx context.Context, target, message string) error {
	payload, err := json.Marshal(map[string]string{"chat_id": target, "text": message})
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint("sendMessage"), bytes.NewReader(payload))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("send Telegram message: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("Telegram API returned %s", response.Status)
	}
	return nil
}

func (c *Client) Poll(ctx context.Context) error {
	var offset int64
	for ctx.Err() == nil {
		updates, err := c.updates(ctx, offset)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		for _, update := range updates {
			offset = update.UpdateID + 1
			if update.Message.Text == "/start" {
				if err := c.Send(ctx, strconv.FormatInt(update.Message.Chat.ID, 10), "Бот подключён. Здесь будут уведомления о заказах."); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (c *Client) updates(ctx context.Context, offset int64) ([]update, error) {
	query := url.Values{"timeout": {"30"}, "offset": {strconv.FormatInt(offset, 10)}}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint("getUpdates")+"?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("get Telegram updates: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("Telegram API returned %s", response.Status)
	}
	var result struct {
		OK     bool     `json:"ok"`
		Result []update `json:"result"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}
	if !result.OK {
		return nil, fmt.Errorf("Telegram API rejected getUpdates")
	}
	return result.Result, nil
}

func (c *Client) endpoint(method string) string {
	return "https://api.telegram.org/bot" + c.token + "/" + method
}

type update struct {
	UpdateID int64 `json:"update_id"`
	Message  struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
	} `json:"message"`
}
