package slack

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"time"
)

type JSONTimeString string

// String converts the unix timestamp into a string
func (t JSONTimeString) String() string {
	if t == "" {
		return ""
	}
	floatN, err := strconv.ParseFloat(string(t), 64)
	if err != nil {
		log.Panicln(err)
		return ""
	}
	timeStr := int64(floatN)
	tm := time.Unix(int64(timeStr), 0)
	return fmt.Sprintf("\"%s\"", tm.Format("Mon Jan _2"))
}

var portMapping = map[string]string{"ws": "80", "wss": "443"}

func websocketizeUrlPort(orig string) (string, error) {
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
