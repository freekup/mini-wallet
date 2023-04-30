package infra

type (

	// @envconfig (prefix:"APP")
	App struct {
		Name string `envconfig:"NAME" default:"Mini Wallet" required:"true"`
		Key  string `envconfig:"KEY" default:"mini-wallet" required:"true"`
		Env  string `envconfig:"ENV" default:"local" required:"true"`

		Address string `envconfig:"ADDRESS" default:":8089" required:"true"`
	}
)
