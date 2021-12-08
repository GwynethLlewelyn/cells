package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pydio/cells/v4/common/config/runtime"
	servicecontext "github.com/pydio/cells/v4/common/service/context"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/emicklei/go-restful"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/rs/cors"
	"go.uber.org/zap"

	"github.com/pydio/cells/v4/common"
	"github.com/pydio/cells/v4/common/log"
	"github.com/pydio/cells/v4/common/proto/rest"
)

var (
	swaggerSyncOnce       = &sync.Once{}
	swaggerJSONStrings    []string
	swaggerMergedDocument *loads.Document
)

// RegisterSwaggerJSON receives a json string and adds it to the swagger definition
func RegisterSwaggerJSON(json string) {
	swaggerSyncOnce = &sync.Once{}
	swaggerJSONStrings = append(swaggerJSONStrings, json)
}

func init() {
	// Instanciate restful framework
	restful.RegisterEntityAccessor("application/json", new(ProtoEntityReaderWriter))
}

// WebHandler defines what functions a web handler must answer to
type WebHandler interface {
	SwaggerTags() []string
	Filter() func(string) string
}

// WithWeb returns a web handler
func WithWeb(handler func() WebHandler) ServiceOption {
	return func(o *ServiceOptions) {
		// Making sure the runtime is correct
		if o.Fork && !runtime.IsFork() {
			return
		}

		o.Server = servicecontext.GetServer(o.Context,"http")
		o.serverStart = func() error {
			var mux *http.ServeMux
			if !o.Server.As(&mux) {
				return fmt.Errorf("server is not a mux")
			}

			// TODO v4
			// if rateLimit, err := strconv.Atoi(os.Getenv("WEB_RATE_LIMIT")); err == nil {
			// 	opts = append(opts, micro.WrapHandler(limiter.NewHandlerWrapper(rateLimit)))
			//}

			// meta := registry.BuildServiceMeta()
			// meta["description"] = o.Description

			// svc.Init(
			// 	micro.Metadata(registry.BuildServiceMeta()),
			// )

			ctx := o.Context

			rootPath := "/a/" + strings.TrimPrefix(o.Name, common.ServiceRestNamespace_)
			log.Logger(ctx).Info("starting", zap.String("service", o.Name), zap.String("hook router to", rootPath))

			ws := new(restful.WebService)
			ws.Consumes(restful.MIME_JSON, "application/x-www-form-urlencoded", "multipart/form-data")
			ws.Produces(restful.MIME_JSON, restful.MIME_OCTET, restful.MIME_XML)
			ws.Path(rootPath)

			h := handler()
			swaggerTags := h.SwaggerTags()
			filter := h.Filter()

			f := reflect.ValueOf(h)

			for path, pathItem := range SwaggerSpec().Spec().Paths.Paths {
				if pathItem.Get != nil {
					shortPath, method := operationToRoute(rootPath, swaggerTags, path, pathItem.Get, filter, f)
					if shortPath != "" {
						ws.Route(ws.GET(shortPath).To(method))
					}
				}
				if pathItem.Delete != nil {
					shortPath, method := operationToRoute(rootPath, swaggerTags, path, pathItem.Delete, filter, f)
					if shortPath != "" {
						ws.Route(ws.DELETE(shortPath).To(method))
					}
				}
				if pathItem.Put != nil {
					shortPath, method := operationToRoute(rootPath, swaggerTags, path, pathItem.Put, filter, f)
					if shortPath != "" {
						ws.Route(ws.PUT(shortPath).To(method))
					}
				}
				if pathItem.Patch != nil {
					shortPath, method := operationToRoute(rootPath, swaggerTags, path, pathItem.Patch, filter, f)
					if shortPath != "" {
						ws.Route(ws.PATCH(shortPath).To(method))
					}
				}
				if pathItem.Head != nil {
					shortPath, method := operationToRoute(rootPath, swaggerTags, path, pathItem.Head, filter, f)
					if shortPath != "" {
						ws.Route(ws.HEAD(shortPath).To(method))
					}
				}
				if pathItem.Post != nil {
					shortPath, method := operationToRoute(rootPath, swaggerTags, path, pathItem.Post, filter, f)
					if shortPath != "" {
						ws.Route(ws.POST(shortPath).To(method))
					}
				}
			}

			wc := restful.NewContainer()
			// Enable globally gzip,deflate encoding globally
			wc.EnableContentEncoding(true)
			wc.Add(ws)

			// var e error
			wrapped := http.Handler(wc)

			// for _, wrap := range o.webHandlerWraps {
			// 	wrapped = wrap(wrapped)
			//}

			// if wrapped, e = NewConfigHTTPHandlerWrapper(wrapped, name); e != nil {
			//	return e
			//}
			// wrapped = NewLogHTTPHandlerWrapper(wrapped, name)

			wrapped = cors.Default().Handler(wrapped)

			mux.Handle(ws.RootPath(), wrapped)
			mux.Handle(ws.RootPath()+"/", wrapped)

			return nil
		}

		return
	}
}

func WithWebStop(handler func() WebHandler) ServiceOption {
	return func(o *ServiceOptions) {
		o.serverStop = func() error {
			// TODO v4 - unregister all services
			return nil
		}
	}
}

func operationToRoute(rootPath string, swaggerTags []string, path string, operation *spec.Operation, pathFilter func(string) string, handlerValue reflect.Value) (string, func(req *restful.Request, rsp *restful.Response)) {

	if !containsTags(operation, swaggerTags) {
		return "", nil
	}

	method := handlerValue.MethodByName(operation.ID)
	if method.IsValid() {
		casted := method.Interface().(func(req *restful.Request, rsp *restful.Response))
		shortPath := strings.TrimPrefix("/a" + path, rootPath)
		if shortPath == "" {
			shortPath = "/"
		}
		if pathFilter != nil {
			shortPath = pathFilter(shortPath)
		}

		log.Logger(context.Background()).Debug("Registering path " + shortPath + " to handler method " + operation.ID)
		return shortPath, casted
	}

	log.Logger(context.Background()).Debug("Cannot find method " + operation.ID + " on handler, ignoring GET for path " + path)
	return "", nil
}

func containsTags(operation *spec.Operation, filtersTags []string) (found bool) {
	for _, tag := range operation.Tags {
		for _, filter := range filtersTags {
			if tag == filter {
				found = true
				break
			}
		}
	}
	return
}

// SwaggerSpec returns the swagger specification as a document
func SwaggerSpec() *loads.Document {
	swaggerSyncOnce.Do(func() {
		var swaggerDocuments []*loads.Document
		for _, data := range append([]string{rest.SwaggerJson}, swaggerJSONStrings...) {
			// Reading swagger json
			rawMessage := new(json.RawMessage)
			if e := json.Unmarshal([]byte(data), rawMessage); e != nil {
				log.Fatal("Failed to load swagger row data", zap.Error(e))
			}
			j, err := loads.Analyzed(*rawMessage, "")
			if err != nil {
				log.Fatal("Failed to load swagger", zap.Error(err))
			}

			swaggerDocuments = append(swaggerDocuments, j)
		}

		for _, j := range swaggerDocuments {
			if swaggerMergedDocument == nil { // First pass
				swaggerMergedDocument = j
			} else { // other passes : merge all Paths
				for p, i := range j.Spec().Paths.Paths {
					if existing, ok := swaggerMergedDocument.Spec().Paths.Paths[p]; ok {
						if i.Get != nil {
							existing.Get = i.Get
						}
						if i.Put != nil {
							existing.Put = i.Put
						}
						if i.Post != nil {
							existing.Post = i.Post
						}
						if i.Options != nil {
							existing.Options = i.Options
						}
						if i.Delete != nil {
							existing.Delete = i.Delete
						}
						if i.Head != nil {
							existing.Head = i.Head
						}
						swaggerMergedDocument.Spec().Paths.Paths[p] = existing
					} else {
						swaggerMergedDocument.Spec().Paths.Paths[p] = i
					}
				}
				for name, schema := range j.Spec().Definitions {
					swaggerMergedDocument.Spec().Definitions[name] = schema
				}
			}
		}
	})

	if swaggerMergedDocument == nil {
		// log.Logger(context.Background()).Fatal("Could not find any valid json spec for swagger")
	}

	return swaggerMergedDocument
}
