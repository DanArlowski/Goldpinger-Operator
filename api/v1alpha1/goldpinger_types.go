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

package v1alpha1

import (
	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GoldpingerSpec defines the desired state of Goldpinger
type GoldpingerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// GoldpingerConfig holds all the Goldpinger Available Configurations.
	GoldpingerConfig *GoldpingerConfig `json:"goldPingerConfig,omitempty"`
}

//GoldpingerConfig holds all the Goldpinger Available Configurations
type GoldpingerConfig struct {
	//the path for the route
	HostPath string `json:"hostPath,omitempty"`

	//the goldpinger Image
	Image string `json:"Image,omitempty"`

	//+optional
	//Tolerations for the Daemonset
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	//+optional
	//NodePort, only used when deployed on kubernetes
	NodePort int32 `json:"nodePort,omitempty"`

	RouteTLS routev1.TLSConfig `json:"routeTLS,omitempty"`

	//+optional
	//Annotations passed by the CR
	Annotations map[string]string `json:"annotations,omitempty"`

	//+optional
	Env *GoldpingerEnv `json:"Env,omitempty"`
}

//GoldpingerEnv holds all the environment variables available for the goldpinger container.
type GoldpingerEnv struct {
	// Host to listen to
	Host string `json:"Host,omitempty"`

	// port for the app
	Port int32 `json:"Port,omitempty"`

	// Hostname to show
	HostName string `json:"HostName,omitempty"`

	// podIP is used to select a randomized subset of nodes to ping.
	PodIP string `json:"PodIP,omitempty"`

	// HostsToResolve is used to ping external Hosts
	HostsToResolve string `json:"HostsToResolve,omitempty"`
}

// GoldpingerStatus defines the observed state of Goldpinger
type GoldpingerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Nodes are the names of the Daemonset pods
	Nodes []string `json:"nodes"`

	//ServiceType is the type of the service deployed
	ServiceType string `json:"serviceType"`

	//RoutePath is the path created by the route
	RoutePath string `json:"routePath,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Goldpinger is the Schema for the goldpingers API
type Goldpinger struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GoldpingerSpec   `json:"spec,omitempty"`
	Status GoldpingerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GoldpingerList contains a list of Goldpinger
type GoldpingerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Goldpinger `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Goldpinger{}, &GoldpingerList{})
}
