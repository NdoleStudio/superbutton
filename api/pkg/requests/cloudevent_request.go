package requests

import "time"

type CloudEvent struct {
	Specversion     string      `json:"specversion"`
	Type            string      `json:"type"`
	Source          string      `json:"source"`
	ID              string      `json:"id"`
	Time            time.Time   `json:"time"`
	Datacontenttype string      `json:"datacontenttype"`
	Data            interface{} `json:"data"`
}
