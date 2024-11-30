package processor

import (
	"net/http"

	"github.com/free5gc/openapi/models"
	"github.com/free5gc/scp/internal/logger"
)

// NOTE: Assume Rand from UDM and ausfInstanceId from AUSF is correct
func (p *Processor) PostGenerateAuthData(
	supiOrSuci string,
	authInfo models.AuthenticationInfoRequest,
) *HandlerResponse {
	logger.ProxyLog.Debugln("[AUSF->UDM] Forward AUSF UE Authentication Request")

	// TODO: Send request to target NF by setting correct uri
	var targetNfUri = "http://udm.free5gc.org:8000"

	// TODO: Verify that the Information Elements (IEs) in the request or response body are correct
	//       Recover and handle errors if the IEs are incorrect
	response, problemDetails, err := p.Consumer().SendGenerateAuthDataRequest(targetNfUri, supiOrSuci, &authInfo)

	if response != nil {
		return &HandlerResponse{http.StatusOK, nil, response}
	} else if problemDetails != nil {
		return &HandlerResponse{int(problemDetails.Status), nil, problemDetails}
	}
	logger.ProxyLog.Errorln(err)
	problemDetails = &models.ProblemDetails{
		Status: http.StatusForbidden,
		Cause:  "UNSPECIFIED",
	}

	return &HandlerResponse{http.StatusForbidden, nil, problemDetails}
}
