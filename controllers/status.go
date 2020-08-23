package controllers

import (
	"context"
	"reflect"

	goldpingerv1alpha1 "github.com/bloomberg/goldpinger/api/v1alpha1"
	"github.com/bloomberg/goldpinger/controllers/resources"
	"github.com/go-logr/logr"
	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//UpdateCRStatus updates the CR status for all the resources
func (r *GoldpingerReconciler) UpdateCRStatus(ctx context.Context, log logr.Logger, err error, req ctrl.Request, goldpinger *goldpingerv1alpha1.Goldpinger, meta metav1.ObjectMeta, isOpenshift bool) (ctrl.Result, error) {

	changed, re, er := r.updatePodsStatus(ctx, log, err, req, goldpinger, meta)
	if changed {
		return ctrl.Result{}, nil
	}
	if err != nil {
		return re, er
	}

	changed, re, er = r.updateServiceStatus(ctx, log, err, req, goldpinger, meta)
	if changed {
		return ctrl.Result{}, nil
	}
	if err != nil {
		return re, er
	}

	if isOpenshift {
		changed, re, er = r.updateRouteStatus(ctx, log, err, req, goldpinger, meta)
		if changed {
			return ctrl.Result{}, nil
		}
		if err != nil {
			return re, er
		}
	}
	return ctrl.Result{}, nil

}

//updatePodsStatus updates the CR pods status
func (r *GoldpingerReconciler) updatePodsStatus(ctx context.Context, log logr.Logger, err error, req ctrl.Request, goldpinger *goldpingerv1alpha1.Goldpinger, meta metav1.ObjectMeta) (bool, ctrl.Result, error) {
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(goldpinger.Namespace),
		client.MatchingLabels(resources.LabelsForGoldpinger(goldpinger.Name)),
	}
	if err = r.List(ctx, podList, listOpts...); err != nil {
		log.Error(err, "Failed to list pods", "Goldpinger.Namespace", goldpinger.Namespace, "Goldpinger.Name", goldpinger.Name)
		return true, ctrl.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, goldpinger.Status.Nodes) {
		goldpinger.Status.Nodes = podNames
		err := r.Status().Update(ctx, goldpinger)
		if err != nil {
			log.Error(err, "Failed to update Goldpinger status")
			return true, ctrl.Result{}, err
		}
	}
	return false, ctrl.Result{}, nil

}

//updateServiceStatus returns the service type deployed
func (r *GoldpingerReconciler) updateServiceStatus(ctx context.Context, log logr.Logger, err error, req ctrl.Request, goldpinger *goldpingerv1alpha1.Goldpinger, meta metav1.ObjectMeta) (bool, ctrl.Result, error) {
	svc := &corev1.Service{}

	if err = r.Get(ctx, types.NamespacedName{Name: goldpinger.Name, Namespace: goldpinger.Namespace}, svc); err != nil {
		log.Error(err, "Failed to get svc", "Goldpinger.Namespace", goldpinger.Namespace, "Goldpinger.Name", goldpinger.Name)
		return true, ctrl.Result{}, err
	}
	svcType := svc.Spec.Type

	if !reflect.DeepEqual(string(svcType), goldpinger.Status.ServiceType) {
		goldpinger.Status.ServiceType = string(svcType)
		err := r.Status().Update(ctx, goldpinger)
		if err != nil {
			log.Error(err, "Failed to update Goldpinger service Status")
			return true, ctrl.Result{}, err
		}
	}
	return false, ctrl.Result{}, nil

}
func (r *GoldpingerReconciler) updateRouteStatus(ctx context.Context, log logr.Logger, err error, req ctrl.Request, goldpinger *goldpingerv1alpha1.Goldpinger, meta metav1.ObjectMeta) (bool, ctrl.Result, error) {
	route := &routev1.Route{}

	if err = r.Get(ctx, types.NamespacedName{Name: goldpinger.Name, Namespace: goldpinger.Namespace}, route); err != nil {
		log.Error(err, "failed to get route", "Goldpinger.Namespace", goldpinger.Namespace, "Goldpinger.Name", goldpinger.Name)
		return true, ctrl.Result{}, err
	}
	routePath := route.Spec.Host
	//Dont look for a custom made comments again Noam!!
	if !reflect.DeepEqual(routePath, goldpinger.Status.RoutePath) {
		goldpinger.Status.RoutePath = routePath
		err := r.Status().Update(ctx, goldpinger)
		if err != nil {
			log.Error(err, "Failed to update Goldpiner route Status")
			return true, ctrl.Result{}, err
		}
	}
	return false, ctrl.Result{}, nil
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}
