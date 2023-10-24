package ceye

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"xyz/config"
)

type Result struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data []struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Remote_addr string `json:"remote_addr"`
		Created_at  string `json:"created_at"`
	} `json:"data"`
}

// client *http.Request
func Get_Ceye_RS(client *http.Client) (data []string, err error) {
	conf, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	//http://api.ceye.io/v1/records?token={token}&type={dns|http}&filter={filter}
	ceye_api := fmt.Sprintf("http://%v/v1/records?token=%v&type=%v&filter=%v", conf.Ceye_Api_Config.Address, conf.Ceye_Api_Config.Token, conf.Ceye_Api_Config.Type, conf.Ceye_Api_Config.Filter)
	resp, err := http.NewRequest("GET", ceye_api, nil)
	if err != nil {
		return nil, err
	}

	resp.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MSIE 6.0; Windows NT 5.0; Trident/4.0)")

	res, err := client.Do(resp)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var rs = Result{}
	err = json.Unmarshal(body, &rs)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if rs.Meta.Code == 200 && len(rs.Data) != 0 {
		for _, v := range rs.Data {
			data = append(data, v.Name)
		}
	}

	return data, nil
}
