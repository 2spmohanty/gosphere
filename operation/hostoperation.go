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

	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

//HostOperation is the Reciever Object For all Host Operation
type HostOperation struct {
	Context context.Context
	Vcenter *VCenter
}

//GetAllVMs Returns All VMs Objects in a Host Object
func (hostops *HostOperation) GetAllVMs(hst mo.HostSystem) ([]mo.VirtualMachine, error) {

	ctx := hostops.Context
	client := hostops.Vcenter.Client.Client

	var vms []types.ManagedObjectReference
	vms = hst.Vm
	pc := property.DefaultCollector(client)

	var vmt []mo.VirtualMachine

	if vms != nil {
		var refs []types.ManagedObjectReference
		for _, vm := range vms {
			refs = append(refs, vm.Reference())
		}

		err := pc.Retrieve(ctx, refs, nil, &vmt)
		if err != nil {
			exit(err)
		}
	}

	return vmt, nil

}
