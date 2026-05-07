package steam

import (
	"back/internal/config"
	steamdto "back/internal/dto/steam"
	"back/internal/model"
	steamrepo "back/internal/repository/steam"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type SteamService struct {
	cfg             *config.Config
	steamRepository *steamrepo.SteamRepository
}

func NewSteamService(cfg *config.Config, steamRepository *steamrepo.SteamRepository) *SteamService {
	return &SteamService{
		cfg:             cfg,
		steamRepository: steamRepository,
	}
}

func (s *SteamService) GetAuthURL(state string) string {
	steamLoginURL := "https://steamcommunity.com/openid/login"
	callbackURL, err := url.Parse(s.cfg.SteamCallbackURL)
	if err != nil {
		callbackURL = &url.URL{Path: s.cfg.SteamCallbackURL}
	}
	callbackParams := callbackURL.Query()
	callbackParams.Set("state", state)
	callbackURL.RawQuery = callbackParams.Encode()

	params := url.Values{}
	params.Set("openid.ns", "http://specs.openid.net/auth/2.0")
	params.Set("openid.mode", "checkid_setup")
	params.Set("openid.return_to", callbackURL.String())
	params.Set("openid.realm", s.cfg.SteamRealm)
	params.Set("openid.identity", "http://specs.openid.net/auth/2.0/identifier_select")
	params.Set("openid.claimed_id", "http://specs.openid.net/auth/2.0/identifier_select")

	return fmt.Sprintf("%s?%s", steamLoginURL, params.Encode())
}

func (s *SteamService) ValidateCallback(params url.Values, userID int) error {
	validationParams := url.Values{}

	for key, values := range params {
		if key != "openid.mode" {
			validationParams[key] = values
		}
	}

	validationParams.Set("openid.mode", "check_authentication")

	resp, err := http.PostForm("https://steamcommunity.com/openid/login", validationParams)
	if err != nil {
		return fmt.Errorf("Ошибка запроса к Steam: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Steam OpenID вернул статус %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Ошибка чтения ответа: %w", err)
	}

	responseStr := string(body)

	isValid := strings.Contains(responseStr, "is_valid:true")
	if !isValid {
		return fmt.Errorf("Ошибка валидации данных от стим")
	}

	claimedID := params.Get("openid.claimed_id")
	steamID := ""
	parts := strings.Split(claimedID, "/")

	if len(parts) > 0 {
		steamID = parts[len(parts)-1]
	}

	userData, err := s.getSteamData(steamID)
	if err != nil {
		return fmt.Errorf("ошибка получения данных %w", err)
	}

	steamUser := &model.SteamUser{
		UserID:      userID,
		SteamID:     userData.Response.Players[0].SteamID,
		PersonaName: userData.Response.Players[0].PersonaName,
		AvatarURL:   userData.Response.Players[0].AvatarFull,
	}

	err = s.steamRepository.CreateSteamUser(steamUser)
	if err != nil {
		return fmt.Errorf("ошибка создания steam пользователя: %w", err)
	}

	return nil
}

func (s *SteamService) getSteamData(steamID string) (steamdto.SteamDataResponse, error) {
	apiURL := fmt.Sprintf("http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s", s.cfg.SteamAPIKey, steamID)

	resp, err := http.Get(apiURL)
	if err != nil {
		return steamdto.SteamDataResponse{}, fmt.Errorf("ошибка запроса к Steam: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return steamdto.SteamDataResponse{}, fmt.Errorf("ошибка к запросу Steam Web API%d", resp.StatusCode)
	}

	var steamData steamdto.SteamDataResponse
	err = json.NewDecoder(resp.Body).Decode(&steamData)
	if err != nil {
		return steamdto.SteamDataResponse{}, fmt.Errorf("ошибка декодирования JSON: %w", err)
	}

	if len(steamData.Response.Players) == 0 {
		return steamdto.SteamDataResponse{}, fmt.Errorf("steam пользователь не найден")
	}

	return steamData, nil
}


func (s *SteamService) Logout(userID int) error {
	err := s.steamRepository.DeleteSteamUser(userID)
	if err != nil {
		return fmt.Errorf("ошибка удаления steam пользователя: %w", err)
	}

	return nil
}