package stdcrypto

type Config struct {
	Salt string `env:"SECURITY_SALT"`
	Algo string `env:"SECURITY_ALGO,default=md5"`
}
