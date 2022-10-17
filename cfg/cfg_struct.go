package cfg

type App struct {
	Name      string
	Port      int
	DebugPort int
	Env       string
	Secret    string
}

type DbConfig struct {
	DSN         string
	MaxOpenCons int
	MinOpenCons int
}

type Telemetry struct {
	URL      string
	Key      string
	Insecure bool
}
