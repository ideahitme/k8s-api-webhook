package authz

import (
	"net/http"

	"github.com/ideahitme/k8s-api-webhook/authz/authorizer"
	"github.com/ideahitme/k8s-api-webhook/authz/v1beta1"
)

// AuthorizationHandler implements the webhook handler
type AuthorizationHandler struct {
	authorizer     authorizer.Authorizer
	resConstructor ResponseConstructor
	reqParser      RequestParser
}

// Option extends default AuthorizationHandler
type Option func(*AuthorizationHandler)

// NewAuthorizationHandler returns authentication http handler
func NewAuthorizationHandler(authz authorizer.Authorizer, opts ...Option) *AuthorizationHandler {
	h := &AuthorizationHandler{
		authorizer:     authz,
		resConstructor: v1beta1.ResponseConstructor{},
		reqParser:      v1beta1.RequestParser{},
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

// WithAPIVersion specify API version to use for handling authentication requests
func WithAPIVersion(apiVersion APIVersion) func(*AuthorizationHandler) {
	return func(h *AuthorizationHandler) {
		if apiVersion == V1Beta1 {
			h.resConstructor = v1beta1.ResponseConstructor{}
			h.reqParser = v1beta1.RequestParser{}
		}
	}
}

func (h *AuthorizationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.reqParser.ReadBody(r.Body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(h.resConstructor.NewFailResponse(err.Error()))
		return
	}
	defer r.Body.Close()

	userSpec := h.reqParser.ExtractUserSpecs()

	if h.reqParser.IsResourceRequest() {
		resourceSpec := h.reqParser.ExtractResourceSpecs()
		allowed, err := h.authorizer.ResourceEnforce(userSpec, resourceSpec)
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
	}

	if h.reqParser.IsNonResourceRequest() {
		nonResourceSpec := h.reqParser.ExtractNonResourceSpecs()
		allowed, err := h.authorizer.NonResourceEnforce(userSpec, nonResourceSpec)
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
	}

	w.WriteHeader(http.StatusOK)
	w.Write(h.resConstructor.NewSuccessResponse())
}
