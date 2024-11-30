package consumer

import (
	"context"
	"net/http"
	"sync"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/Nausf_UEAuthentication"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/scp/internal/logger"
)

type nausfService struct {
	consumer *Consumer

	UEAuthenticationMu sync.RWMutex

	UEAuthenticationClients map[string]*Nausf_UEAuthentication.APIClient
}

func (s *nausfService) getUEAuthenticationClient(uri string) *Nausf_UEAuthentication.APIClient {
	if uri == "" {
		return nil
	}
	s.UEAuthenticationMu.RLock()
	client, ok := s.UEAuthenticationClients[uri]
	if ok {
		s.UEAuthenticationMu.RUnlock()
		return client
	}

	configuration := Nausf_UEAuthentication.NewConfiguration()
	configuration.SetBasePath(uri)
	client = Nausf_UEAuthentication.NewAPIClient(configuration)

	s.UEAuthenticationMu.RUnlock()
	s.UEAuthenticationMu.Lock()
	defer s.UEAuthenticationMu.Unlock()
	s.UEAuthenticationClients[uri] = client
	return client
}

func (s *nausfService) SendUeAuthPostRequest(uri string,
	authInfo *models.AuthenticationInfo) (*models.UeAuthenticationCtx, *models.ProblemDetails, error) {
	logger.ConsumerLog.Debugf("[AMF->AUSF] Forward AMF UE Authentication Request")
	client := s.getUEAuthenticationClient(uri)
	if client == nil {
		return nil, nil, openapi.ReportError("ausf not found")
	}

	// TODO: OAuth AUSF Ue Auth Post
	var ueAuthenticationCtx models.UeAuthenticationCtx
	response := &http.Response{}
	err := error(nil)
	Info := models.AuthenticationInfo(*authInfo)
	ueAuthenticationCtx, response, err = client.DefaultApi.UeAuthenticationsPost(
		context.Background(), Info)
	if response != nil && err != nil {
		rspCode, rspBody := handleAPIServiceResponseError(response, err)
		logger.ConsumerLog.Errorf("UE Authentication Response Error: %+v", rspBody)
		return &ueAuthenticationCtx, &models.ProblemDetails{
			Status: int32(rspCode),
			Cause:  rspBody.(*models.ProblemDetails).Cause,
		}, err
	}
	if err != nil {
		rspCode, rspBody := handleAPIServiceNoResponse(err)
		return &ueAuthenticationCtx, &models.ProblemDetails{
			Status: int32(rspCode),
			Cause:  rspBody.(*models.ProblemDetails).Cause,
		}, err

	}
	logger.ConsumerLog.Debugf("UE Authentication Response: %+v", ueAuthenticationCtx)
	return &ueAuthenticationCtx, nil, nil
}

func (s *nausfService) SendAuth5gAkaConfirmRequest(uri string,
	authCtxId string, confirmationData *models.ConfirmationData) (*models.ConfirmationDataResponse, *models.ProblemDetails, error) {

	client := s.getUEAuthenticationClient(uri)
	if client == nil {
		return nil, nil, openapi.ReportError("ausf not found")
	}

	// TODO: OAuth AUSF Auth 5gAka Confirm Put
	var confirmResult models.ConfirmationDataResponse
	return &confirmResult, nil, nil
}
