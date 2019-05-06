package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/vmware/govmomi/vim25/mo"

	"github.com/2spmohanty/gosphere/operation"
	"github.com/vmware/govmomi/find"
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

	finder := find.NewFinder(vcenter.Client.Client, true)

	pc := property.DefaultCollector(vcenter.Client.Client)

	datacenter, err := finder.Datacenter(ctx, "vcqaDC")

	finder.SetDatacenter(datacenter)

	cluster, err := finder.ClusterComputeResource(ctx, "als")
	if err != nil {
		fmt.Println("Error while Retrieving Cluster Object is ", err)
	}
	var clst mo.ClusterComputeResource
	err = pc.RetrieveOne(ctx, cluster.Reference(), nil, &clst)
	if err != nil {
		fmt.Println("Error while Retrieving Cluster Compute Resourceis ", err)
	}

	host, err := finder.HostSystem(ctx, "10.161.34.101")
	if err != nil {
		fmt.Println("Error while Retrieving Host Object is ", err)
	}
	var hst mo.HostSystem
	err = pc.RetrieveOne(ctx, host.Reference(), nil, &hst)
	if err != nil {
		fmt.Println("Error while Retrieving Host System is ", err)
	}

	vm, err := finder.VirtualMachine(ctx, "standalone-bc66ae0ab-esx.3-vm.0")
	if err != nil {
		fmt.Println("Error while Retrieving VM Object is  ", err)
	}

	ds, err := finder.Datastore(ctx, "sharedVmfs-0")
	if err != nil {
		fmt.Println("Error while Retrieving Datastote Object is ", err)
	}
	var dst mo.Datastore
	err = pc.RetrieveOne(ctx, ds.Reference(), nil, &dst)
	if err != nil {
		fmt.Println("Error while Retrieving Datastore system is ", err)
	}

	fmt.Println("Starting Clone ....")

	vmops := operation.VMOperation{Context: ctx, Vcenter: vcenter}

	vmMor, rs := vmops.CloneVM("Test-VM-Last", true, &hst, vm, &clst, datacenter, &dst)

	var vmt mo.VirtualMachine
	err = pc.RetrieveOne(ctx, vmMor.Reference(), nil, &vmt)
	fmt.Println("The New VM is ", vmt.Name)
	fmt.Println("The state of cloning is ", rs)

}
