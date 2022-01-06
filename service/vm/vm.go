package vm

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/fahmifj/ag-portal/config"
)

type VMControl interface {
	// IsRunning() bool
	// IsStopped() bool
	GetVM(string, chan<- error) vmModel
	StartVM(string, chan<- error)
	StopVM(string, chan<- error)
	ListVM(chan<- error) []vmModel
}

type vmInstance struct{}

type vmModel struct {
	Name     string `json:"name,omitempty"`
	Status   string `json:"status,omitempty"`
	PublicIP string `json:"publicIp,omitempty"`
}

func NewVMInstance() VMControl {
	return &vmInstance{}
}

func (v vmInstance) ListVM(errChan chan<- error) []vmModel {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vm := vmModel{}
	vmList := make([]vmModel, 0)
	vmClient := compute.NewVirtualMachinesClient(config.SubscriptionID())
	vmClient.Authorizer, _ = getAuthorizer()

	for iter, err := vmClient.ListComplete(ctx, config.RG()); iter.NotDone(); err = iter.Next() {
		if err != nil {
			errChan <- err
			return nil
		}
		// fmt.Println(*iter.Value().Name)
		// fmt.Printf(iter.Value().Status)

		// fmt.Printf("%#v", *iter.Value().NetworkProfile.NetworkInterfaces)
		vm.Name = *iter.Value().Name
		// vm.Status = iter.Value().Status

		// for _, v := range *iter.Value().network {
		// 	fmt.Println(*v.ID)
		// }
		vmList = append(vmList, vm)
	}
	return vmList
	// // '123' -> listVM{{vm1}, {vm2}}
	// vmList := make([]vmModel, 0)
	// for x := range `01` {
	// 	vm := vmModel{
	// 		Name:     "VM-" + fmt.Sprint(x+1),
	// 		Status:   "stopped",
	// 		PublicIP: "127.0.0.1",
	// 	}
	// 	vmList = append(vmList, vm)
	// }
	// return vmList
}

func (v vmInstance) StartVM(vmName string, errChan chan<- error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	vmClient := getVMClient()
	result, err := vmClient.Start(ctx, config.RG(), vmName)
	if err != nil {
		errChan <- err
		return
	}
	err = result.WaitForCompletionRef(ctx, vmClient.Client)
	if err != nil {
		errChan <- err
		return
	}
}

// StopVM stops and deallocates VM instances
// https://docs.microsoft.com/en-us/answers/questions/574969/what39s-the-difference-between-deallocated-and-sto.html
func (v vmInstance) StopVM(vmName string, errChan chan<- error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	vmClient := getVMClient()
	result, err := vmClient.Deallocate(ctx, config.RG(), vmName, nil)
	if err != nil {
		errChan <- err
		return
	}
	err = result.WaitForCompletionRef(ctx, vmClient.Client)
	if err != nil {
		errChan <- err
		return
	}
}

// GetVM returns a virtual machine
// currently this called by JS fetch
// Error runtime
func (v vmInstance) GetVM(vmName string, errChan chan<- error) vmModel {
	// ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	// defer cancel()

	// i := vmModel{}
	// // TODO: GET VM status
	// vmClient := getVMClient()

	// result, err := vmClient.Get(ctx, config.RG(), config.VMName(), compute.InstanceViewTypesInstanceView)
	// if err != nil {
	// 	errChan <- err
	// 	return i
	// }

	// status := ""
	// for key, val := range *result.InstanceView.Statuses {
	// 	if *val.DisplayStatus != "" && key == 1 {
	// 		status = *val.DisplayStatus
	// 	}
	// }

	// publicIp, err := GetPublicIP()
	// if err != nil {
	// 	errChan <- err
	// 	return i
	// }
	// p := ""
	// if *publicIp.IPAddress != "" {
	// 	p = *publicIp.IPAddress
	// }

	// i.Name = config.VMName()
	// i.Status = status
	// i.PublicIP = p
	// return i
	i := vmModel{}
	i.Name = vmName
	i.Status = "running"
	i.PublicIP = "127.0.0.1"
	return i
}

func getVMClient() compute.VirtualMachinesClient {
	vmClient := compute.NewVirtualMachinesClient(config.SubscriptionID())
	vmClient.Authorizer, _ = getAuthorizer()
	return vmClient
}

func getAuthorizer() (autorest.Authorizer, error) {
	return auth.NewAuthorizerFromFile(azure.PublicCloud.ResourceManagerEndpoint)
}
