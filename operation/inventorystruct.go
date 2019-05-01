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
	"fmt"
	"os"
)

func exit(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(1)
}

//ClusterStruct is used for storing Cluster Name and the Arrays of Host Name within the cluster.
type ClusterStruct struct {
	Cluster string
	Hosts   []HostStruct
}

//HostStruct is used to store Hostname, its powerstate and the VMs that it contains.
type HostStruct struct {
	Hostname   string
	Connection string
	Vms        []VMStruct
}

//VMStruct contains the VM name and Its power State.
type VMStruct struct {
	VMName     string
	PowerState string
}
