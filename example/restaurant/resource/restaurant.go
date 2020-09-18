package resource

import (
	"github.com/go-chassis/go-chassis/v2/pkg/runtime"
	"github.com/go-chassis/go-chassis/v2/server/restful"
	"net/http"
)

// RestFulRouterA is a struct used for implementation of restfull router program
type RestaurantResource struct {
}

// Equal is method to compare given num and slice sum
func (r *RestaurantResource) Get(context *restful.Context) {
	context.Write([]byte("restaurant version:" + runtime.Version))
}

// URLPatterns helps to respond for corresponding API calls
func (r *RestaurantResource) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/v1/restaurant/{id}", ResourceFuncName: "Get"},
	}
}
