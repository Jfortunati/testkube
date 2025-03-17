/*
 * Testkube API
 *
 * Testkube provides a Kubernetes-native framework for test definition, execution and results
 *
 * API version: 1.0.0
 * Contact: testkube@kubeshop.io
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package testkube

// Projection that may be projected along with other supported volume types. Exactly one of these fields must be set.
type ProjectedVolumeSourceSources struct {
	ClusterTrustBundle  *ProjectedVolumeSourceClusterTrustBundle  `json:"clusterTrustBundle,omitempty"`
	ConfigMap           *ProjectedVolumeSourceConfigMap           `json:"configMap,omitempty"`
	DownwardAPI         *ProjectedVolumeSourceDownwardApi         `json:"downwardAPI,omitempty"`
	Secret              *ProjectedVolumeSourceSecret              `json:"secret,omitempty"`
	ServiceAccountToken *ProjectedVolumeSourceServiceAccountToken `json:"serviceAccountToken,omitempty"`
}
