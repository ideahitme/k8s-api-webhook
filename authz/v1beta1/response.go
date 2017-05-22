package v1beta1

/**

Successfuly response:
{
  "apiVersion": "authorization.k8s.io/v1beta1",
  "kind": "SubjectAccessReview",
  "status": {
    "allowed": true
  }
}

Unauthorized response:
{
  "apiVersion": "authorization.k8s.io/v1beta1",
  "kind": "SubjectAccessReview",
  "status": {
    "allowed": false,
    "reason": "user does not have read access to the namespace"
  }
}

*/

// ResponseConstructor constructs response to authorization requests according to official requirements of v1beta1 version
type ResponseConstructor struct {
}
