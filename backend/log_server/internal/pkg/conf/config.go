package conf

const (
	ModeProd       = "prod"
	ModeDev        = "dev"
)

type Config struct {
	Name      string
	Host      string
	Port      int
	Mode      string
}

type LogConfig struct {
	MeasurementSharePrefix string
	MeasurementBlockPrefix string
}
