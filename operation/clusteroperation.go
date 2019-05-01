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
	"strings"

	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

//GetAllClusterHosts Returns all the Hosts in a given Cluster Object.
func (vcenter *VCenter) GetAllClusterHosts(ctx context.Context, clst mo.ClusterComputeResource, powerstate string) ([]mo.HostSystem, error) {

	pc := property.DefaultCollector(vcenter.Client.Client)

	hosts := clst.Host
	var hostref []types.ManagedObjectReference
	for _, host := range hosts {
		hostref = append(hostref, host.Reference())

	}

	var hst []mo.HostSystem
	err := pc.Retrieve(ctx, hostref, nil, &hst)
	if err != nil {
		exit(err)
	}

	var hstpower []mo.HostSystem

	for _, hs := range hst {
		if string(hs.Runtime.PowerState) == "poweredOn" && strings.Contains(powerstate, "On") {
			hstpower := append(hstpower, hs)
			return hstpower, nil
		} else if string(hs.Runtime.PowerState) == "poweredOff" && strings.Contains(powerstate, "Off") {
			hstpower := append(hstpower, hs)
			return hstpower, nil
		}

	}

	return hst, nil

}

//GetClusterData is used by getvcinventory in example. It Returns Clusters and its child object.
//Do not change this sample.
func GetClusterData(ctx context.Context, pc *property.Collector, clst mo.ClusterComputeResource, clus chan ClusterStruct, getvm bool) {

	//var cg sync.WaitGroup

	Hostarray := []HostStruct{}
	clustername := clst.Name
	hosts := clst.Host

	//cg.Add(len(hosts))
	//defer cg.Done()

	var hostref []types.ManagedObjectReference
	for _, host := range hosts {
		hostref = append(hostref, host.Reference())

	}

	var hst []mo.HostSystem
	err := pc.Retrieve(ctx, hostref, nil, &hst)
	if err != nil {
		exit(err)
	}

	for _, hs := range hst {
		vmchan := make(chan []VMStruct)
		if getvm {
			go GetVMData(ctx, pc, hs, vmchan)
			vmdata := <-vmchan
			Hostarray = append(Hostarray, HostStruct{hs.Name, string(hs.Runtime.PowerState), vmdata})
		} else {
			Hostarray = append(Hostarray, HostStruct{hs.Name, string(hs.Runtime.PowerState), nil})
		}

	}

	Cls := ClusterStruct{clustername, Hostarray}

	clus <- Cls

	//cg.Wait()

}
