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

	"github.com/gosphere/operation"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
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

//ByName Displays the VM
type ByName []mo.VirtualMachine

func (n ByName) Len() int           { return len(n) }
func (n ByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByName) Less(i, j int) bool { return n[i].Name < n[j].Name }

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

	// Create DC Map to STore Value

	//DatacenterValue := make(map[string][]string)
	//ClusterValue := make(map[string][]string)
	//HostValue := make(map[string][]string)

	//finder := find.NewFinder(vcenter.Client.Client, true)

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

		standalonehosts := vcenter.GetStandAloneHosts(ctx, dc)

		for _, hostmor := range standalonehosts {
			fmt.Println(hostmor.Name)
		}

		// Get Cluster

		clst1, _ := vcenter.GetAllCluster(ctx, dc)
		if err != nil {
			log.Println(err)
		}

		for _, cls1 := range clst1 {

			cluschan := make(chan operation.ClusterStruct)
			go operation.GetClusterData(ctx, pc, cls1, cluschan, false)
			hostresult := <-cluschan
			jsonData, err := json.Marshal(hostresult)
			if err != nil {
				log.Println(err)
			}
			fmt.Println(string(jsonData))

		}

	}

}
