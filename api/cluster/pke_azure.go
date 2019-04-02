// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	intPKE "github.com/banzaicloud/pipeline/internal/pke"
	"github.com/banzaicloud/pipeline/internal/providers/azure/pke"
	"github.com/banzaicloud/pipeline/internal/providers/azure/pke/driver"
)

const PKEOnAzure = pke.PKEOnAzure

type NodePool struct {
	Labels map[string]string `json:"labels,omitempty"`
	Name   string            `json:"name" binding:"required"`
	Roles  []string          `json:"roles" binding:"required"`

	Subnet AzureSubnet `json:"subnet"`
	Zones  []string    `json:"zones"`

	InstanceType string `json:"instanceType"`

	Autoscaling bool `json:"autoscaling"`
	Count       int  `json:"count"`
	Min         int  `json:"min"`
	Max         int  `json:"max"`
}

type AzureSubnet struct {
	Name string `json:"name"`
	CIDR string `json:"cidr"`
}

type AzureNetwork struct {
	Name string `json:"name"`
	CIDR string `json:"cidr"`
}

type CreatePKEOnAzureClusterRequest struct {
	CreateClusterRequestBase
	Location      string       `json:"location"`
	ResourceGroup string       `json:"resourceGroup"`
	NodePools     []NodePool   `json:"nodepools,omitempty" binding:"required"`
	Kubernetes    Kubernetes   `json:"kubernetes,omitempty" binding:"required"`
	Network       AzureNetwork `json:"network"`
}

func (req CreatePKEOnAzureClusterRequest) ToAzurePKEClusterCreationParams(organizationID, userID uint) driver.AzurePKEClusterCreationParams {
	nodepools := make([]driver.NodePool, len(req.NodePools))
	for i, node := range req.NodePools {
		nodepools[i] = driver.NodePool{
			CreatedBy:    userID,
			Name:         node.Name,
			InstanceType: node.InstanceType,
			Subnet: driver.Subnet{
				Name: node.Subnet.Name,
				CIDR: node.Subnet.CIDR,
			},
			Zones:       node.Zones,
			Roles:       node.Roles,
			Labels:      node.Labels,
			Autoscaling: node.Autoscaling,
			Count:       node.Count,
			Min:         node.Min,
			Max:         node.Max,
		}
	}
	return driver.AzurePKEClusterCreationParams{
		Name:           req.Name,
		OrganizationID: organizationID,
		CreatedBy:      userID,
		ResourceGroup:  req.ResourceGroup,
		SecretID:       req.SecretID,
		SSHSecretID:    req.SSHSecretID,
		Kubernetes: intPKE.Kubernetes{
			Version: req.Kubernetes.Version,
			RBAC:    req.Kubernetes.RBAC,
			Network: intPKE.Network{
				ServiceCIDR:    req.Kubernetes.Network.ServiceCIDR,
				PodCIDR:        req.Kubernetes.Network.PodCIDR,
				Provider:       req.Kubernetes.Network.Provider,
				ProviderConfig: req.Kubernetes.Network.ProviderConfig,
			},
			CRI: intPKE.CRI{
				Runtime:       req.Kubernetes.CRI.Runtime,
				RuntimeConfig: req.Kubernetes.CRI.RuntimeConfig,
			},
		},
		Network: driver.VirtualNetwork{
			Name:     req.Network.Name,
			CIDR:     req.Network.CIDR,
			Location: req.Location,
		},
		NodePools: nodepools,
	}
}
