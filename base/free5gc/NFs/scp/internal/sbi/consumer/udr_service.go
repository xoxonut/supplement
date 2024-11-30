package consumer

import (
	"sync"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/Nudr_DataRepository"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/scp/internal/logger"
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
	ctx, problemDetails, err := s.consumer.scp.Context().GetTokenCtx(models.ServiceName_NUDR_DR, models.NfType_UDR)
	if err != nil {
		return nil, problemDetails, err
	}
	authSubs, response, err := client.AuthenticationDataDocumentApi.QueryAuthSubsData(
		ctx, supi, nil,
	)
	if response != nil && err != nil {
		rspCode, rspBody := handleAPIServiceResponseError(response, err)
		logger.ConsumerLog.Errorf("UE Authentication Response Error: %+v", rspBody)
		return &authSubs, &models.ProblemDetails{
			Status: int32(rspCode),
			Cause:  rspBody.(*models.ProblemDetails).Cause,
		}, err
	}
	if err != nil {
		rspCode, rspBody := handleAPIServiceNoResponse(err)
		return &authSubs, &models.ProblemDetails{
			Status: int32(rspCode),
			Cause:  rspBody.(*models.ProblemDetails).Cause,
		}, err

	}

	return &authSubs, nil, nil
}
