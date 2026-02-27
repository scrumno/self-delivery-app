package collector

import (
    "net/http"
    "reflect"
    "sync"
	"fmt"
	"strings"

    "github.com/gorilla/mux"
    "github.com/scrumno/scrumno-api/internal/api/utils/postman"
)

type EndpointInfo struct {
    Method    string       `json:"method"`
    Path      string       `json:"path"`
    InputType reflect.Type `json:"-"`
	Group     string       `json:"group,omitempty"`
}

type EndpointCollector struct {
    mu        sync.RWMutex 
    endpoints map[string]EndpointInfo
	groups    map[string][]string
}

func NewEndpointCollector() *EndpointCollector {
    return &EndpointCollector{
        endpoints: make(map[string]EndpointInfo),
		groups:    make(map[string][]string),
    }
}

func (c *EndpointCollector) AddEndpoint(method, path string, inputType reflect.Type, prefix string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    key := fmt.Sprintf("%s:%s", method, path)

	if strings.TrimSpace(prefix) == "" {
        prefix = "other"
    }
	
    c.endpoints[key] = EndpointInfo{
        Method:    method,
        Path:      path,
        InputType: inputType,
		Group:     prefix,
    }

	c.groups[prefix] = append(c.groups[prefix], key)
}

func (c *EndpointCollector) GetAllEndpointsByGroup() map[string][]map[string]interface{} {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    result := make(map[string][]map[string]interface{})
    
    for group, keys := range c.groups {
        groupEndpoints := make([]map[string]interface{}, 0, len(keys))
        for _, key := range keys {
            ep := c.endpoints[key]
            groupEndpoints = append(groupEndpoints, map[string]interface{}{
                "method":    ep.Method,
                "path":      group + ep.Path,
                "inputType": ep.InputType,
            })
        }
        result[group] = groupEndpoints
    }

    return result
}

func (c *EndpointCollector) HandleFuncWithPostman(
	router *mux.Router, 
	prefix string,
    handleFunc http.HandlerFunc,
    inputType reflect.Type, 
	method string, 
	path string,
) *mux.Route {
    c.AddEndpoint(method, path, inputType, prefix)
    return router.HandleFunc(path, handleFunc).Methods(method)
}

func (c *EndpointCollector) GeneratePostmanCollections() error {
    breakpointGroups := c.GetAllEndpointsByGroup()
    
    return utils.Generate(breakpointGroups)
}