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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HardwareReqType is used to enumerate the various hardware requests that can be made for the set
type HardwareReqType string

// BaremetalSetSpec defines the desired state of BaremetalSet
type BaremetalSetSpec struct {
	// Replicas The number of baremetalhosts to attempt to aquire
	// +kubebuilder:validation:Minimum=0
	Replicas int `json:"replicas,omitempty"`
	// Remote URL pointing to desired RHEL qcow2 image
	RhelImageURL string `json:"rhelImageUrl,omitempty"`
	// ProvisionServerName Optional. If supplied will be used as the base Image for the baremetalset instead of RhelImageURL.
	ProvisionServerName string `json:"provisionServerName,omitempty"`
	// Name of secret holding the stack-admin ssh keys
	DeploymentSSHSecret string `json:"deploymentSSHSecret"`
	// Interface to use for ctlplane network
	CtlplaneInterface string `json:"ctlplaneInterface"`
	// BmhLabelSelector allows for a sub-selection of BaremetalHosts based on arbitrary labels
	BmhLabelSelector map[string]string `json:"bmhLabelSelector,omitempty"`
	// Hardware requests for sub-selection of BaremetalHosts with certain hardware specs
	HardwareReqs HardwareReqs `json:"hardwareReqs,omitempty"`
	// Networks the name(s) of the OvercloudNetworks used to generate IPs
	Networks []string `json:"networks"`
	// Role the name of the Overcloud role this IPset is associated with. Used to generate hostnames.
	Role string `json:"role"`
	// PasswordSecret the name of the secret used to optionally set the root pwd by adding
	// NodeRootPassword: <base64 enc pwd>
	// to the secret data
	PasswordSecret string `json:"passwordSecret,omitempty"`
}

// BaremetalSetStatus defines the observed state of BaremetalSet
type BaremetalSetStatus struct {
	BaremetalHosts map[string]BaremetalHostStatus `json:"baremetalHosts"`
}

// BaremetalHostStatus represents the observed state of a particular allocated BaremetalHost resource
type BaremetalHostStatus struct {
	Hostname              string `json:"hostname"`
	UserDataSecretName    string `json:"userDataSecretName"`
	NetworkDataSecretName string `json:"networkDataSecretName"`
	CtlplaneIP            string `json:"ctlplaneIP"`
	Online                bool   `json:"online"`
}

// GetHostnames -
func (bmSet OpenStackBaremetalSet) GetHostnames() map[string]string {
	ret := make(map[string]string)
	for key, val := range bmSet.Status.BaremetalHosts {
		ret[key] = val.Hostname
	}
	return ret
}

// HardwareReqs defines request hardware attributes for the BaremetalHost replicas
type HardwareReqs struct {
	CPUReqs  CPUReqs  `json:"cpuReqs,omitempty"`
	MemReqs  MemReqs  `json:"memReqs,omitempty"`
	DiskReqs DiskReqs `json:"diskReqs,omitempty"`
}

// CPUReqs defines specific CPU hardware requests
type CPUReqs struct {
	// Arch is a scalar (string) because it wouldn't make sense to give it an "exact-match" option
	// Can be either "x86_64" or "ppc64le" if included
	// +kubebuilder:validation:Enum=x86_64;ppc64le
	Arch     string      `json:"arch,omitempty"`
	CountReq CPUCountReq `json:"countReq,omitempty"`
	MhzReq   CPUMhzReq   `json:"mhzReq,omitempty"`
}

// CPUCountReq defines a specific hardware request for CPU core count
type CPUCountReq struct {
	// +kubebuilder:validation:Minimum=1
	Count int `json:"count,omitempty"`
	// If ExactMatch == false, actual count > Count will match
	ExactMatch bool `json:"exactMatch,omitempty"`
}

// CPUMhzReq defines a specific hardware request for CPU clock speed
type CPUMhzReq struct {
	// +kubebuilder:validation:Minimum=1
	Mhz int `json:"mhz,omitempty"`
	// If ExactMatch == false, actual mhz > Mhz will match
	ExactMatch bool `json:"exactMatch,omitempty"`
}

// MemReqs defines specific memory hardware requests
type MemReqs struct {
	GbReq MemGbReq `json:"gbReq,omitempty"`
}

// MemGbReq defines a specific hardware request for memory size
type MemGbReq struct {
	// +kubebuilder:validation:Minimum=1
	Gb int `json:"gb,omitempty"`
	// If ExactMatch == false, actual GB > Gb will match
	ExactMatch bool `json:"exactMatch,omitempty"`
}

// DiskReqs defines specific disk hardware requests
type DiskReqs struct {
	GbReq DiskGbReq `json:"gbReq,omitempty"`
	// SSD is scalar (bool) because it wouldn't make sense to give it an "exact-match" option
	SSDReq DiskSSDReq `json:"ssdReq,omitempty"`
}

// DiskGbReq defines a specific hardware request for disk size
type DiskGbReq struct {
	// +kubebuilder:validation:Minimum=1
	Gb int `json:"gb,omitempty"`
	// If ExactMatch == false, actual GB > Gb will match
	ExactMatch bool `json:"exactMatch,omitempty"`
}

// DiskSSDReq defines a specific hardware request for disk of type SSD (true) or rotational (false)
type DiskSSDReq struct {
	SSD bool `json:"ssd,omitempty"`
	// We only actually care about SSD flag if it is true or ExactMatch is set to true.
	// This second flag is necessary as SSD's bool zero-value (false) is indistinguishable
	// from it being explicitly set to false
	ExactMatch bool `json:"exactMatch,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// OpenStackBaremetalSet represent a set of baremetal hosts for a specific role within the Overcloud deployment
type OpenStackBaremetalSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BaremetalSetSpec   `json:"spec,omitempty"`
	Status BaremetalSetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OpenStackBaremetalSetList contains a list of BaremetalSet
type OpenStackBaremetalSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenStackBaremetalSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OpenStackBaremetalSet{}, &OpenStackBaremetalSetList{})
}
