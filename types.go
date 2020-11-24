package main

// Status represents SteamStatus' response
type Status struct {
	// Time represents when this status update was collected
	Time int64 `json:"time"`
	// Online represents the percentage of online services
	Online float32 `json:"online"`
	// Services represents each service's status
	Services [][]interface{} `json:"services"`
}

// Remap represents a remapping of Status that's better for templates
type Remap struct {
	Time     int64
	Online   float32
	Statuses map[string]struct {
		Name   string
		Good   string
		Status string
	}
}
