/*
Author : Smruti P Mohanty

Copyright (c) 2017 VMware, Inc. All Rights Reserved.

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

package operation

import (
	"context"
	"fmt"
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
)

// VCenter represents a vCenter server
type VCenter struct {
	Hostname string
	Username string
	Password string
	Client   *govmomi.Client
}

// NewVCenter returns a new VCenter instance
func NewVCenter(hostname string, username string, password string) *VCenter {
	return &VCenter{
		Hostname: hostname,
		Username: username,
		Password: password,
	}
}

//Connect object is used to connect vCenter
func (vcenter *VCenter) Connect(ctx context.Context) error {
	fmt.Printf("Connecting to vcenter: %s\n", vcenter.Hostname)
	if len(vcenter.Hostname) == 0 {
		return fmt.Errorf("No vCenter host defined")
	}
	u, err := url.Parse("https://" + vcenter.Username + ":" + vcenter.Password + "@" + vcenter.Hostname + "/sdk")
	if err != nil {
		fmt.Printf("Error with URL: %s\n", err)
		return err
	}
	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		fmt.Printf("Could not connect to vcenter: %s. Error: %v\n", vcenter.Hostname, err)
		return err
	}
	vcenter.Client = client
	return nil
}

//GetAllDatacenter Returns all the datacenter in a Given vCenter.
//The vCenter reciever must be connected using vcenter.Connect(ctx) where ctx is the context of type context.Context
func (vcenter *VCenter) GetAllDatacenter(ctx context.Context) ([]*object.Datacenter, error) {

	finder := find.NewFinder(vcenter.Client.Client, true)

	datacenters, err := finder.DatacenterList(ctx, "*")
	if err != nil {
		fmt.Printf("Datacenters errors: %s", err)
		return nil, err
	}
	return datacenters, nil

}

//GetDatacenter Returns datacenter object when a datacenter name in a Given vCenter is passed to it.
//The vCenter reciever must be connected using vcenter.Connect(ctx) where ctx is the context of type context.Context
func (vcenter *VCenter) GetDatacenter(ctx context.Context, dcname string) (*object.Datacenter, error) {

	finder := find.NewFinder(vcenter.Client.Client, true)

	datacenter, err := finder.Datacenter(ctx, dcname)
	if err != nil {
		fmt.Printf("Datacenters errors: %s", err)
		return nil, err
	}
	return datacenter, nil

}
