package entities

type Arguments struct {
	ConfigFile   string
	Full         bool
	Force        bool
	LogLevel     string
	Providers    []string
	Repositories []string
}

func (a *Arguments) Process(args []string) *Arguments {
	for i, v := range args {
		switch v {
		case "--config-file":
			a.ConfigFile = args[i+1]
		case "--full", "-f":
			a.Full = true
		case "--force":
			a.Force = true
		case "--log-level":
			a.LogLevel = args[i+1]
		case "--provider", "-p":
			a.Providers = append(a.Providers, args[i+1])
		case "--repository", "-r":
			a.Repositories = append(a.Repositories, args[i+1])
		}
	}
	return a

}
