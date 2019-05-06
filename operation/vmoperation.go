package operation

import (
	"context"
	"fmt"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

//VMOperation is the Reciever Object For all VM Operation
type VMOperation struct {
	Context context.Context
	Vcenter *VCenter
}

//GetVMData is used by getvcinventory in example. It Returns VM Objects and its power state.
func GetVMData(ctx context.Context, pc *property.Collector, hst mo.HostSystem, vmchan chan []VMStruct) {

	//var wg sync.WaitGroup
	var vms []types.ManagedObjectReference
	vms = hst.Vm

	vmarray := []VMStruct{}

	//wg.Add(len(vms))
	//defer wg.Done()
	if vms != nil {
		var refs []types.ManagedObjectReference
		for _, vm := range vms {
			refs = append(refs, vm.Reference())
		}

		var vmt []mo.VirtualMachine
		err := pc.Retrieve(ctx, refs, nil, &vmt)
		if err != nil {
			exit(err)
		}

		for _, vm := range vmt {
			vmarray = append(vmarray, VMStruct{vm.Name, string(vm.Summary.Runtime.PowerState)})

		}

	}

	vmchan <- vmarray

	//

}

//CloneVM Clones a VM
func (vmops *VMOperation) CloneVM(newVMName string, poweron bool, host *mo.HostSystem, template *object.VirtualMachine, cluster *mo.ClusterComputeResource, datacenter *object.Datacenter, datastore *mo.Datastore) (*object.VirtualMachine, types.TaskInfoState) {

	ctx := vmops.Context
	c := vmops.Vcenter.Client.Client

	//Get the Resourcepool
	resourcepool := cluster.ResourcePool

	dcfolder, err := datacenter.Folders(ctx)
	if err != nil {
		exit(err)
	}
	vmfolder := dcfolder.VmFolder

	relocationSpec := types.VirtualMachineRelocateSpec{

		Pool: resourcepool,
	}
	hostref := host.Reference()
	relocationSpec.Host = &hostref

	datastoreRef := datastore.Reference()
	relocationSpec.Datastore = &datastoreRef

	cloneSpec := &types.VirtualMachineCloneSpec{
		PowerOn:  poweron,
		Template: false,
	}

	cloneSpec.Location = relocationSpec

	task, _ := template.Clone(ctx, vmfolder, newVMName, *cloneSpec)

	info, err := task.WaitForResult(ctx, nil)

	if err != nil {
		fmt.Println("Task failed  with error ", err)
	}

	fmt.Printf("%s Cloning completed with %s.\n", newVMName, info.State)

	return object.NewVirtualMachine(c, info.Result.(types.ManagedObjectReference)), info.State

}
