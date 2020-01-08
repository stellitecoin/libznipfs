// package name: libznipfs
package main

import "C"
import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"github.com/scala-network/libznipfs/src/ipfs"
)

var ipfsNode *ipfs.IPFS

// Result holds the seedlist and any error that occurred in the process
// for the daemon to use
type Result struct {
	// Status for the result
	Status string
	// Message to be displayed
	Message string
	// The seedlist
	Seedlist []string
}

func main() {
	// We need the main function to make possible
	// CGO compiler to compile the package as C library
}

/**
 * libznipfs implements the C-style library for fetching information
 * from ZeroNet and IPFS.
 * Here we only have 3 exported functions that can be called from C
 */

//export IPFSStartNode
// IPFSStartNode starts the IPFS node and initializes ZeroNet
func IPFSStartNode(dataPath *C.char) *C.char {
	// result is marshalled to JSON before being returned to the daemon
	result := Result{
		Status:  "ok",
		Message: fmt.Sprintf("IPFS node started on port 5001"),
	}
	var err error
	basePath := C.GoString(dataPath)

	ipfsNode, err = ipfs.New(filepath.Join(basePath, "ipfs"))
	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("Unable to create IPFS node: %s\n", err)
		return toCJSONString(result)
	}

	err = ipfsNode.Start()
	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("Unable to start IPFS node: %s\n", err)
	}

	return toCJSONString(result)
}

/*
//export ZNIPFSGetSeedList
// ZNIPFSGetSeedList retrieves the seedlist using ZeroNet and IPFS and returns
// it as JSON to the daemon. We use a named return here to ensure any
// lower level panic's (from 3rd party libs) are captured back to the daemon
func ZNIPFSGetSeedList(zeroNetAddress *C.char) (resultJSON *C.char) {
	// This defer/recover block captures any lower lever panics that might
	// occur in the 3rd party IPFS and ZeroNet libraries. It prevents the
	// daemon from crashing should such an error occur.
	defer func() {
		if r := recover(); r != nil {
			resultJSON = toCJSONString(Result{
				Status:  "err",
				Message: fmt.Sprintf("Unable to fetch seedlist from IPFS and ZeroNet: %s", r),
			})
			return
		}
	}()

	// Returns the address list from the given ZeroNet address
	result := Result{
		Status:  "ok",
		Message: fmt.Sprintf("Seedlist retrieved from ZeroNet and IPFS"),
	}

	ipfsHash := "Qmaisz6NMhDB51cCvNWa1GMS7LU1pAxdF4Ld6Ft9kZEP2a"

	data, err := ipfsNode.Get(ipfsHash)
	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("Unable fetch data from IPFS node: %s\n", err)
		resultJSON = toCJSONString(result)
		return
	}

	// data contains a JSON array with the seed list
	err = json.Unmarshal(data, &result.Seedlist)
	if err != nil {
		result.Status = "err"
		result.Message = fmt.Sprintf("Invalid seedlist format: %s\n", err)
	}

	// If the seedlist was in the correct format is has been stored in
	// result.Seedlist and can be returned without reassigning
	resultJSON = toCJSONString(result)
	return
}
*/
/*
//export IPFSStopNode
func IPFSStopNode() {
	// Stop the ZN/IPFS node
	ipfsNode.Stop()
}
*/
/* Should implement this function into the ipfs.go file instead */
/*
//export resolve
func resolve() (resultJSON *C.char) {
	defer func() {
		if r := recover(); r != nil {
			resultJSON = toCJSONString(Result{
				Status:  "err",
				Message: fmt.Sprintf("Resolution failed HORRIBLY. %s", r),
			})
			return
		}
	}()

	sh := shell.NewShell("localhost:5001")
	result, err := sh.Resolve("QmNW7Db89EGQZmf6cDBrFANQWL2XCDCMddneQjqXV6ssUC")
	if err != nil {
		resultJSON = toCJSONString(Result{
			Status:  "err",
			Message: fmt.Sprintf("Resolution failed HORRIBLY. %s", err),
		})
		return
	}
	resultJSON = toCJSONString(Result{
		Status:  "OK",
		Message: fmt.Sprintf(result),
	})
	return
}
*/
// toCJSONString marshals the error result into JSON for the daemon to
// understand and returns it in the required C format
func toCJSONString(result Result) *C.char {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		panic(fmt.Errorf("Fatal error converting result: %s", err))
	}
	return C.CString(string(resultJSON))
}
