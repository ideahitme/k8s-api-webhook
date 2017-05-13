# k8s-api-webhook

[![Coverage Status](https://coveralls.io/repos/github/ideahitme/k8s-api-webhook/badge.svg?branch=master)](https://coveralls.io/github/ideahitme/k8s-api-webhook?branch=master)
[![Build Status](https://travis-ci.org/ideahitme/k8s-api-webhook.svg?branch=master)](https://travis-ci.org/ideahitme/k8s-api-webhook)
[![Go Report Card](https://goreportcard.com/badge/github.com/ideahitme/k8s-api-webhook)](https://goreportcard.com/report/github.com/ideahitme/k8s-api-webhook)

Library for quickly bootstrapping authentication/authorisation webhook for Kubernetes API

- [x] Move `main.go` to `cmd/main.go` 
- [x] Think how to restructure packages
- [ ] Provide detailed readme how to extend it 
- [ ] Write authorization module
- [ ] Complete tests
  - [x] Authentication
  - [ ] Authorization
- [x] Get rid of kingpin. as far as i can see this lib can be completely vendor free