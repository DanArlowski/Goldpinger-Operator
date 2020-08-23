package resources

import (
	goldpingerv1alpha1 "github.com/bloomberg/goldpinger/api/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// GetGoldpingerRouteDefinition return the route definition for goldpinger
func GetGoldpingerRouteDefinition(meta metav1.ObjectMeta, g *goldpingerv1alpha1.Goldpinger) *routev1.Route {
	return &routev1.Route{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Route",
			APIVersion: routev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: meta,
		Spec: routev1.RouteSpec{
			Host: g.Spec.GoldpingerConfig.HostPath,
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: meta.Name,
			},
			Port: &routev1.RoutePort{
				TargetPort: intstr.FromInt(int(g.Spec.GoldpingerConfig.Env.Port)),
			},
			TLS: &g.Spec.GoldpingerConfig.RouteTLS,
		},
	}

}
