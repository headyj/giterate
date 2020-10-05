package command

type Arguments struct {
	ConfigFile string
	Full       bool
	Force      bool
	LogLevel   string
}

func (a *Arguments) process(args []string) *Arguments {
	for i, v := range args {
		switch v {
		case "--config-file":
			a.ConfigFile = args[i+1]
		case "--full":
			a.Full = true
		case "--force":
			a.Force = true
		case "--log-level":
			a.LogLevel = args[i+1]
		}
	}
	return a

}
