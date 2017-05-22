package main

import (
	"fmt"

	"github.com/casbin/casbin"
)

func main() {
	//test non resource access control
	nonResource := casbin.NewEnforcer("./non-resource/non-resource.conf")
	fmt.Println(nonResource.Enforce("kubelet", "", "/api", "post"))         //true
	fmt.Println(nonResource.Enforce("ideahitme", "", "/metrics", "delete")) //true
	fmt.Println(nonResource.Enforce("yerken", "", "/api", "post"))          //true
	fmt.Println(nonResource.Enforce("yerken", "", "/api", "get"))           //true
	fmt.Println(nonResource.Enforce("yerken", "", "/api/print", "get"))     //true
	fmt.Println(nonResource.Enforce("yerken", "", "/metrics", "get"))       //false

	//group
	fmt.Println(nonResource.Enforce("yerken", "PowerUser", "/metrics", "delete")) //true

	//test resource based access control
	resource := casbin.NewEnforcer("./resource/resource.conf")
	fmt.Println(resource.Enforce("kubelet", "", "kube-system", "pods", "list"))
	fmt.Println(resource.Enforce("yerken", "", "kube-system", "pods", "list"))
	fmt.Println(resource.Enforce("yerken", "", "kube-system", "pods", "delete"))
}
