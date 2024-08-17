// Package conf provides functionality to initialize and parse configuration for the URL shortener application.
package conf

import (
	"flag"
	"github.com/caarlos0/env/v9"
)

// Conf represents the application configuration.
var Conf = struct {
	ServerCertFile string `env:"SERVER_CERT_FILE" envDefault:"cert/server-cert.pem"`
	ServerKeyFile  string `env:"SERVER_KEY_FILE" envDefault:"cert/server-key.pem"`
	CAFile         string `env:"CA_FILE" envDefault:"cert/ca-cert.pem"`
	GRPCPort       string `env:"GRPC_PORT"`
	PgDsn          string `env:"DATABASE_URI"`
	JwtSecret      string `env:"JWT_SECRET"`
	S3Endpoint     string `env:"S3_ENDPOINT" envDefault:"localhost:9000"`
	S3Bucket       string `env:"S3_BUCKET" envDefault:"mybucket"`
	S3AccessKey    string `env:"S3_ACCESS_KEY" envDefault:"minioadmin"`
	S3SecretKey    string `env:"S3_SECRET_KEY" envDefault:"minioadmin"`
	EnableTLS      bool   `env:"ENABLE_TLS" envDefault:"true"`
}{}

// init initializes the configuration for the application by setting up command-line flags
// and parsing environment variables. The flags include options for specifying the gRPC server address
// and port, as well as the database connection string.
//
// The function uses the `env.Parse` method to load configuration values from environment variables.
// If any error occurs during parsing, the application panics.
func init() {
	flag.StringVar(&Conf.GRPCPort, "a", ":5050", "address and port where grpc server start")
	flag.StringVar(&Conf.PgDsn, "d", "", "database connection line")

	if err := env.Parse(&Conf); err != nil {
		panic(err)
	}
}
