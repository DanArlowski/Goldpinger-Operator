/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	goldpingerv1alpha1 "github.com/bloomberg/goldpinger/api/v1alpha1"
	"github.com/go-logr/logr"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GoldpingerReconciler reconciles a Goldpinger object
type GoldpingerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=goldpinger.bloomberg.com,resources=goldpingers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=goldpinger.bloomberg.com,resources=goldpingers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;
func (r *GoldpingerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("goldpinger", req.NamespacedName)
	log.Info("Reconcile Loop...")
	goldpinger := &goldpingerv1alpha1.Goldpinger{}

	err := r.Get(ctx, req.NamespacedName, goldpinger)

	isOpenshift := r.isClusterOpenshift(ctx, log, err)

	meta := metav1.ObjectMeta{
		Name:      goldpinger.Name,
		Namespace: goldpinger.Namespace,
	}

	if goldpinger.Namespace == "" {
		log.Info("Namespace is not set, probably an object deletion.")
		return ctrl.Result{}, nil
	}

	SAfound := &corev1.ServiceAccount{}
	err = r.Get(ctx, types.NamespacedName{Name: goldpinger.Name, Namespace: goldpinger.Namespace}, SAfound)
	if err != nil && errors.IsNotFound(err) {
		result, Rerror := r.createServiceAccount(ctx, log, err, req, goldpinger, meta)
		if Rerror == nil {
			return result, Rerror
		}
	}
	found := &appsv1.DaemonSet{}
	err = r.Get(ctx, types.NamespacedName{Name: goldpinger.Name, Namespace: goldpinger.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		result, Rerror := r.createDaemonSet(ctx, log, err, req, goldpinger, meta)
		if Rerror == nil {
			return result, Rerror
		}
	}

	svcFound := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: goldpinger.Name, Namespace: goldpinger.Namespace}, svcFound)
	if err != nil && errors.IsNotFound(err) {
		result, Rerror := r.createService(ctx, log, err, req, goldpinger, meta, isOpenshift)
		if Rerror == nil {
			return result, Rerror
		}
	}
	if isOpenshift {
		log.Info("deploying route for openshift")
		routeFound := &routev1.Route{}
		err = r.Get(ctx, types.NamespacedName{Name: goldpinger.Name, Namespace: goldpinger.Namespace}, routeFound)
		if err != nil && errors.IsNotFound(err) {
			result, Rerror := r.createRoute(ctx, log, err, req, goldpinger, meta)
			if Rerror == nil {
				return result, Rerror
			}
		}
	}

	return r.UpdateCRStatus(ctx, log, err, req, goldpinger, meta, isOpenshift)
}

func (r *GoldpingerReconciler) isClusterOpenshift(ctx context.Context, log logr.Logger, err error) bool {
	openshift := &corev1.ServiceAccount{}
	err = r.Get(ctx, types.NamespacedName{Name: "default", Namespace: "openshift"}, openshift)
	if err != nil {
		return false
	}
	return true
}

//SetupWithManager setup the reconciler to the manager
func (r *GoldpingerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&goldpingerv1alpha1.Goldpinger{}).
		Complete(r)
}
