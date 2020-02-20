package hcloud

import (
	"time"

	"github.com/appscode/go-hetzner"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/sirupsen/logrus"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

// HetznerClient is a client that uses the Hcloud API and Hetzner Robot API
type HetznerClient struct {
	cloudClient          *hcloud.Client
	robotClient          *hetzner.Client
	robotCache           []*hetzner.ServerSummary
	robotCacheLastUpdate time.Time
}

// GetServerByProviderID returns the server based on the providerID. The providerID contains the server ID
// returns cloudprovider.InstanceNotFound if the instance is not found.
func (hc *HetznerClient) GetServerByProviderID(providerID string) (*Server, error) {
	logrus.Debug(" returning server by provider id")
	serverID, err := providerIDToServerID(providerID)
	if err != nil {
		logrus.Errorf("failed parsing provider id to server id")
		return nil, err
	}
	logrus.Debug("Try cloud server")
	server, err := hc.getCloudServerByID(serverID)
	if err != cloudprovider.InstanceNotFound {
		if err != nil {
			logrus.Debugf("error getting cloud servers: %v", err)
			return nil, err
		}
		return convertCloudServerToServer(server), nil
	}
	logrus.Debug("Try robot")
	if hc.robotClient != nil {
		server, err := hc.getRobotServerByID(serverID)
		if err != cloudprovider.InstanceNotFound {
			if err != nil {
				logrus.Debugf("error getting robot servers: %v", err)
				return nil, err
			}
			return convertRobotServerToServer(server), nil
		}
	}
	logrus.Debug("Server not found @Hetzner, still sending empty server object to prevent deleting non Hetzner nodes")
	return &Server{}, nil
	// return nil, cloudprovider.InstanceNotFound
}

// GetServerByName returns the server based on the server name. We check hcloud first so if there is a match this will be returned.
// returns cloudprovider.InstanceNotFound if the instance is not found.
func (hc *HetznerClient) GetServerByName(name string) (*Server, error) {
	logrus.Debug("Try getting server by name")
	logrus.Debug("Trying hcloud")
	server, err := hc.getCloudServerByName(name)
	if err != cloudprovider.InstanceNotFound {
		if err != nil {
			logrus.Debugf("error getting cloud servers: %v", err)
			return nil, err
		}
		return convertCloudServerToServer(server), nil
	}
	logrus.Debug("Trying Hetzner robot")
	if hc.robotClient != nil {
		server, err := hc.getRobotServerByName(name)
		if err != cloudprovider.InstanceNotFound {
			if err != nil {
				logrus.Debugf("error getting robot servers: %v", err)
				return nil, err
			}
			return convertRobotServerToServer(server), nil
		}
	}
	logrus.Debug("Server not found @Hetzner, still sending empty server object to prevent deleting non Hetzner nodes")
	return &Server{Name: name}, nil
	// return nil, cloudprovider.InstanceNotFound
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
