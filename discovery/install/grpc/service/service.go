package service

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/pydio/cells/v5/common"
	"github.com/pydio/cells/v5/common/client/commons/idmc"
	"github.com/pydio/cells/v5/common/proto/idm"
	service2 "github.com/pydio/cells/v5/common/proto/service"
	"github.com/pydio/cells/v5/common/runtime"
	"github.com/pydio/cells/v5/common/server/generic"
	"github.com/pydio/cells/v5/common/service"
	"github.com/pydio/cells/v5/common/telemetry/log"
	gh "github.com/pydio/cells/v5/discovery/install/grpc"
)

var Name = common.ServiceGrpcNamespace_ + common.ServiceInstall

func init() {

	runtime.Register("main", func(ctx context.Context) {
		service.NewService(
			service.Name(Name),
			service.Context(ctx),
			service.Tag(common.ServiceTagDiscovery),
			service.Description("Services Migration"),
			service.WithGRPC(func(ctx context.Context, server grpc.ServiceRegistrar) error {
				handler := new(gh.Handler)
				service2.RegisterMigrateServiceServer(server, handler)
				return nil
			}),
		)
	})

	runtime.Register("controller", func(ctx context.Context) {
		service.NewService(
			service.Name(common.ServiceGenericNamespace_+common.ServiceInstall),
			service.Context(ctx),
			service.Tag(common.ServiceTagGateway),
			service.Description("Cells kubernetes controller"),
			service.WithGeneric(func(ctx context.Context, srv *generic.Server) error {

				if os.Getenv("CELLS_START_K8S_MANAGER") == "true" {
					scheme := k8sruntime.NewScheme()
					_ = clientgoscheme.AddToScheme(scheme)
					_ = batchv1.AddToScheme(scheme)
					_ = corev1.AddToScheme(scheme)

					mgr, err := manager.New(ctrl.GetConfigOrDie(), manager.Options{
						BaseContext: func() context.Context {
							return ctx
						},
						Scheme: scheme,
					})

					if err != nil {
						fmt.Println("Unable to start manager", err)
						os.Exit(1)
					}

					reconciler := &SecretReconciler{
						Client:     mgr.GetClient(),
						Scheme:     mgr.GetScheme(),
						TargetName: os.Getenv("CELLS_START_K8S_MANAGER_SECRET"),
					}

					// Predicates: only reconcile on real changes.
					// Note: Secrets don't bump .metadata.generation on data changes,
					// so we compare relevant fields manually.
					changed := predicate.Funcs{
						CreateFunc: func(e event.CreateEvent) bool {
							//return false
							return allowByTarget(reconciler, e.Object)
						},
						UpdateFunc: func(e event.UpdateEvent) bool {
							if !allowByTarget(reconciler, e.ObjectNew) {
								return false
							}
							oldS, ok1 := e.ObjectOld.(*corev1.Secret)
							newS, ok2 := e.ObjectNew.(*corev1.Secret)
							if !ok1 || !ok2 {
								return true // be safe
							}
							dataChanged := !reflect.DeepEqual(oldS.Data, newS.Data)
							typeChanged := oldS.Type != newS.Type
							lblChanged := !reflect.DeepEqual(oldS.Labels, newS.Labels)
							annChanged := !reflect.DeepEqual(oldS.Annotations, newS.Annotations)
							return dataChanged || typeChanged || lblChanged || annChanged
						},
						DeleteFunc: func(e event.DeleteEvent) bool {
							return allowByTarget(reconciler, e.Object)
						},
						GenericFunc: func(e event.GenericEvent) bool { return false },
					}

					// Build controller.
					if err := builder.
						ControllerManagedBy(mgr).
						For(&corev1.Secret{}, builder.WithPredicates(changed)).
						Watches(&corev1.Secret{}, &handler.EnqueueRequestForObject{}, builder.WithPredicates(changed)). // redundant but explicit
						Complete(reconciler); err != nil {
						panic(err)
					}

					zl := zap.New(log.Logger(ctx).Core())

					ctrl.SetLogger(zapr.NewLogger(zl))
					if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
						fmt.Println("Problem running manager", err)
						os.Exit(1)
					}
				}

				return nil
			}))
	})
}

func allowByTarget(r *SecretReconciler, obj client.Object) bool {
	if r.TargetName == "" {
		// All secrets (optionally limited by mgr Namespace)
		return true
	}
	// Watch only a single Secret name (and namespace if set)
	if r.TargetNamespace != "" && obj.GetNamespace() != r.TargetNamespace {
		return false
	}

	return obj.GetName() == r.TargetName
}

// Reconciler that reacts to Secret changes.
type SecretReconciler struct {
	client.Client
	Scheme *k8sruntime.Scheme

	// Optional: limit to a single secret
	// If both fields are non-empty, only that Secret will trigger reconciles.
	TargetNamespace string
	TargetName      string
}

func (r *SecretReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	logger := log.Logger(ctx)
	logger.Info("Reconciling Secret", zap.String("ns", req.Namespace), zap.String("name", req.Name))

	var secret corev1.Secret
	if err := r.Get(ctx, req.NamespacedName, &secret); err != nil {
		if errors.IsNotFound(err) {
			// Secret deleted; you can react here if you need to clean up.
			logger.Info("Secret deleted", zap.String("ns", req.Namespace), zap.String("name", req.Name))
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// ---- Your logic goes here ----
	username := string(secret.Data["username"])
	password := string(secret.Data["password"])

	if username == "" {
		return reconcile.Result{}, nil
	}

	client := idmc.UserServiceClient(ctx)
	users, err := searchUser(ctx, client, string(username))
	if err != nil {
		logger.Error("Cannot list users for login", zap.String("username", username), zap.Error(err))
		return reconcile.Result{}, err
	}

	for _, user := range users {
		user.Password = string(password)
		if user.Attributes == nil {
			user.Attributes = make(map[string]string, 1)
		}
		user.Attributes["profile"] = common.PydioProfileAdmin

		if _, err := client.CreateUser(ctx, &idm.CreateUserRequest{
			User: user,
		}); err != nil {
			logger.Error("could not update password, skipping and continuing", zap.String("username", username), zap.Error(err))
			return reconcile.Result{}, err
		} else {
			logger.Info("user successfully updated", zap.String("username", username))
		}
	}
	// ------------------------------

	return reconcile.Result{}, nil
}

func searchUser(ctx context.Context, cli idm.UserServiceClient, login string) ([]*idm.User, error) {

	singleQ := &idm.UserSingleQuery{Login: login}
	query, _ := anypb.New(singleQ)

	mainQuery := &service2.Query{SubQueries: []*anypb.Any{query}}

	stream, err := cli.SearchUser(ctx, &idm.SearchUserRequest{Query: mainQuery})
	if err != nil {
		return nil, err
	}

	users := []*idm.User{}

	for {
		response, e := stream.Recv()
		if e != nil {
			break
		}
		if response == nil {
			continue
		}

		currUser := response.GetUser()
		if currUser.IsGroup {
			continue
		}

		if len(users) >= 50 {
			fmt.Println("Maximum of users that can be edited at a time reached. Truncating the list. Please refine you search.")
			break
		}
		users = append(users, currUser)
	}
	return users, nil
}
