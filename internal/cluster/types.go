package cluster

type Sample struct {
	TS  string `json:"ts,omitempty"`
	Raw string `json:"raw,omitempty"`
}

type Cluster struct {
	Count   int      `json:"count"`
	Repr    string   `json:"repr"`
	FirstTS string   `json:"first_ts,omitempty"`
	LastTS  string   `json:"last_ts,omitempty"`
	Samples []Sample `json:"samples,omitempty"`
}
