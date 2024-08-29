// Package conf manages the configuration settings for the GophKeeper client,
// allowing configuration through both command-line flags and environment variables.
package conf

import (
	"flag"
	"github.com/caarlos0/env/v9"
)

// Conf holds the configuration settings for the client, including the gRPC server address,
// paths to the CA and client certificates, and the option to enable TLS.
var Conf = struct {
	ServerAddress  string `env:"server_address"`
	RedisAddress   string `env:"redis_address" envDefault:"localhost:6379"`
	RedisPassword  string `env:"redis_password" envDefault:"password"`
	RedisDb        int    `env:"redis_db" envDefault:"0"`
	CAFile         string `env:"CA_FILE" envDefault:"cert/ca-cert.pem"`
	ClientCertFile string `env:"client_cert_file" envDefault:"cert/client-cert.pem"`
	ClientKeyFile  string `env:"client_key_file" envDefault:"cert/client-key.pem"`
	EnableTLS      bool   `env:"ENABLE_TLS" envDefault:"true"`
}{}

// init initializes the configuration by parsing command-line flags and environment variables.
// It sets default values for fields like the gRPC server address and certificate paths.
func init() {
	flag.StringVar(&Conf.ServerAddress, "a", "localhost:5050", "address and port where grpc server start")

	if err := env.Parse(&Conf); err != nil {
		panic(err)
	}
}
