package inspire

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type quoteStruct struct {
	Quote  string `json:"q"`
	Author string `json:"a"`
	H      string `json:"h"`
}

func GetQuote() string {

	url := "https://zenquotes.io/api/random"

	lizBot := http.Client{
		Timeout: time.Second * 2,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Discord Bot")
	res, getErr := lizBot.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	/*text :=
	textBytes := []byte(text)*/

	quote1 := []quoteStruct{}
	err = json.Unmarshal(body, &quote1)
	if err != nil {
		fmt.Println(err)
		return "Error with getting quote!"
	}

	return (quote1[0].Quote + " - " + quote1[0].Author)
}
