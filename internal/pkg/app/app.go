package app

import (
	"fmt"
	"github.com/pheezz/medods-jwt-test/internal/app/endpoint"
	"github.com/pheezz/medods-jwt-test/internal/app/service"
	"net/http"
)

type App struct {
	endpoint *endpoint.Endpoint
	service  *service.Service
}

func New() (*App, error) {
	app := &App{}
	app.service = service.New()
	app.endpoint = endpoint.New(app.service)
	http.HandleFunc("/keyPair", app.endpoint.GetKeyPair)
	http.HandleFunc("/refreshPair", app.endpoint.RefreshKeyPair)
	return app, nil
}

func (a *App) Run() error {
	const port = ":8080"
	fmt.Println("Server is running on port ", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		return err
	}
	return nil
}
