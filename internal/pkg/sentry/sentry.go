package sentry

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/Drafteame/draft/internal/config"
	http2 "github.com/Drafteame/draft/internal/pkg/http"
)

type Project struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

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

	res, errPost := http2.Post(context.Background(), url, jb, headers)
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

	res, errGet := http2.Get(context.Background(), url, headers)
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

func CreateStages(serviceName, dsn string) error {
	transport := sentry.NewHTTPSyncTransport()
	transport.Timeout = time.Second * 1

	c, err := sentry.NewClient(sentry.ClientOptions{
		ServerName:  serviceName,
		Dsn:         dsn,
		Environment: "dev",
		Transport:   transport,
	})

	if err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		c.CaptureException(errors.New("test error"), nil, nil)
	}

	c, err = sentry.NewClient(sentry.ClientOptions{
		ServerName:  serviceName,
		Dsn:         dsn,
		Environment: "prod",
		Transport:   transport,
	})

	if err != nil {
		return err
	}

	for i := 0; i < 5; i++ {
		c.CaptureException(errors.New("test error"), nil, nil)
	}

	return nil
}

func DeleteProject(projectID string) error {
	cfg := config.Get().Sentry

	token := cfg.Token
	org := cfg.Organization

	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	url := baseURL + "/api/0/projects/" + org + "/" + projectID + "/"

	res, errDelete := http2.Delete(context.Background(), url, headers)
	if errDelete != nil {
		return errDelete
	}

	if res.StatusCode != 204 {
		return errors.New("sentry: failed to delete project")
	}

	return nil
}

func ListProjects() (map[string]string, error) {
	cfg := config.Get().Sentry
	token := cfg.Token
	org := cfg.Organization

	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	url := baseURL + "/api/0/projects/" + org + "/"

	res, err := http2.Get(context.Background(), url, headers)
	if err != nil {
		return nil, err
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
		return nil, errors.New("sentry: failed to list projects")
	}

	var projects []Project
	if err := json.Unmarshal(resBody, &projects); err != nil {
		return nil, err
	}

	names := make(map[string]string, len(projects))
	for _, p := range projects {
		names[p.Name] = p.ID
	}

	return names, nil
}
