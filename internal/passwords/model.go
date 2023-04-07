package passwords

type Password struct {
	DeviceID string `json:"deviceId,omitempty"`
	ID       string `json:"id,omitempty"`
	Password string `json:"password,omitempty"`
}
