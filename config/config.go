package config

type Admin struct {
	Name     string
	Email    string
	Password string
}

type LinkUrl struct {
	Alias       string `json:"alias"`
	RedirectURL string `json:"redirectURL"`
}

type Config struct {
	Admin Admin
	Links []LinkUrl
}

func GetConfig() Config {
	config := Config{
		Admin: Admin{
			Name:     "Admin User",
			Password: "password",
			Email:    "admin@test.com",
		},
		Links: []LinkUrl{
			{
				Alias:       "hydyco",
				RedirectURL: "https://hydyco.com",
			},
			{
				Alias:       "heeico",
				RedirectURL: "https://heeico.com",
			},
			{
				Alias:       "google",
				RedirectURL: "https://google.com",
			},
		},
	}
	return config
}
