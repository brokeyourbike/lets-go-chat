package configurations

type Config struct {
	Database struct {
		Url string `env:"DATABASE_URL" envDefault:"postgres://postgres:secret@localhost:5432/test"`
	}
	Host string `env:"HOST" envDefault:""`
	Port string `env:"PORT" envDefault:"8080"`
}
