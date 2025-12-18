package sendsms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type UseCase struct{}

func NewUseCase() *UseCase {
	return &UseCase{}
}

func (u *UseCase) SendSMS(phone, code, token string) error {
	body := map[string]string{
		"phone_number": phone,
		"message":      fmt.Sprintf("Tasdiqlash kodi: %s", code),
		"from":         "OnlineMarket",
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "https://notify.eskiz.uz/api/message/sms/send", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil

}
