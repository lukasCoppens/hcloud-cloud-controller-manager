package hcloud

import (
	"github.com/appscode/go-hetzner"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

// HetznerClient is a client that uses the Hcloud API and Hetzner Robot API
type HetznerClient struct {
	cloudClient *hcloud.Client
	robotClient *hetzner.Client
}

// GetServerByProviderID returns the server based on the providerID. The providerID contains the server ID
// returns cloudprovider.InstanceNotFound if the instance is not found.
func (hc *HetznerClient) GetServerByProviderID(providerID string) (*Server, error) {
	serverID, err := providerIDToServerID(providerID)
	if err != nil {
		return nil, err
	}
	server, err := getCloudServerByID(hc.cloudClient, serverID)
	if err != cloudprovider.InstanceNotFound {
		if err != nil {
			return nil, err
		}
		return convertCloudServerToServer(server), nil
	}
	if hc.robotClient != nil {
		server, err := getRobotServerByID(hc.robotClient, serverID)
		if err != nil {
			return nil, err
		}
		return convertRobotServerToServer(server), nil
	}
	return nil, cloudprovider.InstanceNotFound
}

// GetServerByName returns the server based on the server name. We check hcloud first so if there is a match this will be returned.
// returns cloudprovider.InstanceNotFound if the instance is not found.
func (hc *HetznerClient) GetServerByName(name string) (*Server, error) {
	server, err := getCloudServerByName(hc.cloudClient, name)
	if err != cloudprovider.InstanceNotFound {
		if err != nil {
			return nil, err
		}
		return convertCloudServerToServer(server), nil
	}
	if hc.robotClient != nil {
		server, err := getRobotServerByName(hc.robotClient, name)
		if err != nil {
			return nil, err
		}
		return convertRobotServerToServer(server), nil
	}
	return nil, cloudprovider.InstanceNotFound
}

// Server is an internal struct to contain the data
type Server struct {
	Name          string
	Ipv4          string
	Region        string
	Failuredomain string
	InstanceType  string
	ID            int
}
