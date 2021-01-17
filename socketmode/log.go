package socketmode

func (smc *Client) Debugf(msg string, args ...interface{}) {
	smc.apiClient.Debugf(msg, args...)
}

func (smc *Client) Debugln(v ...interface{}) {
	smc.apiClient.Debugln(v...)
}
