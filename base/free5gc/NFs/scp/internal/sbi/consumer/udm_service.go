package consumer

import (
	"sync"

	"github.com/free5gc/openapi"
	Nudm_UEAU "github.com/free5gc/openapi/Nudm_UEAuthentication"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/scp/internal/logger"
)

type nudmService struct {
	consumer *Consumer

	ueauMu sync.RWMutex

	ueauClients map[string]*Nudm_UEAU.APIClient
}

func (s *nudmService) getUdmUeauClient(uri string) *Nudm_UEAU.APIClient {
	if uri == "" {
		return nil
	}
	s.ueauMu.RLock()
	client, ok := s.ueauClients[uri]
	if ok {
		s.ueauMu.RUnlock()
		return client
	}

	configuration := Nudm_UEAU.NewConfiguration()
	configuration.SetBasePath(uri)
	client = Nudm_UEAU.NewAPIClient(configuration)

	s.ueauMu.RUnlock()
	s.ueauMu.Lock()
	defer s.ueauMu.Unlock()
	s.ueauClients[uri] = client
	return client
}

func (s *nudmService) SendGenerateAuthDataRequest(uri string,
	supiOrSuci string, authInfoReq *models.AuthenticationInfoRequest) (*models.AuthenticationInfoResult, *models.ProblemDetails, error) {

	client := s.getUdmUeauClient(uri)
	if client == nil {
		return nil, nil, openapi.ReportError("udm not found")
	}

	// TODO: OAuth UDM Generate Auth Data Post
	var authInfoResult models.AuthenticationInfoResult
	ctx, problemDetails, err := s.consumer.scp.Context().GetTokenCtx(models.ServiceName_NUDM_UEAU, models.NfType_UDM)
	if err != nil {
		return nil, problemDetails, err
	}
	authInfoResult, response, err := client.GenerateAuthDataApi.GenerateAuthData(
		ctx, supiOrSuci, *authInfoReq,
	)
	if response != nil && err != nil {
		rspCode, rspBody := handleAPIServiceResponseError(response, err)
		logger.ConsumerLog.Errorf("UE Authentication Response Error: %+v", rspBody)
		return &authInfoResult, &models.ProblemDetails{
			Status: int32(rspCode),
			Cause:  rspBody.(*models.ProblemDetails).Cause,
		}, err
	}
	if err != nil {
		rspCode, rspBody := handleAPIServiceNoResponse(err)
		return &authInfoResult, &models.ProblemDetails{
			Status: int32(rspCode),
			Cause:  rspBody.(*models.ProblemDetails).Cause,
		}, err

	}
	return &authInfoResult, nil, nil
}
