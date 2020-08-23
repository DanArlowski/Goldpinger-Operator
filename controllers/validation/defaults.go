package validation

import (
	"reflect"

	goldpingerv1alpha1 "github.com/bloomberg/goldpinger/api/v1alpha1"
	"github.com/bloomberg/goldpinger/controllers/constants"
	"github.com/go-logr/logr"
	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
)

// SetDefaults generates the default Values not passed by the user, Triggered from logic.go, when needed in reconcile loop.
func SetDefaults(goldpinger *goldpingerv1alpha1.Goldpinger, log logr.Logger) {
	found := true

	if !found || goldpinger.Spec.GoldpingerConfig == nil {
		log.Info("goldpingerConfig not found, setting default")
		goldpinger.Spec.GoldpingerConfig = &goldpingerv1alpha1.GoldpingerConfig{}
		found = false
	}

	if goldpinger.Spec.GoldpingerConfig.Tolerations == nil {
		goldpinger.Spec.GoldpingerConfig.Tolerations = []corev1.Toleration{}
		log.Info("No tolerations specified.")
	} else {
		log.Info("Using specified tolerations")
	}

	if goldpinger.Spec.GoldpingerConfig.Annotations == nil {
		log.Info("No Annotations, using default...")
		goldpinger.Spec.GoldpingerConfig.Annotations = getDefaultAnnotationsForGoldpinger(goldpinger)
	} else {
		log.Info("using specified Annotations...")
	}

	if !found || isZeroOfUnderlyingType(goldpinger.Spec.GoldpingerConfig.Image) {
		goldpinger.Spec.GoldpingerConfig.Image = constants.Image
	} else {
		log.Info("found value of Image")
	}

	if !found || isZeroOfUnderlyingType(goldpinger.Spec.GoldpingerConfig.HostPath) {
		goldpinger.Spec.GoldpingerConfig.HostPath = constants.HostPath
	} else {
		log.Info("found value of HostPath")
	}

	if !found || isZeroOfUnderlyingType(goldpinger.Spec.GoldpingerConfig.NodePort) {
		goldpinger.Spec.GoldpingerConfig.NodePort = constants.NodePort
	} else {
		log.Info("found value of NodePort")
	}

	if !found || isZeroOfUnderlyingType(goldpinger.Spec.GoldpingerConfig.RouteTLS) {
		goldpinger.Spec.GoldpingerConfig.RouteTLS = routev1.TLSConfig{}
	} else {
		log.Info("found RouteTLS definition")
	}
	if !found || goldpinger.Spec.GoldpingerConfig.Env == nil {
		goldpinger.Spec.GoldpingerConfig.Env = &goldpingerv1alpha1.GoldpingerEnv{}
		log.Info("goldpingerEnv not found, setting default")
		found = false
	}

	if !found || isZeroOfUnderlyingType(goldpinger.Spec.GoldpingerConfig.Env.Host) {
		goldpinger.Spec.GoldpingerConfig.Env.Host = constants.Host
	} else {
		log.Info("found value of Env.Host:")
	}
	if !found || isZeroOfUnderlyingType(goldpinger.Spec.GoldpingerConfig.Env.HostName) {
		goldpinger.Spec.GoldpingerConfig.Env.HostName = constants.HostName
	} else {
		log.Info("found value of Env.HostName")
	}
	if !found || isZeroOfUnderlyingType(goldpinger.Spec.GoldpingerConfig.Env.PodIP) {
		goldpinger.Spec.GoldpingerConfig.Env.PodIP = constants.PodIP
	} else {
		log.Info("found value of Env.PodIP")
	}
	if !found || isZeroOfUnderlyingType(goldpinger.Spec.GoldpingerConfig.Env.Port) {
		goldpinger.Spec.GoldpingerConfig.Env.Port = constants.Port
	} else {
		log.Info("found value of Env.Port")
	}
	if !found || isZeroOfUnderlyingType(goldpinger.Spec.GoldpingerConfig.Env.HostsToResolve) {
		goldpinger.Spec.GoldpingerConfig.Env.HostsToResolve = constants.HostNameToResolve
	} else {
		log.Info("found value of Env.HostsToResolve")
	}
}

func isZeroOfUnderlyingType(x interface{}) bool {
	return reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

//GetDefaultAnnotationsForGoldpinger returns default map if not passed by the CR.
func getDefaultAnnotationsForGoldpinger(g *goldpingerv1alpha1.Goldpinger) map[string]string {
	return map[string]string{"prometheus.io/scrape": "true", "prometheus.io/port": string(g.Spec.GoldpingerConfig.Env.Port)}
}
