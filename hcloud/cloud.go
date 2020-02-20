/*
Copyright 2018 Hetzner Cloud GmbH.

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

package hcloud

import (
	"fmt"
	"io"
	"os"

	"github.com/appscode/go-hetzner"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"github.com/sirupsen/logrus"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
)

const (
	hcloudTokenENVVar    = "HCLOUD_TOKEN"
	hcloudEndpointENVVar = "HCLOUD_ENDPOINT"
	hrobotUsername       = "HROBOT_USERNAME"
	hrobotPassword       = "HROBOT_PASSWORD"
	nodeNameENVVar       = "NODE_NAME"
	providerName         = "hetzner"
	logLevel             = "LOG_LEVEL"
)

type cloud struct {
	client    *HetznerClient
	instances cloudprovider.Instances
	zones     cloudprovider.Zones
}

func newCloud(config io.Reader) (cloudprovider.Interface, error) {
	logrus.Debug("Creating new cloud-provider")
	logrus.Debug("newCloud - getting hcloud token from env")
	token := os.Getenv(hcloudTokenENVVar)
	if token == "" {
		return nil, fmt.Errorf("environment variable %q is required", hcloudTokenENVVar)
	}
	logrus.Debug("newCloud - getting node name from env")
	nodeName := os.Getenv(nodeNameENVVar)
	if nodeName == "" {
		return nil, fmt.Errorf("environment variable %q is required", nodeNameENVVar)
	}
	logrus.Debugf("newCloud - nodeName: %v", nodeName)
	logrus.Debug("newCloud - setting options to create cloud client")
	opts := []hcloud.ClientOption{
		hcloud.WithToken(token),
	}
	logrus.Debug("newCloud - getting endpoint from env vars")
	if endpoint := os.Getenv(hcloudEndpointENVVar); endpoint != "" {
		opts = append(opts, hcloud.WithEndpoint(endpoint))
	}
	logrus.Debug("newCloud - creating new hcloud client")
	client := hcloud.NewClient(opts...)

	logrus.Debug("newCloud - getting Hetzner robot creds")
	hetznerRobotUsername := os.Getenv(hrobotUsername)
	hetznerRobotPassword := os.Getenv(hrobotPassword)
	var robotClient *hetzner.Client
	logrus.Debug("newCloud - creating Hetzner robot client")
	if hetznerRobotUsername != "" && hetznerRobotPassword != "" {
		robotClient = hetzner.NewClient(hetznerRobotUsername, hetznerRobotPassword)
	}
	hetznerClient := &HetznerClient{
		cloudClient: client,
		robotClient: robotClient,
	}
	logrus.Debug("newCloud - returning cloud provider")
	return &cloud{
		client:    hetznerClient,
		zones:     newZones(hetznerClient, nodeName),
		instances: newInstances(hetznerClient),
	}, nil
}

func (c *cloud) Initialize(clientBuilder controller.ControllerClientBuilder) {}

func (c *cloud) Instances() (cloudprovider.Instances, bool) {
	logrus.Debug("returning instances")
	return c.instances, true
}

func (c *cloud) Zones() (cloudprovider.Zones, bool) {
	logrus.Debug("returning zones")
	return c.zones, true
}

func (c *cloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	logrus.Debug("LoadBalancer not implemented")
	return nil, false
}

func (c *cloud) Clusters() (cloudprovider.Clusters, bool) {
	logrus.Debug("Cluster not implemented")
	return nil, false
}

func (c *cloud) Routes() (cloudprovider.Routes, bool) {
	logrus.Debug("Routes not implemented")
	return nil, false
}

func (c *cloud) ProviderName() string {
	logrus.Debug("Returning providerName")
	return providerName
}

func (c *cloud) ScrubDNS(nameservers, searches []string) (nsOut, srchOut []string) {
	logrus.Debug("ScrubDNS not implemented")
	return nil, nil
}

func (c *cloud) HasClusterID() bool {
	logrus.Debug("setting hasClusterID to false")
	return false
}

func init() {
	ll := os.Getenv(logLevel)
	logrus.SetLevel(logrus.InfoLevel)
	level, err := logrus.ParseLevel(ll)
	if err == nil {
		logrus.SetLevel(level)
	}
	logrus.Info("registering our own sofware as cloud provider")
	cloudprovider.RegisterCloudProvider(providerName, func(config io.Reader) (cloudprovider.Interface, error) {
		return newCloud(config)
	})
}
