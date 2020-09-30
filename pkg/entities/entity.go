package entities

type Entity struct {
	Path         string   `json:"Path" yaml:"Path"`
	Name         string   `json:"Name" yaml:"Name"`
	Destination  string   `json:"Destination" yaml:"Destination"`
	Type         string   `json:"Type" yaml:"Type"`
	Recurse      bool     `json:"Recurse" yaml:"Recurse"`
	CloneOptions []Option `json:"CloneOptions" yaml:"CloneOptions"`
}
