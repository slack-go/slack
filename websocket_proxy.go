package slack

import (
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

// Taken and reworked from: https://gist.github.com/madmo/8548738
func websocketHTTPConnect(proxy, urlString string) (net.Conn, error) {
	p, err := net.Dial("tcp", proxy)
	if err != nil {
		return nil, err
	}

	turl, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	req := http.Request{
		Method: "CONNECT",
		URL:    &url.URL{},
		Host:   turl.Host,
	}

	cc := httputil.NewProxyClientConn(p, nil)
	if _, err := cc.Do(&req); err != nil {
		return nil, err
	}

	rwc, _ := cc.Hijack()

	return rwc, nil
}

func websocketProxyDial(urlString, origin string) (ws *websocket.Conn, err error) {
	config, err := websocket.NewConfig(urlString, origin)
	if err != nil {
		return nil, err
	}
	config.Dialer = &net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}

	if os.Getenv("HTTP_PROXY") == "" {
		return websocket.DialConfig(config)
	}

	purl, err := url.Parse(os.Getenv("HTTP_PROXY"))
	if err != nil {
		return nil, err
	}

	client, err := websocketHTTPConnect(purl.Host, urlString)
	if err != nil {
		return nil, err
	}

	switch config.Location.Scheme {
	case "ws":
	case "wss":
		tlsClient := tls.Client(client, &tls.Config{
			ServerName: strings.Split(config.Location.Host, ":")[0],
		})
		err := tlsClient.Handshake()
		if err != nil {
			tlsClient.Close()
			return nil, err
		}
		client = tlsClient

	default:
		return nil, errors.New("invalid websocket schema")
	}

	return websocket.NewClient(config, client)
}
