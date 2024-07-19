package v1

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func GetAvatar() ([]byte, error) {
	url := fmt.Sprintf("https://api.dicebear.com/9.x/pixel-art/svg?seed=%v", time.Now().Unix())

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received non-200 status code: %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return bodyBytes, nil
}
