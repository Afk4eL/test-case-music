package error_struct

type Error struct {
	Err string `json:"Error"`
	Msg string `json:"Message,omitempty"`
}
