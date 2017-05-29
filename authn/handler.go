package authn

import (
	"net/http"

	"github.com/ideahitme/k8s-api-webhook/authn/authenticator"
	"github.com/ideahitme/k8s-api-webhook/authn/versioned"
	"github.com/ideahitme/k8s-api-webhook/authn/versioned/v1beta1"
)

// AuthenticationHandler implements the webhook handler
type AuthenticationHandler struct {
	authenticator  authenticator.Authenticator
	resConstructor versioned.ResponseConstructor
	reqParser      versioned.RequestParser
}

// CreateAuthenticationHandler returns default authentication http handler
func CreateAuthenticationHandler() *AuthenticationHandler {
	h := &AuthenticationHandler{
		authenticator:  authenticator.Noop{},
		resConstructor: v1beta1.ResponseConstructor{},
		reqParser:      v1beta1.RequestParser{},
	}

	return h
}

// WithAuthenticator adds authenticator to overwrite default noop authenticator
func (h *AuthenticationHandler) WithAuthenticator(p authenticator.Authenticator) *AuthenticationHandler {
	h.authenticator = p
	return h
}

// WithAPIVersion specify API version to use for handling authentication requests
func (h *AuthenticationHandler) WithAPIVersion(apiVersion versioned.APIVersion) *AuthenticationHandler {
	if apiVersion == versioned.V1Beta1 {
		h.resConstructor = v1beta1.ResponseConstructor{}
		h.reqParser = v1beta1.RequestParser{}
	}
	return h
}

func (h *AuthenticationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, err := h.reqParser.ExtractToken(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(h.resConstructor.NewFailResponse())
		return
	}

	user, err := h.authenticator.Authenticate(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(h.resConstructor.NewFailResponse())
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(h.resConstructor.NewFailResponse())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(h.resConstructor.NewSuccessResponse(user))
}
