package consensus

type Command struct {
	Operation string `json:"operation,omitempty"`
	Key       string `json:"key,omitempty"`
	Value     any    `json:"value,omitempty"`
}
