package main

import (
	"log"

	"net/http"

	"github.com/ideahitme/k8s-api-webhook/authn"
	"github.com/ideahitme/k8s-api-webhook/authn/provider"
)

func main() {
	cfg := NewConfig()
	cfg.ParseFlags()

	//create authn object
	staticAuthenticator, err := provider.ReadTokensFile(cfg.TokenFile)
	if err != nil {
		log.Fatal(err)
	}
	authnProvider := provider.NewAggregator(staticAuthenticator)
	authnHandler := authn.NewAuthenticationHandler(authnProvider)
	http.Handle("/authentication", authnHandler)
	//create authz object

	//register routes and attach handlers
}
