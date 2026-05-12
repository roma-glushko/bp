package config

type Config struct {
	Port    int
	DataDir string
	NoOpen  bool
	Debug   bool
}

func DefaultConfig() Config {
	dataDir, _ := DataDir()
	return Config{
		Port:    7391,
		DataDir: dataDir,
		NoOpen:  false,
		Debug:   false,
	}
}
