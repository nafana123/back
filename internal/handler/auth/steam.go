package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

func (h *AuthHandler) SteamLogin(w http.ResponseWriter, r *http.Request) {
	steamLoginURL := "https://steamcommunity.com/openid/login"
	callbackURL := "http://localhost:8080/api/auth/steam/callback"

	params := url.Values{}
	params.Set("openid.ns", "http://specs.openid.net/auth/2.0")
	params.Set("openid.mode", "checkid_setup")
	params.Set("openid.return_to", callbackURL)
	params.Set("openid.realm", "http://localhost:8080")
	params.Set("openid.identity", "http://specs.openid.net/auth/2.0/identifier_select")
	params.Set("openid.claimed_id", "http://specs.openid.net/auth/2.0/identifier_select")

	redirectURL := fmt.Sprintf("%s?%s", steamLoginURL, params.Encode())
	h.Logger.Info("Редирект на Steam", zap.String("url", redirectURL))
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *AuthHandler) SteamCallback(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("Получен callback от Steam", zap.Any("params", r.URL.Query()))

	queryParams := r.URL.Query()
	h.Logger.Info("queryParams", zap.Any("queryParams", queryParams))

	isValid, steamID, err := h.validateSteamResponse(queryParams)
	if err != nil {
		h.Logger.Error("Ошибка валидации Steam ответа", zap.Error(err))
		http.Error(w, "Validation failed", http.StatusInternalServerError)
		return
	}

	if !isValid {
		h.Logger.Error("Невалидный ответ от Steam")
		http.Error(w, "Invalid Steam response", http.StatusUnauthorized)
		return
	}

	h.Logger.Info("Успешная валидация Steam", zap.String("steamID", steamID))

	userProfile, err := h.getSteamUserProfile(steamID)
	profileJSON, _ := json.Marshal(userProfile)
	encodedProfile := base64.URLEncoding.EncodeToString(profileJSON)

	// TODO юзеров наверное надо записывать в базу при первичной авторизации через стим, а последующие атворизации проверять есть в базе такой или нет

	frontendURL := fmt.Sprintf("http://localhost:5173/auth/steam/callback?profile=%s", encodedProfile)
	http.Redirect(w, r, frontendURL, http.StatusFound)
}

func (h *AuthHandler) validateSteamResponse(params url.Values) (bool, string, error) {
	validationParams := url.Values{}

	for key, values := range params {
		if key != "openid.mode" {
			validationParams[key] = values
		}
	}

	validationParams.Set("openid.mode", "check_authentication")

	resp, err := http.PostForm("https://steamcommunity.com/openid/login", validationParams)
	if err != nil {
		return false, "", fmt.Errorf("ошибка запроса к Steam: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	responseStr := string(body)
	h.Logger.Debug("Ответ от Steam при валидации", zap.String("response", responseStr))

	isValid := strings.Contains(responseStr, "is_valid:true")

	claimedID := params.Get("openid.claimed_id")
	steamID := ""
	parts := strings.Split(claimedID, "/")

	if len(parts) > 0 {
		steamID = parts[len(parts)-1]
	}

	return isValid, steamID, nil
}

func (h *AuthHandler) getSteamUserProfile(steamID string) (map[string]interface{}, error) {
	apiURL := fmt.Sprintf("http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s", "87A337265F5625C3FEC8913E7FAB81E7", steamID)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к Steam API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа Steam API: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON %w", err)
	}

	return result, nil
}
