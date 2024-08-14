// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build kubeapiserver

// Package alwaysadmit is a validation webhook that allows all pods into the cluster.
// Its behavior is the same as if there were no validation at all.
package alwaysadmit

import (
	"github.com/DataDog/datadog-agent/cmd/cluster-agent/admission"
	"github.com/DataDog/datadog-agent/pkg/clusteragent/admission/common"
	validatecommon "github.com/DataDog/datadog-agent/pkg/clusteragent/admission/validate/common"
	admiv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/dynamic"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	webhookName     = "always_admit"
	webhookEndpoint = "/always-admit"
)

// Webhook is a validation webhook that allows all pods into the cluster.
type Webhook struct {
	name        string
	webhookType common.WebhookType
	isEnabled   bool
	endpoint    string
	resources   []string
	operations  []admissionregistrationv1.OperationType
}

// NewWebhook returns a new webhook
func NewWebhook() *Webhook {
	return &Webhook{
		name:        webhookName,
		webhookType: common.ValidatingWebhook,
		isEnabled:   true,
		endpoint:    webhookEndpoint,
		resources:   []string{"pods"},
		operations:  []admissionregistrationv1.OperationType{admissionregistrationv1.Create},
	}
}

// Name returns the name of the webhook
func (w *Webhook) Name() string {
	return w.name
}

// WebhookType returns the type of the webhook
func (w *Webhook) WebhookType() common.WebhookType {
	return w.webhookType
}

// IsEnabled returns whether the webhook is enabled
func (w *Webhook) IsEnabled() bool {
	return w.isEnabled
}

// Endpoint returns the endpoint of the webhook
func (w *Webhook) Endpoint() string {
	return w.endpoint
}

// Resources returns the kubernetes resources for which the webhook should
// be invoked
func (w *Webhook) Resources() []string {
	return w.resources
}

// Operations returns the operations on the resources specified for which
// the webhook should be invoked
func (w *Webhook) Operations() []admissionregistrationv1.OperationType {
	return w.operations
}

// LabelSelectors returns the label selectors that specify when the webhook
// should be invoked
func (w *Webhook) LabelSelectors(useNamespaceSelector bool) (namespaceSelector *metav1.LabelSelector, objectSelector *metav1.LabelSelector) {
	return common.DefaultLabelSelectors(useNamespaceSelector)
}

// WebhookFunc returns the function that validate the resources
func (w *Webhook) WebhookFunc() func(request *admission.Request) *admiv1.AdmissionResponse {
	return func(request *admission.Request) *admiv1.AdmissionResponse {
		return common.ValidationResponse(validatecommon.Validate(request.Raw, request.Namespace, w.Name(), w.alwaysAdmit, request.DynamicClient))
	}
}

// alwaysAdmit is a function that always admits the pod.
func (w *Webhook) alwaysAdmit(_ *corev1.Pod, _ string, _ dynamic.Interface) (bool, error) {
	return true, nil
}
