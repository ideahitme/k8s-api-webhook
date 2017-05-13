package main

import (
	"log"

	"github.com/ideahitme/k8s-api-webhook/authn"
	"github.com/ideahitme/k8s-api-webhook/authn/provider"
)

func main() {
	staticAuthn, err := provider.NewStaticAuthenticator("./tmp")
	if err != nil {
		log.Fatal(err)
	}
	authn.NewAuthenticationHandler(staticAuthn, authn.WithAPIVersion(authn.V1Beta1))
}
