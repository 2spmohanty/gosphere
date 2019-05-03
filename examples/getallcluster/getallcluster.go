package main

import (
	"context"
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/2spmohanty/gosphere/operation"
)

func main() {

	now := time.Now()
	defer func() {
		fmt.Println("Inventory Fetch Took ", time.Now().Sub(now))
	}()

	runtime.GOMAXPROCS(4)

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

	dcops := operation.DatacenterOperation{Context: ctx, Vcenter: vcenter}

	datacenters, err := vcenter.GetAllDatacenter(ctx)
	if err != nil {
		fmt.Printf("Datacenters errors: %s", err)
		return
	}

	for _, dc := range datacenters {
		cls, err := dcops.GetAllCluster(dc)
		if err != nil {
			fmt.Printf("Datacenters errors: %s", err)
			return
		}

		for _, clsref := range cls {
			fmt.Println(clsref.Name)
		}

	}

}
