/*
Package main contains code to fetch Inventory Object of a Vcenter in JSON format.
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

/*
go build createinv createinv.go
usage: ./createinv -vc vcip -user username -pass

*/
package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/gosphere/operation"
	"github.com/vmware/govmomi/property"
)

func main() {
	vc := flag.String("vc", "10.161.50.3", "Enter vCenter IP/ FQDN")
	user := flag.String("user", "Administrator@vsphere.local", "vCenter User")
	pass := flag.String("pass", "Admin!23", "Enter vCenter pass")
	flag.Parse()

	vcenter := operation.NewVCenter(*vc, *user, *pass)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := vcenter.Connect(ctx)

	if err != nil {
		fmt.Printf("Failed to connect to vCenter: %s\n", err)
		return
	}

	fmt.Printf("Connected to vCenter: %s\n", *vc)

	pc := property.DefaultCollector(vcenter.Client.Client)

	//Create a Datacenter

	dcmor, err := vcenter.CreateDataCenter(ctx, "MyDatacenter")

	if err != nil {
		fmt.Printf("Failed to connect to vCenter: %s\n", err)
		return
	}

	fmt.Printf("Dataceter  created ")

	// Create a Cluster In the Datacenter

	clustermor, err := vcenter.CreateCluster(ctx, dcmor, "MyCluster")

	if err != nil {
		fmt.Printf("Failed to connect to vCenter: %s\n", err)
		return
	}

	var clusterref types.ManagedObjectReference
	clusterref = clustermor.Reference()
	var cls mo.ClusterComputeResource
	pc.RetrieveOne(ctx, clusterref, nil, &cls)
	fmt.Println("Cluster Created Sucessfully ", cls.Name)

}
