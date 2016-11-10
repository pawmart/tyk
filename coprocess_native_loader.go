// +build coprocess
// +build native

package main

/*
#cgo python CFLAGS: -DENABLE_PYTHON

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include <dlfcn.h>
typedef void (*sym)(void*, void*);
sym processRequestSym;

#include "coprocess/sds/sds.h"

#include "coprocess/api.h"

static struct CoProcessMessage* DispatchHook(struct CoProcessMessage* object) {
	struct CoProcessMessage* outputObject = malloc(sizeof *outputObject);
	// char *output = malloc(object->length);
	processRequestSym(object, outputObject);
	// strncpy(output, object->p_data, object->length);
	// outputObject->p_data = (void*)output;
	// outputObject->length = object->length;
	return outputObject;
}

static void LoadMiddleware() {
	void* lib;
	lib = dlopen("/Users/matias/dev/tyk-native-plugin/tyk_middleware/middleware.so", RTLD_NOW);
	if(lib == NULL) {
		printf("Couldn't load library!\n");
		return;
	}
	processRequestSym = dlsym(lib, "ProcessRequest");
	if(processRequestSym == NULL){
		printf("Couldn't load symbol!\n");
		return;
	}
}

static void DispatchEvent(char* event_json) {
}

*/
import "C"

import (
	// "errors"
	"os"
	"path"
	// "strings"
	"unsafe"

	"github.com/Sirupsen/logrus"
	"github.com/TykTechnologies/tyk/coprocess"
	"github.com/TykTechnologies/tykcommon"
)

// CoProcessName declares the driver name.
const CoProcessName string = "python"

// MessageType sets the default message type.
var MessageType = coprocess.ProtobufMessage

// PythonDispatcher implements a coprocess.Dispatcher
type NativeDispatcher struct {
	coprocess.Dispatcher
}

// Dispatch takes a CoProcessMessage and sends it to the CP.
func (d *NativeDispatcher) Dispatch(objectPtr unsafe.Pointer) unsafe.Pointer {
	var object *C.struct_CoProcessMessage
	object = (*C.struct_CoProcessMessage)(objectPtr)

	var newObjectPtr *C.struct_CoProcessMessage
	newObjectPtr = C.DispatchHook(object)

	return unsafe.Pointer(newObjectPtr)
	// return objectPtr
}

// DispatchEvent dispatches a Tyk event.
func (d *NativeDispatcher) DispatchEvent(eventJSON []byte) {
	var CEventJSON *C.char
	CEventJSON = C.CString(string(eventJSON))
	C.DispatchEvent(CEventJSON)
	C.free(unsafe.Pointer(CEventJSON))
	return
}

// Reload triggers a reload affecting CP middlewares and event handlers.
func (d *NativeDispatcher) Reload() {
	return
}

// HandleMiddlewareCache isn't used by Python.
func (d* NativeDispatcher) HandleMiddlewareCache(b *tykcommon.BundleManifest, basePath string) {
	// var CBundlePath *C.char
	// CBundlePath = C.CString(basePath)
	// C.Python_HandleMiddlewareCache(CBundlePath)
	return
}

// PythonInit initializes the Python interpreter.
func PythonInit() (err error) {
	/*
	result := C.Python_Init()
	if result == 0 {
		err = errors.New("Can't Py_Initialize()")
	}
	return err
	*/
	return err
}

// PythonLoadDispatcher creates reference to the dispatcher class.
func PythonLoadDispatcher() (err error) {
	/*
	result := C.Python_LoadDispatcher()
	if result == -1 {
		err = errors.New("Can't load dispatcher")
	}
	return err
	*/
	return err
}

// PythonNewDispatcher creates an instance of TykDispatcher.
func PythonNewDispatcher(middlewarePath string, eventHandlerPath string, bundlePaths []string) (dispatcher coprocess.Dispatcher, err error) {
	/*
	var CMiddlewarePath *C.char
	CMiddlewarePath = C.CString(middlewarePath)

	var CEventHandlerPath *C.char
	CEventHandlerPath = C.CString(eventHandlerPath)

	var CBundlePaths *C.char
	CBundlePaths = C.CString(strings.Join(bundlePaths, ":"))

	result := C.Python_NewDispatcher(CMiddlewarePath, CEventHandlerPath, CBundlePaths)

	if result == -1 {
		err = errors.New("Can't initialize a dispatcher")
	} else {
		dispatcher = &NativeDispatcher{}
	}

	C.free(unsafe.Pointer(CMiddlewarePath))
	C.free(unsafe.Pointer(CEventHandlerPath))
	*/
	return dispatcher, err
}

// PythonSetEnv sets PYTHONPATH, it's called before initializing the interpreter.
func PythonSetEnv(pythonPaths ...string) {
	/*
	var CPythonPath *C.char
	CPythonPath = C.CString(strings.Join(pythonPaths, ":"))
	C.Python_SetEnv(CPythonPath)

	C.free(unsafe.Pointer(CPythonPath))
	*/
}

// NewCoProcessDispatcher wraps all the actions needed for this CP.
func NewCoProcessDispatcher() (dispatcher coprocess.Dispatcher, err error) {

	workDir, _ := os.Getwd()

	dispatcherPath := path.Join(workDir, "coprocess/python")
	middlewarePath := path.Join(workDir, "middleware/python")
	eventHandlerPath := path.Join(workDir, "event_handlers")
	protoPath := path.Join(workDir, "coprocess/python/proto")

	paths := []string{dispatcherPath, middlewarePath, eventHandlerPath, protoPath}

	// Append bundle paths:
	bundlePaths := getBundlePaths()
	for _, v := range bundlePaths {
		paths = append(paths, v)
	}

	/*

	PythonSetEnv(paths...)

	PythonInit()
	PythonLoadDispatcher()

	dispatcher, err = PythonNewDispatcher(middlewarePath, eventHandlerPath, bundlePaths)

	C.PyEval_ReleaseLock()
	*/

	dispatcher = &NativeDispatcher{}
	C.LoadMiddleware()

	if err != nil {
		log.WithFields(logrus.Fields{
			"prefix": "coprocess",
		}).Error(err)
	}

	return dispatcher, err
}
