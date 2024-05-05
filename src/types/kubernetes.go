package types

import (
	corev1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Define the structure for Kubernetes Service
type Service struct {
	Kind       string
	APIVersion string `yaml:"apiVersion"`
	Metadata   metav1.ObjectMeta
	Spec       corev1.ServiceSpec
}

// Define the structure for Kubernetes Ingress
type Ingress struct {
	Kind       string
	APIVersion string `yaml:"apiVersion"`
	Metadata   metav1.ObjectMeta
	Spec       networking.IngressSpec
}
