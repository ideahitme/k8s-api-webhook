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

// Option extends default AuthenticationHandler
type Option func(*AuthenticationHandler)

// NewAuthenticationHandler returns authentication http handler
func NewAuthenticationHandler(p authenticator.Authenticator, opts ...Option) *AuthenticationHandler {
	h := &AuthenticationHandler{
		authProvider:   p,
		resConstructor: v1beta1.ResponseConstructor{},
		reqParser:      v1beta1.RequestParser{},
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

// WithAPIVersion specify API version to use for handling authentication requests
func WithAPIVersion(apiVersion APIVersion) func(*AuthenticationHandler) {
	return func(h *AuthenticationHandler) {
		if apiVersion == V1Beta1 {
			h.resConstructor = v1beta1.ResponseConstructor{}
			h.reqParser = v1beta1.RequestParser{}
		}
	}
}

func (h *AuthenticationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, err := h.reqParser.ExtractToken(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(h.resConstructor.NewFailResponse())
		return
	}
	defer r.Body.Close()

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
