package message_broker

type QueueBinding struct {
	Source          string                 `json:"source"`
	VHost           string                 `json:"vhost"`
	Destination     string                 `json:"destination"`
	DestinationType string                 `json:"destination_type"`
	RoutingKey      string                 `json:"routing_key"`
	Arguments       map[string]interface{} `json:"arguments"`
	PropertiesKey   string                 `json:"properties_key"`
}
