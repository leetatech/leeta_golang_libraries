package states

type State struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Region  string   `json:"region"`
	Capital string   `json:"capital"`
	Lgas    []string `json:"lgas,omitempty"`
	Slogan  string   `json:"slogan"`
	Towns   []string `json:"towns,omitempty"`
} // @name State
