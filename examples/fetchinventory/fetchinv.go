package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/vmware/govmomi/vim25/mo"

	"github.com/2spmohanty/gosphere/operation"
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

	//Get Datacenter Operation Level object
	dcops := operation.DatacenterOperation{Context: ctx, Vcenter: vcenter}

	//Get Cluster Operation Level object
	clops := operation.ClusterOperation{Context: ctx, Vcenter: vcenter}

	datacenters, err := vcenter.GetAllDatacenter(ctx)
	if err != nil {
		fmt.Printf("Datacenters errors: %s", err)
		return
	}

	for _, dc := range datacenters {

		dcName := dc.Name()

		fmt.Printf(" Datacenter %s\n", dcName)

		standalonehosts := dcops.GetStandAloneHosts(dc)

		if standalonehosts != nil {
			fmt.Printf("Standalone Hosts on Datacenter %s\n", dcName)
			for _, hostmor := range standalonehosts {
				fmt.Println(hostmor.Name)
			}
		}

		var cls []mo.ClusterComputeResource

		cls, _ = dcops.GetAllCluster(dc)

		if cls != nil {

			for _, clsref := range cls {

				fmt.Printf("Datcenter Clusters ***** %s ******\n", clsref.Name)

				var hosts []mo.HostSystem
				hosts, _ = clops.GetAllClusterHosts(clsref, "")

				if hosts != nil {
					fmt.Printf("Cluster Hosts")
					for _, hostref := range hosts {
						fmt.Printf("**** %s ****\n", hostref.Name)

					}
				}

			}

		}

	}
}
