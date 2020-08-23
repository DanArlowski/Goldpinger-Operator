package constants

const (
	//HostPath is the path for the route
	HostPath = ""
	//Image is the Goldpinger Image
	Image = "docker.io/bloomberg/goldpinger:v3.0.0"
	//Port is the port used by the app
	Port int32 = 8080
	//NodePort is the exposed port on the node when deployed on kubernetes
	NodePort = 30800
	//Host Goldpinger should listen to
	Host = "0.0.0.0"
	//HostName to Show
	HostName = ""
	//PodIP is used to selec a subset of nodes to ping
	PodIP = ""
	//HostNameToResolve is the environment variable to ping external hosts
	HostNameToResolve = ""

	//EnvHost is the environment variable name for host
	EnvHost = "HOST"
	//EnvPort is the environment variable name for port
	EnvPort = "PORT"
	//EnvHostName is the environment variable name for HostName
	EnvHostName = "HOSTNAME"
	//EnvPodIP is the environment variable name for PodIP
	EnvPodIP = "POD_IP"
	//ENVHOSTNAME is the environment variable name for HostNames to Resolve
	ENVHostNameToResolve = "HOSTS_TO_RESOLVE"

	//LimitMemory is the pod memory limit
	LimitMemory = "80Mi"
	//RequestMemory is the pod memory request
	RequestMemory = "40Mi"
	//RequestCPU is the pod CPU Request
	RequestCPU = "1m"
	//ImagePullPolicy for the image
	ImagePullPolicy = "Always"

	//ReadinessPath is the path for the readiness probe
	ReadinessPath = "/healthz"
	//LivenessPath is the path for the liveness probe
	LivenessPath = "/healthz"

	//UpdateStrategy is the startegy for updating the Daemonset
	UpdateStrategy = "RollingUpdate"
)
