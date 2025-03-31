package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TgErrorBody struct {
	Gateway string
	Source  string
	ErrText string
	Time    string
	Request interface{}
}

func SendInfo(messageBody TgErrorBody) {
	client := http.Client{}
	d := struct {
		ChatId    string `json:"chat_id"`
		Text      string `json:"text"`
		ParseMode string `json:"parse_mode"`
	}{
		ChatId:    "-1002412508267",
		ParseMode: "HTML",
	}
	d.Text = fmt.Sprintf(`
	 <b>Gateway:</b> %s
	 <b>Source:</b> %s
	 <b>Error Text:</b> %s
	 <b>Time:</b> %s
	 
   `, "test-api", messageBody.Source, messageBody.ErrText, messageBody.Time)
	if messageBody.Request != nil {
		d.Text += fmt.Sprintf(`
	 <b>Request:</b> <pre>%+v</pre>`, messageBody.Request)
	}
	body, _ := json.Marshal(d)
	req, err := http.NewRequest("POST", "https://api.telegram.org/bot7849240216:AAFCLJwclWCdpI-1YJRCVHKlOvDczZx0w30/sendMessage", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("sendInfo:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = client.Do(req)
	if err != nil {
		fmt.Println("sendInfo:", err)
	}
}