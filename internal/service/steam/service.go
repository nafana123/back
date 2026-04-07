package steam

import (
	"back/internal/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type SteamService interface {
	GetAuthURL() string
	ValidateAndGetProfile(params url.Values) ([]byte, error)
}

type steamService struct {
	cfg *config.Config
}

func NewSteamService(cfg *config.Config) SteamService {
	return &steamService{
		cfg: cfg,
	}
}

func (s *steamService) GetAuthURL() string {
	steamLoginURL := "https://steamcommunity.com/openid/login"
	callbackURL := s.cfg.SteamCallbackURL

	params := url.Values{}
	params.Set("openid.ns", "http://specs.openid.net/auth/2.0")
	params.Set("openid.mode", "checkid_setup")
	params.Set("openid.return_to", callbackURL)
	params.Set("openid.realm", s.cfg.SteamRealm)
	params.Set("openid.identity", "http://specs.openid.net/auth/2.0/identifier_select")
	params.Set("openid.claimed_id", "http://specs.openid.net/auth/2.0/identifier_select")

	return fmt.Sprintf("%s?%s", steamLoginURL, params.Encode())
}

func (s *steamService) ValidateAndGetProfile(params url.Values) ([]byte, error) {
	validationParams := url.Values{}

	for key, values := range params {
		if key != "openid.mode" {
			validationParams[key] = values
		}
	}

	validationParams.Set("openid.mode", "check_authentication")

	resp, err := http.PostForm("https://steamcommunity.com/openid/login", validationParams)
	if err != nil {
		return nil, fmt.Errorf("Ошибка запроса к Steam: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Ошибка чтения ответа: %w", err)
	}

	responseStr := string(body)

	isValid := strings.Contains(responseStr, "is_valid:true")
	if !isValid {
		return nil, fmt.Errorf("Ошибка валидации данных от стим")
	}

	claimedID := params.Get("openid.claimed_id")
	steamID := ""
	parts := strings.Split(claimedID, "/")

	if len(parts) > 0 {
		steamID = parts[len(parts)-1]
	}

	userData, err := s.getSteamUserProfile(steamID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения данных %w", err)
	}

	jsonData, err := json.Marshal(userData)
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга JSON: %w", err)
	}

	return jsonData, nil
}

func (s *steamService) getSteamUserProfile(steamID string) (map[string]interface{}, error) {
	apiURL := fmt.Sprintf("http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s", s.cfg.SteamAPIKey, steamID)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}