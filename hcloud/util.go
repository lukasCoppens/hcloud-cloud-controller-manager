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
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/appscode/go-hetzner"
	"github.com/hetznercloud/hcloud-go/hcloud"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

func getCloudServerByName(c *hcloud.Client, name string) (server *hcloud.Server, err error) {
	server, _, err = c.Server.GetByName(context.Background(), name)
	if err != nil {
		return
	}
	if server == nil {
		err = cloudprovider.InstanceNotFound
		return
	}
	return
}

func getRobotServerByName(c *hetzner.Client, name string) (*hetzner.ServerSummary, error) {
	summary, _, err := c.Server.ListServers()
	if err != nil {
		return nil, err
	}
	for _, entry := range summary {
		if entry.ServerName == name {
			return entry, nil
		}
	}
	// We did not find the server at this stage
	return nil, cloudprovider.InstanceNotFound
}

func getCloudServerByID(c *hcloud.Client, id int) (server *hcloud.Server, err error) {
	server, _, err = c.Server.GetByID(context.Background(), id)
	if err != nil {
		return
	}
	if server == nil {
		err = cloudprovider.InstanceNotFound
		return
	}
	return
}

func getRobotServerByID(c *hetzner.Client, id int) (*hetzner.ServerSummary, error) {
	summary, _, err := c.Server.ListServers()
	if err != nil {
		return nil, err
	}
	for _, entry := range summary {
		if entry.ServerNumber == id {
			return entry, nil
		}
	}
	// We did not find the server at this stage
	return nil, cloudprovider.InstanceNotFound
}

func convertFailureDomainToRegion(failuredomain string) string {
	splitted := strings.Split(failuredomain, "-")
	if len(splitted) > 0 {
		return splitted[0]
	}
	return ""
}

func providerIDToServerID(providerID string) (id int, err error) {
	providerPrefix := providerName + "://"
	if !strings.HasPrefix(providerID, providerPrefix) {
		err = fmt.Errorf("providerID should start with %s://: %s", providerName, providerID)
		return
	}
	idString := strings.TrimPrefix(providerID, providerPrefix)
	if idString == "" {
		err = fmt.Errorf("missing server id in providerID: %s", providerID)
		return
	}
	id, err = strconv.Atoi(idString)
	return
}

func convertCloudServerToServer(server *hcloud.Server) *Server {
	return &Server{
		Name:          server.Name,
		Ipv4:          server.PublicNet.IPv4.IP.String(),
		Region:        server.Datacenter.Location.Name,
		Failuredomain: server.Datacenter.Name,
		InstanceType:  server.ServerType.Name,
		ID:            server.ID,
	}
}

func convertRobotServerToServer(server *hetzner.ServerSummary) *Server {
	lowerCaseDc := strings.ToLower(server.Dc)
	return &Server{
		Name:          server.ServerName,
		Ipv4:          server.ServerIP,
		Failuredomain: lowerCaseDc,
		Region:        convertFailureDomainToRegion(lowerCaseDc),
		InstanceType:  strings.ToLower(server.Product),
		ID:            server.ServerNumber,
	}
}
