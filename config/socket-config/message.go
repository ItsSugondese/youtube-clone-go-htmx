package socket_config

type Message struct {
	FromClientId  string `json:"fromClientId"`
	ClientID      string `json:"clientID"`
	Text          string `json:"text"`
	UseRedis      bool   `json:"useRedis"`
	OperationType string `json:"operationType"`
	Response      any    `json:"response"`
}
