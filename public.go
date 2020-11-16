// License Header

package manager

import (
	"sync"
	"errors"
)

// NewManager will return a blank manager for use.
// Buffer size is the number of requests for the manager to hold onto until it starts blocking
//	requests.
func NewManager(name string, bufferSize int) *Manager {

	// Create a pointer to a new manager for clients to use.
	newManager := &Manager{

		Name: name,
		Requests: make(chan *Request, bufferSize),
		Functions: make(map[string]func(request interface{}) interface{}),
		StateLock: sync.Mutex{},

	}

	// Mutex management
	managersLock.Lock()
	defer managersLock.Unlock()

	// Add it to the managers map and return it
	managers[name] = newManager
	return newManager

}

// NewRequest will return a new request with the given Route and input Data
func NewRequest(route string, data interface{}) *Request {

	// Create a new request and return it with the give values
	return &Request{

		Route: route,
		Data: data,
		Response: make(chan Response, 1),

	}

}

// Send will create and send a request to a defined manager
func Send(managerName string, route string, data interface{}) error {

	// Get the manager
	manager, ok := getManager(managerName)

	// If the manager doesn't exist, respond with an error
	if !ok {
		return errors.New(managerName + " manager is not created.")
	}

	// Send a job to the manager and return with no errors
	manager.Send(route, data)
	return nil

}

// Await will create and send a request to a defined manager and respond with the completed data
func Await(managerName string, route string, data interface{}) (*Response, error) {

	// Get the manager
	manager, ok := getManager(managerName)

	// If the manager doesn't exist, respond with an error
	if !ok {
		return nil, errors.New(managerName + " manager is not created.")
	}

	// Send a job to the manager and return with no errors
	response := manager.Await(route, data)
	return response, nil

}