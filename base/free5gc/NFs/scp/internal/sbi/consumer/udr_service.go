package consumer

import (
	"context"
	"sync"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/Nudr_DataRepository"
	"github.com/free5gc/openapi/models"
)

type nudrService struct {
	consumer *Consumer

	mu      sync.RWMutex
	clients map[string]*Nudr_DataRepository.APIClient
}

func (s *nudrService) getClient(uri string) *Nudr_DataRepository.APIClient {
	s.mu.RLock()
	if client, ok := s.clients[uri]; ok {
		defer s.mu.RUnlock()
		return client
	} else {
		configuration := Nudr_DataRepository.NewConfiguration()
		configuration.SetBasePath(uri)
		cli := Nudr_DataRepository.NewAPIClient(configuration)

		s.mu.RUnlock()
		s.mu.Lock()
		defer s.mu.Unlock()
		s.clients[uri] = cli
		return cli
	}
}

func (s *nudrService) SendAuthSubsDataGet(uri string,
	supi string) (*models.AuthenticationSubscription, *models.ProblemDetails, error) {

	client := s.getClient(uri)
	if client == nil {
		return nil, nil, openapi.ReportError("udr not found")
	}

	// TODO: OAuth UDR Auth Subs Data Get
	var authSubs models.AuthenticationSubscription
	ctx := context.Background()
	authSubs, _, err := client.AuthenticationDataDocumentApi.QueryAuthSubsData(
		ctx, supi, nil,
	)
	if err != nil {
		return nil, nil, err
	}
	return &authSubs, nil, nil
}
