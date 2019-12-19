package configuration

func paths() []string {
	paths := make([]string, 0)
	return append(paths,
		"/etc/dionysos",
		"XDG_CONFIG_HOME/dionysos",
		"$HOME/.dionysos",
		".")
}
