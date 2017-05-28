package authn

import (
	"net/http"

	"github.com/ideahitme/k8s-api-webhook/authn/authenticator"
	"github.com/ideahitme/k8s-api-webhook/authn/v1beta1"
)

// AuthenticationHandler implements the webhook handler
type AuthenticationHandler struct {
	authProvider   authenticator.Authenticator
	resConstructor ResponseConstructor
	reqParser      RequestParser
}

// NewAuthenticationHandler returns authentication http handler
func NewAuthenticationHandler(p authenticator.Authenticator) *AuthenticationHandler {
	h := &AuthenticationHandler{
		authProvider:   p,
		resConstructor: v1beta1.ResponseConstructor{},
		reqParser:      v1beta1.RequestParser{},
	}

	return h
}

// WithAPIVersion specify API version to use for handling authentication requests
func (h *AuthenticationHandler) WithAPIVersion(apiVersion APIVersion) *AuthenticationHandler {
	if apiVersion == V1Beta1 {
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

	user, err := h.authProvider.Authenticate(token)
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
