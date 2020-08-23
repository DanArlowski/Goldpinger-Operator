package resources

// LabelsForGoldpinger returns the labels for selecting the resources
// belonging to the given goldpinger CR name.
func LabelsForGoldpinger(name string) map[string]string {
	return map[string]string{"app": "goldpinger", "goldpinger_cr": name}
}
