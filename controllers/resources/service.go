package resources

import (
	goldpingerv1alpha1 "github.com/bloomberg/goldpinger/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// GetGoldpingerServiceDefinition return Goldpinger Service
func GetGoldpingerServiceDefinition(meta metav1.ObjectMeta, g *goldpingerv1alpha1.Goldpinger, isOpenshift bool) *corev1.Service {
	serviceType := getServiceType(isOpenshift)
	servicePort := getRelevantServicePort(isOpenshift, g)
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: meta,
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(serviceType),

			ClusterIP: "",
			Selector:  LabelsForGoldpinger(meta.Name),
			Ports: []corev1.ServicePort{
				servicePort,
			},
		},
	}

}

func getServiceType(isOpenshift bool) string {
	if isOpenshift {
		return "ClusterIP"
	}
	return "NodePort"
}
func getRelevantServicePort(isOpenshift bool, g *goldpingerv1alpha1.Goldpinger) corev1.ServicePort {
	if isOpenshift {
		return corev1.ServicePort{
			Port:       g.Spec.GoldpingerConfig.Env.Port,
			Protocol:   "TCP",
			TargetPort: intstr.FromInt(int(g.Spec.GoldpingerConfig.Env.Port)),
		}
	}
	return corev1.ServicePort{
		TargetPort: intstr.FromInt(int(g.Spec.GoldpingerConfig.Env.Port)),
		Protocol:   "TCP",
		NodePort:   g.Spec.GoldpingerConfig.NodePort,
		Port:       g.Spec.GoldpingerConfig.Env.Port,
	}
}
