package integration_test

import (
	"log"
	"net/http/httptest"

	"github.com/anrisys/quicket/internal/router"
	"github.com/anrisys/quicket/pkg/di"
)

type TestServer struct {
	Server *httptest.Server
}

func NewTestServer() *TestServer {
	app, err := di.InitializeApp()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	r := router.SetupRouter(app)
	ts := httptest.NewServer(r)

	return &TestServer{Server: ts}
}

func (ts *TestServer) Close() {
	ts.Server.Close()
}

func (ts *TestServer) URL() string {
	return ts.Server.URL
}