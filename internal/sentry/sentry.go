package sentry

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"github.com/Drafteame/draft/internal/config"
	"github.com/Drafteame/draft/internal/http"
)

const baseURL = "https://sentry.io"

func CreateProject(name string) (string, error) {
	cfg := config.Get().Sentry

	token := cfg.Token
	org := cfg.Organization
	team := cfg.Team

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}

	url := baseURL + "/api/0/teams/" + org + "/" + team + "/projects/"

	reqBody := map[string]any{
		"name":     name,
		"platform": "go",
	}

	jb, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	res, errPost := http.Post(context.Background(), url, jb, headers)
	if errPost != nil {
		return "", errPost
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 201 {
		println("[sentry] error with status code", res.StatusCode)
		println("[sentry] error body", string(resBody))
		println("[sentry] request url", url)
		println("[sentry] request body", string(jb))
		println("[sentry] request auth token", token)
		return "", errors.New("sentry: failed to create project")
	}

	body := map[string]any{}

	if errUnm := json.Unmarshal(resBody, &body); errUnm != nil {
		return "", errUnm
	}

	return body["id"].(string), nil
}

func GetClientKeys(projectID string) (map[string]string, error) {
	cfg := config.Get().Sentry

	token := cfg.Token
	org := cfg.Organization

	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	url := baseURL + "/api/0/projects/" + org + "/" + projectID + "/keys/"

	res, errGet := http.Get(context.Background(), url, headers)
	if errGet != nil {
		return nil, errGet
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		println("[sentry] error with status code", res.StatusCode)
		println("[sentry] error body", string(resBody))
		println("[sentry] request url", url)
		println("[sentry] request auth token", token)
		return nil, errors.New("sentry: failed to get client keys")
	}

	var body []map[string]any

	if errUnm := json.Unmarshal(resBody, &body); errUnm != nil {
		return nil, errUnm
	}

	keys := map[string]string{
		"public": body[0]["public"].(string),
		"secret": body[0]["secret"].(string),
		"dsn":    body[0]["dsn"].(map[string]any)["public"].(string),
	}

	return keys, nil
}
