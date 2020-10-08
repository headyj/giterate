package entities

type Service struct {
	BaseURL           string   `json:"BaseURL" yaml:"BaseURL"`
	Name              string   `json:"Name" yaml:"Name"`
	API               string   `json:"API" yaml:"API"`
	APIURI            string   `json:"ApiURI" yaml:"ApiURI"`
	APIToken          string   `json:"ApiToken" yaml:"ApiToken"`
	CloneType         string   `json:"CloneType" yaml:"CloneType"`
	SSHPrivateKeyPath string   `json:"SSHPrivateKeyPath" yaml:"SSHPrivateKeyPath"`
	Username          string   `json:"Username" yaml:"Username"`
	Password          string   `json:"Password" yaml:"Password"`
	Destination       string   `json:"Destination" yaml:"Destination"`
	Entities          []Entity `json:"Entities" yaml:"Entities"`
	Repositories      []Repository
	CloneOptions      []Option `json:"CloneOptions" yaml:"CloneOptions"`
}
