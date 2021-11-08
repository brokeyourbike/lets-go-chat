package configurations

type Config struct {
	Server struct {
		Host string `env:"HOST" envDefault:""`
		Port string `env:"PORT" envDefault:"8080"`
	}
}
