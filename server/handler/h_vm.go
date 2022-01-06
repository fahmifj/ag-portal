package handler

import (
	"net/http"

	"github.com/fahmifj/ag-portal/logger"
	"github.com/fahmifj/ag-portal/util"
	"github.com/gorilla/mux"
)

const (
	start = "start"
	stop  = "stop"
)

type response struct {
	Message string
}

var c = make(chan error)

func (r Router) VMControlHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	action := req.URL.Query().Get("action")
	vmName := getVMname(req)

	switch action {
	case start:
		logger.Log.Info("Starting VM", "VM", vmName)
		go r.service.StartVM(vmName, c)
	case stop:
		logger.Log.Info("Stopping VM", "VM", vmName)
		go r.service.StopVM(vmName, c)
	default:
		util.ToJSON(response{Message: "Invalid action"}, resp)
		return
	}

	util.ToJSON(response{Message: "Ok"}, resp)

	go logPrinter(action, vmName)
}

func (r Router) FetchVM(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	d := r.service.ListVM(c)
	if len(d) == 0 {
		util.ToJSON(response{Message: "Failed"}, resp)
		logger.Log.Info("Fetching VM failed")
		return
	}

	util.ToJSON(d, resp)

}

func logPrinter(action string, vmName string) {
	err := <-c
	if err != nil {
		logger.Log.Error("Operation failed", "action", action, "VM", vmName, "reason", err.Error())
		return
	}
	logger.Log.Info(err.Error(), "VM", vmName)

}

func getVMname(req *http.Request) string {
	return mux.Vars(req)["vm"]
}
