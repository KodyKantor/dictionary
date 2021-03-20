package main

// DictionaryClient provides an HTTP client for accessing dictionary definitions.
type DictionaryClient struct {
	remotePort string
}

func (dc *DictionaryClient) Help() string {
	helpText := `
	Request to set or get a dictionary definition from a remote server.
	`

	return helpText
}

func (dc *DictionaryClient) Synopsis() string {
	return "Get or set definitions."
}

func (dc *DictionaryClient) Run(args []string) int {
	return 0
}
