package slack

// Config contains some config parameters needed
// Token always needs to be set for the api to function
// Origin and Protocol are optional and only needed for websocket
type Config struct {
	token    string
	origin   string
	protocol string
}
