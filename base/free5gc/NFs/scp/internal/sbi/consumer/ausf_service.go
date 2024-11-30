package consumer

import (
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
	logger.ConsumerLog.Debugf("[AMF->AUSF] Forward AMF UE Authentication Request")
	client := s.getUEAuthenticationClient(uri)
	if client == nil {
		return nil, nil, openapi.ReportError("ausf not found")
	}

	// TODO: OAuth AUSF Ue Auth Post
	var ueAuthenticationCtx models.UeAuthenticationCtx
	response := &http.Response{}
	err := error(nil)
	ctx, problemDetails, err := s.consumer.scp.Context().GetTokenCtx(models.ServiceName_NAUSF_AUTH, models.NfType_AUSF)
	if err != nil {
		return &ueAuthenticationCtx, problemDetails, err
	}
	ueAuthenticationCtx, response, err = client.DefaultApi.UeAuthenticationsPost(
		ctx, *authInfo)
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
	ctx, problemDetails, err := s.consumer.scp.Context().GetTokenCtx(
		models.ServiceName_NAUSF_AUTH, models.NfType_AUSF)
	if err != nil {
		return &confirmResult, problemDetails, err

	}
	logger.ConsumerLog.Debugf("[AMF->AUSF] Forward AMF UE Authentication 5gAka Confirm Request")
	data := optional.NewInterface(*confirmationData)
	logger.ConsumerLog.Debugf("ConfirmationData: %+v", confirmationData)
	logger.ConsumerLog.Debugf("Data: %+v", data)
	confirmResult, response, err := client.DefaultApi.UeAuthenticationsAuthCtxId5gAkaConfirmationPut(
		ctx, authCtxId, &Nausf_UEAuthentication.UeAuthenticationsAuthCtxId5gAkaConfirmationPutParamOpts{
			ConfirmationData: data,
		})
	if response != nil && err != nil {
		rspCode, rspBody := handleAPIServiceResponseError(response, err)
		logger.ConsumerLog.Errorf("Auth 5gAka Confirm Response Error: %+v", rspBody)
		return &confirmResult, &models.ProblemDetails{
			Status: int32(rspCode),
			Cause:  rspBody.(*models.ProblemDetails).Cause,
		}, err
	}
	if err != nil {
		logger.ConsumerLog.Errorf("Auth 5gAka Confirm Response Error: %+v", err)
		rspCode, rspBody := handleAPIServiceNoResponse(err)
		return &confirmResult, &models.ProblemDetails{
			Status: int32(rspCode),
			Cause:  rspBody.(*models.ProblemDetails).Cause,
		}, err
	}
	return &confirmResult, nil, nil
}
