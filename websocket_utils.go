package slack

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"time"
)

// JSONTimeString is an auxiliary type to allow us to format the time as we wish
type JSONTimeString string

// String converts the unix timestamp into a string
func (t JSONTimeString) String() string {
	tm := t.Time()
	if tm.IsZero() {
		return ""
	}
	return fmt.Sprintf("\"%s\"", tm.Format("Mon Jan _2"))
}

// Time converts the timestamp string to time.Time
func (t JSONTimeString) Time() time.Time {
	if t == "" {
		return time.Time{}
	}
	floatN, err := strconv.ParseFloat(string(t), 64)
	if err != nil {
		log.Println("ERROR parsing a JSONTimeString!", err)
		return time.Time{}
	}
	return time.Unix(int64(floatN), 0)
}

var portMapping = map[string]string{"ws": "80", "wss": "443"}

func websocketizeURLPort(orig string) (string, error) {
	urlObj, err := url.ParseRequestURI(orig)
	if err != nil {
		return "", err
	}
	_, _, err = net.SplitHostPort(urlObj.Host)
	if err != nil {
		return urlObj.Scheme + "://" + urlObj.Host + ":" + portMapping[urlObj.Scheme] + urlObj.Path, nil
	}
	return orig, nil
}
