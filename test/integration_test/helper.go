package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/anrisys/quicket/internal/router"
	"github.com/anrisys/quicket/pkg/di"
)

type TestServer struct {
	Server *httptest.Server
	App *di.App
}

func NewTestServer() *TestServer {
	app, err := di.InitializeApp()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	r := router.SetupRouter(app)
	ts := httptest.NewServer(r)

	return &TestServer{Server: ts, App: app}
}

func (ts *TestServer) Close() {
	ts.Server.Close()
}

func (ts *TestServer) URL() string {
	return ts.Server.URL
}

func (ts *TestServer) MakeRequest(method, path, token string, body interface{}) (*http.Response, error) {
    var buf bytes.Buffer
    if body != nil {
        if err := json.NewEncoder(&buf).Encode(body); err != nil {
            return nil, err
        }
    }
	fmt.Printf("PAYLOAD REQUEST IN JSON : %s", &buf)

    req, err := http.NewRequest(method, ts.URL()+path, &buf)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    if token != "" {
        req.Header.Set("Authorization", "Bearer "+token)
    }
    
    return http.DefaultClient.Do(req)
}