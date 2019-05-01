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

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

//GetAllCluster Returns all the cluster object under a data center
func (vcenter *VCenter) GetAllCluster(ctx context.Context, datacenter *object.Datacenter) ([]mo.ClusterComputeResource, error) {

	finder := find.NewFinder(vcenter.Client.Client, true)
	finder.SetDatacenter(datacenter)

	pc := property.DefaultCollector(vcenter.Client.Client)

	clusters, err := finder.ClusterComputeResourceList(ctx, "*")
	if err != nil {
		exit(err)
	}

	var clusterref []types.ManagedObjectReference
	for _, cluster := range clusters {
		clusterref = append(clusterref, cluster.Reference())

	}

	// Retrieve All property for all Clusters
	var clst []mo.ClusterComputeResource
	err = pc.Retrieve(ctx, clusterref, nil, &clst)
	if err != nil {
		exit(err)
	}

	return clst, nil
}

//GetCluster Returns cluster object under a data center when a cluster name is passed to it.
func (vcenter *VCenter) GetCluster(ctx context.Context, clustername string) (mo.ClusterComputeResource, error) {

	finder := find.NewFinder(vcenter.Client.Client, true)

	pc := property.DefaultCollector(vcenter.Client.Client)

	cluster, err := finder.ClusterComputeResource(ctx, clustername)
	if err != nil {
		exit(err)
	}

	var clusterref types.ManagedObjectReference

	clusterref = cluster.Reference()

	// Retrieve All property for the Clusters
	var clst mo.ClusterComputeResource

	err = pc.RetrieveOne(ctx, clusterref, nil, &clst)
	if err != nil {
		exit(err)
	}

	return clst, nil
}

//GetStandAloneHosts Returns all the Standalone Hosts Objects in a Cluster
func (vcenter *VCenter) GetStandAloneHosts(ctx context.Context, datacenter *object.Datacenter) []mo.HostSystem {

	pc := property.DefaultCollector(vcenter.Client.Client)

	dcfolder, err := datacenter.Folders(ctx)
	if err != nil {
		exit(err)
	}

	hostfolder := dcfolder.HostFolder
	standalonehosts, _ := WalkFolder(ctx, hostfolder)

	var hst []mo.ComputeResource
	err = pc.Retrieve(ctx, standalonehosts, nil, &hst)
	if err != nil {
		exit(err)
	}

	var hostref [][]types.ManagedObjectReference
	for _, hst := range hst {
		hostref = append(hostref, hst.Host)
	}

	var hs []mo.HostSystem
	for _, hosts := range hostref {
		err := pc.Retrieve(ctx, hosts, nil, &hs)
		if err != nil {
			exit(err)
		}
	}

	return hs

}

//WalkFolder is used to Walk the Folder of a Inventory Object
func WalkFolder(ctx context.Context, f *object.Folder) ([]types.ManagedObjectReference, []types.ManagedObjectReference) {
	var standalonehosts []types.ManagedObjectReference
	var clusters []types.ManagedObjectReference
	childEntities, err := f.Children(ctx)
	if err != nil {
		exit(err)
	}
	for _, childEntity := range childEntities {

		stdhost, cluster := WalkManagedEntity(childEntity.Reference())
		if stdhost {
			standalonehosts = append(standalonehosts, childEntity.Reference())
		} else if cluster {
			clusters = append(clusters, childEntity.Reference())
		}

	}

	return standalonehosts, clusters

}

//WalkManagedEntity is used to Walk through the child entity of a Managed ENtity
func WalkManagedEntity(childEntity types.ManagedObjectReference) (bool, bool) {

	var standalonehost bool
	var clusterhost bool
	if childEntity.Type == "ComputeResource" {
		standalonehost = true
	} else if childEntity.Type == "ClusterComputeResource" {
		clusterhost = true
	}

	return standalonehost, clusterhost

}

//CreateCluster creates a Cluster under a specified datacenter object.
func (vcenter *VCenter) CreateCluster(ctx context.Context, datacenter *object.Datacenter, clustername string) (*object.ClusterComputeResource, error) {
	dcfolder, err := datacenter.Folders(ctx)
	if err != nil {
		exit(err)
	}

	hostfolder := dcfolder.HostFolder

	var clsConfigSpec types.ClusterConfigSpecEx
	clustermor, err := hostfolder.CreateCluster(ctx, clustername, clsConfigSpec)

	if err != nil {
		exit(err)
	}

	return clustermor, nil
}
