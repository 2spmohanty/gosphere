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
go build getvcinv main.go
usage: ./getvcinv -vc vcip -user username -pass

*/

package main

import (
	"context"
	"encoding/json"
	"log"

	"flag"
	"fmt"

	"os"
	"runtime"
	"time"

	"github.com/2spmohanty/gosphere/operation"
	"github.com/vmware/govmomi/property"
)

type InventoryStruct struct {
	Vcenter     string
	Datacenters []DatacenterStruct
}

type DatacenterStruct struct {
	Dcname          string
	Clusters        []ClusterStruct
	StandaloneHosts []HostStruct
}

type ClusterStruct struct {
	Clustername string
	Hosts       []HostStruct
}

type HostStruct struct {
	Hostname   string
	Connection string
	Vms        []VmStruct
}

type VmStruct struct {
	Vmname       string
	Vmconnection string
}

func exit(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(1)
}

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

	pc := property.DefaultCollector(vcenter.Client.Client)

	datacenters, err := vcenter.GetAllDatacenter(ctx)
	if err != nil {
		fmt.Printf("Datacenters errors: %s", err)
		return
	}

	for _, dc := range datacenters {

		dcName := dc.Name()

		fmt.Printf("Working on Datacenter %s\n", dcName)

		//Get StandALone Hosts

		fmt.Printf("Getting Standalone Hosts for %s\n", dcName)

		standalonehosts := dcops.GetStandAloneHosts(dc)

		if standalonehosts != nil {
			for _, hostmor := range standalonehosts {
				fmt.Println(hostmor.Name)
			}
		}

		clst1, _ := dcops.GetAllCluster(dc)
		//fmt.Println("Cls ", clst1)
		if err != nil {
			log.Println(err)
		}

		if clst1 != nil {
			for _, cls1 := range clst1 {

				cluschan := make(chan operation.ClusterStruct)
				go operation.GetClusterData(ctx, pc, cls1, cluschan, true)
				hostresult := <-cluschan
				jsonData, err := json.Marshal(hostresult)
				if err != nil {
					log.Println(err)
				}
				fmt.Println(string(jsonData))

			}
		}

	}
	pc.Destroy(ctx)

}
