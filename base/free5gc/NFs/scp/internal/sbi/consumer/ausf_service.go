package consumer

import (
	"context"
	"net/http"
	"sync"

	"github.com/antihax/optional"
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
	logger.ConsumerLog.Infof("SendUeAuthPostRequest: %s", uri)
	client := s.getUEAuthenticationClient(uri)
	if client == nil {
		logger.ConsumerLog.Errorf("ausf not found")
		return nil, nil, openapi.ReportError("ausf not found")
	}

	// TODO: OAuth AUSF Ue Auth Post
	var ueAuthenticationCtx models.UeAuthenticationCtx
	ctx := context.Background()
	err := error(nil)
	ueAuthenticationCtx, _, err = client.DefaultApi.UeAuthenticationsPost(
		ctx, *authInfo,
	)
	if err != nil {
		return nil, nil, err
	}
	return &ueAuthenticationCtx, nil, nil
}

func (s *nausfService) SendAuth5gAkaConfirmRequest(uri string,
	authCtxId string, confirmationData *models.ConfirmationData) (*models.ConfirmationDataResponse, *models.ProblemDetails, error) {
	logger.ConsumerLog.Infof("SendAuth5gAkaConfirmRequest: %s", uri)
	client := s.getUEAuthenticationClient(uri)
	if client == nil {
		return nil, nil, openapi.ReportError("ausf not found")
	}

	// TODO: OAuth AUSF Auth 5gAka Confirm Put
	var confirmResult models.ConfirmationDataResponse
	ctx := context.Background()
	err := error(nil)
	var res *http.Response
	confirmResult, res, err = client.DefaultApi.UeAuthenticationsAuthCtxId5gAkaConfirmationPut(
		ctx, authCtxId, &Nausf_UEAuthentication.UeAuthenticationsAuthCtxId5gAkaConfirmationPutParamOpts{
			ConfirmationData: optional.NewInterface(*confirmationData),
		},
	)
	if err != nil {
		handleAPIServiceResponseError(res, err)
		return nil, nil, err
	}
	return &confirmResult, nil, nil
}
