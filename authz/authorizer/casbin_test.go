package authorizer

var _ ResourceAuthorizer = &ResourceCasbin{}
var _ NonResourceAuthorizer = &NonResourceCasbin{}
