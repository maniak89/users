package hydra

type Config struct {
	HydraAdminUrl string `env:"HYDRA_ADMIN_URL,required"`
}
