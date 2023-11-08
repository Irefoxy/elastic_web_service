package app

import (
	"elastic_web_service/internal/auth"
	"elastic_web_service/internal/model"
	"elastic_web_service/internal/repo"
	"elastic_web_service/internal/service"
	"gopkg.in/yaml.v3"
	"os"
)

type provider struct {
	srv  service.Service
	auth auth.Auth
	repo repo.Repo
	cfg  model.Config
}

func NewProvider(path string) (*provider, error) {
	var cfg = model.Config{}
	configFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		return nil, err
	}
	return &provider{cfg: cfg}, nil
}

func (p *provider) Service() service.Service {
	if p.srv == nil {
		p.srv = service.NewHandler(p.Repo(), p.Auth(), p.cfg.Address)
	}
	return p.srv
}

func (p *provider) Repo() repo.Repo {
	if p.repo == nil {
		p.repo = repo.NewRepo()
	}
	return p.repo
}

func (p *provider) Auth() auth.Auth {
	if p.auth == nil {
		p.auth = auth.NewJWT(p.cfg.Secret)
	}
	return p.auth
}
