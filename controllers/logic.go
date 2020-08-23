package controllers

import (
	"context"
	"time"

	goldpingerv1alpha1 "github.com/bloomberg/goldpinger/api/v1alpha1"
	"github.com/bloomberg/goldpinger/controllers/validation"
	"github.com/go-logr/logr"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *GoldpingerReconciler) createServiceAccount(ctx context.Context, log logr.Logger, err error, req ctrl.Request, goldpinger *goldpingerv1alpha1.Goldpinger, meta metav1.ObjectMeta) (ctrl.Result, error) {
	SAfound := &corev1.ServiceAccount{}
	err = r.Get(ctx, types.NamespacedName{Name: goldpinger.Name, Namespace: goldpinger.Namespace}, SAfound)
	if err != nil && errors.IsNotFound(err) {
		//create a new service account
		sa := ServiceAccountForGoldpinger(goldpinger, meta, r)
		log.Info("Creating new ServiceAccount", "ServiceAccount.Namespace", sa.Namespace, "ServiceAccount.Name", sa.Name)
		err = r.Create(ctx, sa)
		if err != nil {
			log.Error(err, "Failed to create new ServiceAccount", "ServiceAccount.Namespace", sa.Namespace, "ServiceAccount.Name", sa.Name)
			return ctrl.Result{}, err
		}
		// create a new Role
		role := RoleForGoldpinger(goldpinger, meta, r)
		log.Info("Creating new ClusterRole", "ClusterRole.Namespace", role.Namespace, "ClusterRole.Name", role.Name)
		err = r.Create(ctx, role)
		if err != nil {
			log.Error(err, "Failed to create new ClusterRole", "ClusterRole.Namespace", role.Namespace, "ClusterRole.Name", role.Name)
			return ctrl.Result{}, err
		}
		//create a new RoleBinding
		rolebind := RoleBindingForGoldpinger(goldpinger, meta, r)
		log.Info("Creating new rolebind", "rolebind.Namespace", rolebind.Namespace, "rolebind.Name", rolebind.Name)
		err = r.Create(ctx, rolebind)
		if err != nil {
			log.Error(err, "Failed to create new rolebind", "rolebind.Namespace", rolebind.Namespace, "rolebind.Name", rolebind.Name)
			return ctrl.Result{}, err
		}
		//ServiceAccount and security created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get ServiceAccount")
		return ctrl.Result{}, err
	}
	return ctrl.Result{Requeue: true}, nil
}

func (r *GoldpingerReconciler) createDaemonSet(ctx context.Context, log logr.Logger, err error, req ctrl.Request, goldpinger *goldpingerv1alpha1.Goldpinger, meta metav1.ObjectMeta) (ctrl.Result, error) {

	found := &appsv1.DaemonSet{}
	validation.SetDefaults(goldpinger, log)
	err = r.Get(ctx, types.NamespacedName{Name: goldpinger.Name, Namespace: goldpinger.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new daemonset
		ds := DaemonsetForGoldpinger(goldpinger, meta, r)
		log.Info("Creating a new Daemonset", "Daemonset.Namespace", ds.Namespace, "Daemonset.Name", ds.Name)
		err = r.Create(ctx, ds)
		if err != nil {
			log.Error(err, "Failed to create new Daemonset", "Daemonset.Namespace", ds.Namespace, "Daemonset.Name", ds.Name)
			return ctrl.Result{}, err
		}
		// Daemonset created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Daemonset")
		return ctrl.Result{}, err
	}
	return ctrl.Result{Requeue: true}, nil
}

func (r *GoldpingerReconciler) createService(ctx context.Context, log logr.Logger, err error, req ctrl.Request, goldpinger *goldpingerv1alpha1.Goldpinger, meta metav1.ObjectMeta, isOpenshift bool) (ctrl.Result, error) {
	svcFound := &corev1.Service{}
	validation.SetDefaults(goldpinger, log)
	err = r.Get(ctx, types.NamespacedName{Name: goldpinger.Name, Namespace: goldpinger.Namespace}, svcFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new service
		svc := ServiceForGoldpinger(goldpinger, meta, r, isOpenshift)
		log.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		err = r.Create(ctx, svc)
		if err != nil {
			log.Error(err, "Failed to create new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
			return ctrl.Result{}, err
		}
		// Service created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Service")
		return ctrl.Result{}, err
	}
	return ctrl.Result{Requeue: true}, nil
}

func (r *GoldpingerReconciler) createRoute(ctx context.Context, log logr.Logger, err error, req ctrl.Request, goldpinger *goldpingerv1alpha1.Goldpinger, meta metav1.ObjectMeta) (ctrl.Result, error) {
	routeFound := &routev1.Route{}
	validation.SetDefaults(goldpinger, log)
	time.Sleep(2 * time.Second)
	err = r.Get(ctx, types.NamespacedName{Name: goldpinger.Name, Namespace: goldpinger.Namespace}, routeFound)
	if err != nil && errors.IsNotFound(err) {
		// Define a new Route
		route := RouteForGoldpinger(goldpinger, meta, r)
		log.Info("Creating a new Route", "Route.Namespace", route.Namespace, "Route.Name", route.Name)
		err = r.Create(ctx, route)
		if err != nil {
			log.Error(err, "Failed to create new Route", "Route.Namespace", route.Namespace, "Route.Name", route.Name)
			return ctrl.Result{RequeueAfter: time.Second * 5}, err
		}
		// Route created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Route")
		return ctrl.Result{}, err
	}

	return ctrl.Result{Requeue: true}, nil
}
