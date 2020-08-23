package controllers

import (
	goldpingerv1alpha1 "github.com/bloomberg/goldpinger/api/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"

	resources "github.com/bloomberg/goldpinger/controllers/resources"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// DaemonsetForGoldpinger returns a Goldpinger Daemonset object
func DaemonsetForGoldpinger(g *goldpingerv1alpha1.Goldpinger, meta metav1.ObjectMeta, r *GoldpingerReconciler) *appsv1.DaemonSet {
	ds := resources.GetDaemonSetResourceDefinition(g, meta)
	// Set Goldpinger instance as the owner and controller
	ctrl.SetControllerReference(g, ds, r.Scheme)
	return ds
}

// ServiceAccountForGoldpinger return a service Account object
func ServiceAccountForGoldpinger(g *goldpingerv1alpha1.Goldpinger, meta metav1.ObjectMeta, r *GoldpingerReconciler) *corev1.ServiceAccount {
	sa := resources.GetServiceAccountResourceDefinition(meta)
	// Set Goldpinger instance as the owner and controller
	ctrl.SetControllerReference(g, sa, r.Scheme)
	return sa
}

// RoleForGoldpinger creates a Role for the Goldpinger Daemonset
func RoleForGoldpinger(g *goldpingerv1alpha1.Goldpinger, meta metav1.ObjectMeta, r *GoldpingerReconciler) *rbacv1.Role {
	role := resources.GetRoleResourceDefinition(meta)
	// Set Goldpinger instance as the owner and controller
	ctrl.SetControllerReference(g, role, r.Scheme)
	return role
}

// RoleBindingForGoldpinger Binds the role and the Service Account
func RoleBindingForGoldpinger(g *goldpingerv1alpha1.Goldpinger, meta metav1.ObjectMeta, r *GoldpingerReconciler) *rbacv1.RoleBinding {
	rolebinding := resources.GetRoleBindingResourceDefinition(meta)
	// Set Goldpinger instance as the owner and controller
	ctrl.SetControllerReference(g, rolebinding, r.Scheme)
	return rolebinding
}

// ServiceForGoldpinger creates a service for goldpinger
func ServiceForGoldpinger(g *goldpingerv1alpha1.Goldpinger, meta metav1.ObjectMeta, r *GoldpingerReconciler, isOpenshift bool) *corev1.Service {
	svc := resources.GetGoldpingerServiceDefinition(meta, g, isOpenshift)
	// Set Goldpinger instance as the owner and controller
	ctrl.SetControllerReference(g, svc, r.Scheme)
	return svc
}

// RouteForGoldpinger exposes a route to connect to the daemonset
func RouteForGoldpinger(g *goldpingerv1alpha1.Goldpinger, meta metav1.ObjectMeta, r *GoldpingerReconciler) *routev1.Route {
	route := resources.GetGoldpingerRouteDefinition(meta, g)
	// Set Goldpinger instance as the owner and controller
	ctrl.SetControllerReference(g, route, r.Scheme)
	return route
}
