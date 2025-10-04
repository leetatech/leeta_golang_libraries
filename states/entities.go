package states

type State struct {
	ID       string   `json:"id" bson:"id"`
	Name     string   `json:"name" bson:"name"`
	Region   string   `json:"region" bson:"region"`
	Capital  string   `json:"capital" bson:"capital"`
	Lgas     []string `json:"lgas,omitempty" bson:"-"`
	Slogan   string   `json:"slogan" bson:"slogan"`
	Towns    []string `json:"towns,omitempty" bson:"-"`
	StatusTS int64    `json:"status_ts" bson:"status_ts"`
	TS       int64    `json:"ts" bson:"ts"`
} // @name State
