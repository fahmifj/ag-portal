package service

import (
	"github.com/fahmifj/ag-portal/service/vm"
)

type Service struct {
	vm.VMControl
}

func NewServices() *Service {
	return &Service{
		VMControl: vm.NewVMInstance(),
	}
}
