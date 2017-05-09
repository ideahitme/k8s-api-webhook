package authn

import (
	"net/http"

	"io/ioutil"

	"github.com/ideahitme/k8s-api-webhook/authn/provider"
	"github.com/ideahitme/k8s-api-webhook/authn/v1beta1"
)

// AuthenticationHandler implements the webhook handler
type AuthenticationHandler struct {
	authProvider provider.Authenticator
}

// NewAuthenticationHandler returns authentication http handler
func NewAuthenticationHandler(p provider.Authenticator) *AuthenticationHandler {
	return &AuthenticationHandler{
		authProvider: p,
	}
}

func (h *AuthenticationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// extract and check the payload
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(v1beta1.NewFailResponse())
		return
	}
	defer r.Body.Close()

	token, err := v1beta1.ParseToken(payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(v1beta1.NewFailResponse())
		return
	}

	user, err := h.authProvider.Authenticate(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(v1beta1.NewFailResponse())
		return
	}

	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(v1beta1.NewFailResponse())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(v1beta1.NewSuccessResponse(user))
}
