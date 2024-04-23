// apiVersion: deployto.dev/v1beta1
// kind: Target
// meta:
//   name: dev-target
// spec:
//   cluster: file://./kubeconfig
//   namespace: coolapp
// helm:
//   repository: file://./targets-helm/dev
// ---
// apiVersion: deployto.dev/v1beta1
// kind: target
// meta:
//   name: prod-target-geo1
// spec:
//   cluster: file://./kubeconfig2
//   namespace: coolapp
// helm:
//   repository: file://./targets-helm/prod-geo1
// ---
// apiVersion: deployto.dev/v1beta1
// kind: target
// meta:
//   name: prod-target-geo2
// spec:
//   cluster: file://./kubeconfig3
//   namespace: coolapp
// helm:
//   repository: file://./targets-helm/prod-geo2

package types

// TODO
type Target struct {
	Meta `json:",inline" yaml:",inline"`
	Helm *Helm
}
