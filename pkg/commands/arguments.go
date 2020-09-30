package command

type Arguments struct {
	ConfigFile string
}

func (a *Arguments) process(args []string) *Arguments {
	for i, v := range args {
		switch v {
		case "--config-file":
			a.ConfigFile = args[i+1]
		}
	}
	return a

}
