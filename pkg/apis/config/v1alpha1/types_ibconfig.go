package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type InfobloxConfig struct {
	metav1.TypeMeta `json:"inline"`

	// View is the view in which the default DNS details are listed
	View string `json:"view,omitempty"`
}