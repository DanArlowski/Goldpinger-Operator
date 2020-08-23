package resources

import (
	"strconv"

	goldpingerv1alpha1 "github.com/bloomberg/goldpinger/api/v1alpha1"
	"github.com/bloomberg/goldpinger/controllers/constants"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// GetDaemonSetResourceDefinition returns a Goldpinger Daemonset object
func GetDaemonSetResourceDefinition(g *goldpingerv1alpha1.Goldpinger, meta metav1.ObjectMeta) *appsv1.DaemonSet {
	ls := LabelsForGoldpinger(meta.Name)
	env := envForGoldpinger(g.Spec.GoldpingerConfig)

	return &appsv1.DaemonSet{
		ObjectMeta: meta,
		Spec: appsv1.DaemonSetSpec{
			UpdateStrategy: appsv1.DaemonSetUpdateStrategy{
				Type: constants.UpdateStrategy,
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: g.Spec.GoldpingerConfig.Annotations,
					Labels:      ls,
				},
				Spec: corev1.PodSpec{
					Tolerations:        g.Spec.GoldpingerConfig.Tolerations,
					ServiceAccountName: g.Name,
					Containers: []corev1.Container{{
						Image: g.Spec.GoldpingerConfig.Image,
						Name:  g.Name,
						Ports: []corev1.ContainerPort{{
							ContainerPort: g.Spec.GoldpingerConfig.Env.Port,
							Name:          g.Name,
						}},
						ImagePullPolicy: constants.ImagePullPolicy,
						Env:             env,
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								"memory": resource.MustParse(constants.LimitMemory),
							},
							Requests: corev1.ResourceList{
								"memory": resource.MustParse(constants.RequestMemory),
								"cpu":    resource.MustParse(constants.RequestCPU),
							},
						},
						ReadinessProbe: &corev1.Probe{
							InitialDelaySeconds: 20,
							PeriodSeconds:       5,
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: constants.ReadinessPath,
									Port: intstr.FromInt(int(g.Spec.GoldpingerConfig.Env.Port)),
								},
							},
						},
						LivenessProbe: &corev1.Probe{
							InitialDelaySeconds: 20,
							PeriodSeconds:       5,
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: constants.LivenessPath,
									Port: intstr.FromInt(int(g.Spec.GoldpingerConfig.Env.Port)),
								},
							},
						},
					}},
				},
			},
		},
	}
}

func envForGoldpinger(gc *goldpingerv1alpha1.GoldpingerConfig) []corev1.EnvVar {
	env := []corev1.EnvVar{{
		Name:  constants.EnvHost,
		Value: gc.Env.Host,
	},
		{
			Name:  constants.EnvPort,
			Value: strconv.Itoa(int(gc.Env.Port)),
		},
		{
			Name:  constants.ENVHostNameToResolve,
			Value: gc.Env.HostsToResolve,
		},
	}

	if gc.Env.HostName != "" {
		env = append(env, corev1.EnvVar{
			Name:  constants.EnvHostName,
			Value: gc.Env.HostName,
		})
	} else {
		env = append(env, corev1.EnvVar{
			Name: constants.EnvHostName,
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "spec.nodeName",
				},
			},
		})
	}
	if gc.Env.PodIP != "" {
		env = append(env, corev1.EnvVar{
			Name:  constants.EnvPodIP,
			Value: gc.Env.PodIP,
		})
	} else {
		env = append(env, corev1.EnvVar{
			Name: constants.EnvPodIP,
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "status.podIP",
				},
			},
		})
	}
	return env
}
