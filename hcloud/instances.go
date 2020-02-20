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
	"errors"
	"strconv"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

type instances struct {
	client *HetznerClient
}

func newInstances(client *HetznerClient) *instances {
	logrus.Debug("Returning new instances")
	return &instances{client}
}

func (i *instances) NodeAddressesByProviderID(providerID string) ([]v1.NodeAddress, error) {
	logrus.Debug("NodeAddressesByProvider")
	logrus.Debug("Searching for server")
	server, err := i.client.GetServerByProviderID(providerID)
	if err != nil {
		return nil, err
	}
	logrus.Debug("parsing server to node addresses")
	return nodeAddresses(server.Name, server.Ipv4), nil
	// return nil, cloudprovider.InstanceNotFound
}

func (i *instances) NodeAddresses(nodeName types.NodeName) ([]v1.NodeAddress, error) {
	logrus.Debug("NodeAddresses")
	logrus.Debug("Getting server by name")
	server, err := i.client.GetServerByName(string(nodeName))
	if err != nil {
		return nil, err
	}
	logrus.Debug("Return nodeAddresses")
	return nodeAddresses(server.Name, server.Ipv4), nil
}

func (i *instances) ExternalID(nodeName types.NodeName) (string, error) {
	logrus.Debug("Return externalID")
	return i.InstanceID(nodeName)
}

func (i *instances) InstanceID(nodeName types.NodeName) (string, error) {
	logrus.Debug("InstanceID")
	logrus.Debug("Searching server by name")
	server, err := i.client.GetServerByName(string(nodeName))
	if err != nil {
		return "", err
	}
	return strconv.Itoa(server.ID), nil
}

func (i *instances) InstanceType(nodeName types.NodeName) (string, error) {
	logrus.Debug("InstanceType")
	logrus.Debug("Searching server by name")
	server, err := i.client.GetServerByName(string(nodeName))
	if err != nil {
		return "", err
	}
	logrus.Debug("Returning instanceType")
	return server.InstanceType, nil
}

func (i *instances) InstanceTypeByProviderID(providerID string) (string, error) {
	logrus.Debug("Getting ingress type by provider ID")
	logrus.Debug("Getting server")
	server, err := i.client.GetServerByProviderID(providerID)
	if err != nil {
		return "", err
	}
	logrus.Debug("Retruning instance type")
	return server.InstanceType, nil
}

func (i *instances) AddSSHKeyToAllInstances(user string, keyData []byte) error {
	logrus.Debug("AddSSHKeyToAllInstances --> not implemented")
	return errors.New("not implemented")
}

func (i *instances) CurrentNodeName(hostname string) (types.NodeName, error) {
	logrus.Debug("Returning current node name")
	return types.NodeName(hostname), nil
}

func (i instances) InstanceExistsByProviderID(providerID string) (exists bool, err error) {
	logrus.Debug("Checking existence by provider ID")
	var server *Server
	logrus.Debug("Getting server")
	server, err = i.client.GetServerByProviderID(providerID)
	if err != nil {
		return
	}
	logrus.Debug("Checking existence")
	exists = server != nil
	return
}

func nodeAddresses(name, ipv4 string) []v1.NodeAddress {
	logrus.Debug("Generating node address")
	var addresses []v1.NodeAddress
	addresses = append(
		addresses,
		v1.NodeAddress{Type: v1.NodeHostName, Address: name},
		v1.NodeAddress{Type: v1.NodeExternalIP, Address: ipv4},
	)
	return addresses
}
