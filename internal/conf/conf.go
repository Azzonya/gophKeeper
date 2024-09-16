package conf

import (
	"flag"
	"github.com/caarlos0/env/v9"
)

// Conf represents the application configuration.
var Conf = struct {
	ServerAddress string `env:"server_address"`

	RedisAddress   string `env:"redis_address" envDefault:"localhost:6379"`
	RedisPassword  string `env:"redis_password" envDefault:"password"`
	RedisDB        int    `env:"redis_db" envDefault:"0"`
	ClientCertFile string `env:"client_cert_file" envDefault:"cert/client-cert.pem"`
	ClientKeyFile  string `env:"client_key_file" envDefault:"cert/client-key.pem"`

	ServerCertFile string `env:"SERVER_CERT_FILE" envDefault:"cert/server-cert.pem"`
	ServerKeyFile  string `env:"SERVER_KEY_FILE" envDefault:"cert/server-key.pem"`
	CAFile         string `env:"CA_FILE" envDefault:"cert/ca-cert.pem"`

	GRPCPort string `env:"GRPC_PORT"`

	PgDsn string `env:"DATABASE_URI"`

	JwtSecret string `env:"JWT_SECRET"`

	S3Endpoint  string `env:"S3_ENDPOINT" envDefault:"localhost:9000"`
	S3Bucket    string `env:"S3_BUCKET" envDefault:"mybucket"`
	S3AccessKey string `env:"S3_ACCESS_KEY" envDefault:"minioadmin"`
	S3SecretKey string `env:"S3_SECRET_KEY" envDefault:"minioadmin"`

	EnableTLS bool `env:"ENABLE_TLS" envDefault:"true"`
}{}

// init initializes the configuration by parsing command-line flags and environment variables.
// It sets default values for fields like the gRPC server address and certificate paths.
func init() {
	flag.StringVar(&Conf.ServerAddress, "a", "localhost:5050", "address and port where grpc server start")
	flag.StringVar(&Conf.GRPCPort, "g", ":5050", "address and port where grpc server start")
	flag.StringVar(&Conf.PgDsn, "d", "", "database connection line")

	if err := env.Parse(&Conf); err != nil {
		panic(err)
	}
}
