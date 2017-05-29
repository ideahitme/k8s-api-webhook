package authz

import (
	"net/http"

	"github.com/ideahitme/k8s-api-webhook/authz/authorizer"
	"github.com/ideahitme/k8s-api-webhook/authz/v1beta1"
)

// AuthorizationHandler implements the webhook handler
type AuthorizationHandler struct {
	resourceAuthorizer    authorizer.ResourceAuthorizer
	nonResourceAuthorizer authorizer.NonResourceAuthorizer
	resConstructor        ResponseConstructor
	reqParser             RequestParser
}

// CreateAuthorizationHandler returns authentication http handler
func CreateAuthorizationHandler() *AuthorizationHandler {
	h := &AuthorizationHandler{
		resourceAuthorizer:    authorizer.ResourceUnauthorizer{},
		nonResourceAuthorizer: authorizer.NonResourceUnauthorizer{},
		resConstructor:        &v1beta1.ResponseConstructor{},
		reqParser:             &v1beta1.RequestParser{},
	}

	return h
}

// WithAPIVersion specify API version to use for handling authentication requests
func (h *AuthorizationHandler) WithAPIVersion(apiVersion APIVersion) *AuthorizationHandler {
	if apiVersion == V1Beta1 {
		h.resConstructor = &v1beta1.ResponseConstructor{}
		h.reqParser = &v1beta1.RequestParser{}
	}
	return h
}

// WithResourceAuthorizer specify API version to use for handling authentication requests
func (h *AuthorizationHandler) WithResourceAuthorizer(authz authorizer.ResourceAuthorizer) *AuthorizationHandler {
	h.resourceAuthorizer = authz
	return h
}

// WithNonResourceAuthorizer specify API version to use for handling authentication requests
func (h *AuthorizationHandler) WithNonResourceAuthorizer(authz authorizer.NonResourceAuthorizer) *AuthorizationHandler {
	h.nonResourceAuthorizer = authz
	return h
}

func (h *AuthorizationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.reqParser.ReadBody(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(h.resConstructor.NewFailResponse(err.Error()))
		return
	}

	userSpec := h.reqParser.ExtractUserSpecs()
	var allowed bool
	var err error

	if h.reqParser.IsResourceRequest() {
		resourceSpec := h.reqParser.ExtractResourceSpecs()
		allowed, err = h.resourceAuthorizer.IsAuthorized(userSpec, resourceSpec)
	}

	if h.reqParser.IsNonResourceRequest() {
		nonResourceSpec := h.reqParser.ExtractNonResourceSpecs()
		allowed, err = h.nonResourceAuthorizer.IsAuthorized(userSpec, nonResourceSpec)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(h.resConstructor.NewFailResponse(err.Error()))
		return
	}
	if !allowed {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(h.resConstructor.NewFailResponse("Unauthorized"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(h.resConstructor.NewSuccessResponse())
}
