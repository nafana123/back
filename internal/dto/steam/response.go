package steam

type SteamDataResponse struct {
	Response struct {
		Players []struct {
			SteamID     string `json:"steamid"`
			PersonaName string `json:"personaname"`
			AvatarFull  string `json:"avatarfull"`
		} `json:"players"`
	} `json:"response"`
}

type SteamRedirectResponse struct {
	RedirectURL string `json:"redirect_url"`
}