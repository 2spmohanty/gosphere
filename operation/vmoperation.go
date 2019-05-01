package operation

import (
	"context"

	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

//GetVMData is used by getvcinventory in example. It Returns VM Objects and its power state.
func GetVMData(ctx context.Context, pc *property.Collector, hst mo.HostSystem, vmchan chan []VMStruct) {

	//var wg sync.WaitGroup

	vms := hst.Vm

	//wg.Add(len(vms))
	//defer wg.Done()

	var refs []types.ManagedObjectReference
	for _, vm := range vms {
		refs = append(refs, vm.Reference())
	}

	var vmt []mo.VirtualMachine
	err := pc.Retrieve(ctx, refs, nil, &vmt)
	if err != nil {
		exit(err)
	}

	vmarray := []VMStruct{}
	for _, vm := range vmt {
		vmarray = append(vmarray, VMStruct{vm.Name, string(vm.Summary.Runtime.PowerState)})

	}

	vmchan <- vmarray

	//

}
