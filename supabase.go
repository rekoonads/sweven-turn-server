package edgeturn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

)

type SupabaseAgent struct {
	url string
	anon_key string
}

func NewSupabaseAgent(url string, 
					  key string, 
					  ) (*SupabaseAgent) {
	return &SupabaseAgent{
		url: url,
		anon_key: key,
	}
}


func (agent *SupabaseAgent)	Ping( uid string)( err error)  {
	body,_ := json.Marshal( struct {
		AccountID string `json:"account_uid"`
	}{
		AccountID: uid,
	})

	req,err := http.NewRequest("POST",fmt.Sprintf("%s/rest/v1/rpc/ping_account",agent.url),bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Profile", "public")
	req.Header.Set("Content-Profile", "public")
	req.Header.Set("Authorization", "Bearer "+agent.anon_key)
	req.Header.Set("apikey", agent.anon_key)
	resp,err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		data,_ := io.ReadAll(resp.Body)
		return fmt.Errorf(string(data))
	}

	return
}