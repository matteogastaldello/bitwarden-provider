package secrets

import (
	"context"
	"os"

	prv1 "github.com/krateoplatformops/provider-runtime/apis/common/v1"
	bwclient "github.com/matteogastaldello/bitwarden-provider/internal/clients"
	"github.com/matteogastaldello/bitwarden-provider/internal/clients/secrets"
	"github.com/pkg/errors"

	"k8s.io/client-go/tools/record"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/krateoplatformops/provider-runtime/pkg/controller"
	"github.com/krateoplatformops/provider-runtime/pkg/event"
	"github.com/krateoplatformops/provider-runtime/pkg/logging"
	"github.com/krateoplatformops/provider-runtime/pkg/ratelimiter"
	"github.com/krateoplatformops/provider-runtime/pkg/reconciler"
	"github.com/krateoplatformops/provider-runtime/pkg/resource"

	bwv1 "github.com/matteogastaldello/bitwarden-provider/api/secret/v1"
)

const (
	errNotBitwardenSecret = "managed resource is not a BitwardenSecret custom resource"
	errToUnlockVault      = "error unlocking the vault"
)

func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := reconciler.ControllerName(bwv1.BitwardenSecretGroupKind)

	log := o.Logger.WithValues("controller", name)

	recorder := mgr.GetEventRecorderFor(name)

	r := reconciler.NewReconciler(mgr,
		resource.ManagedKind(bwv1.BitwardenSecretGroupVersionKind),
		reconciler.WithExternalConnecter(&connector{
			kube:     mgr.GetClient(),
			log:      log,
			recorder: recorder,
		}),
		reconciler.WithPollInterval(o.PollInterval),
		reconciler.WithLogger(log),
		reconciler.WithRecorder(event.NewAPIRecorder(recorder)))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		For(&bwv1.BitwardenSecret{}).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}

type connector struct {
	kube     client.Client
	log      logging.Logger
	recorder record.EventRecorder
}

func (c *connector) Connect(ctx context.Context, mg resource.Managed) (reconciler.ExternalClient, error) {
	_, ok := mg.(*bwv1.BitwardenSecret)
	if !ok {
		return nil, errors.New(errNotBitwardenSecret)
	}

	// opts, err := resolvers.ResolveConnectorConfig(ctx, c.kube, cr.Spec.ConnectorConfigRef)
	// if err != nil {
	// 	return nil, err
	// }
	// opts.Verbose = meta.IsVerbose(cr)

	return &external{
		kube:  c.kube,
		log:   c.log,
		bwCli: bwclient.NewClient(),
		rec:   c.recorder,
	}, nil
}

type external struct {
	kube  client.Client
	log   logging.Logger
	bwCli *bwclient.Client
	rec   record.EventRecorder
}

func (e *external) Observe(ctx context.Context, mg resource.Managed) (reconciler.ExternalObservation, error) {
	resp, err := bwclient.Unlock(ctx, e.bwCli, os.Getenv("BW_PASSWORD"))
	if err != nil {
		return reconciler.ExternalObservation{}, err
	}
	if !resp.Success {
		return reconciler.ExternalObservation{}, errors.New(errToUnlockVault)
	}

	cr, ok := mg.(*bwv1.BitwardenSecret)
	if !ok {
		return reconciler.ExternalObservation{}, errors.New(errNotBitwardenSecret)
	}

	status := cr.Status.DeepCopy()
	spec := cr.Spec.DeepCopy()

	ok, err = secrets.Exists(ctx, e.bwCli, status.SecretId)
	if err != nil {
		return reconciler.ExternalObservation{}, err
	}

	if ok {
		e.log.Debug("Secret already exists", "id", status.SecretId, "name", spec.Name)
		e.rec.Eventf(cr, corev1.EventTypeNormal, "AlredyExists", "Secret '%s/%s' already exists", status.SecretId, spec.Name)

		cr.SetConditions(prv1.Available())
		return reconciler.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: true,
		}, nil
	}

	e.log.Debug("Repo does not exists", "id", status.SecretId, "name", spec.Name)

	return reconciler.ExternalObservation{
		ResourceExists:   false,
		ResourceUpToDate: true,
	}, nil
}

func (e *external) Update(ctx context.Context, mg resource.Managed) error {
	return nil // noop
}

func (e *external) Delete(ctx context.Context, mg resource.Managed) error {
	return nil // noop
}

func (e *external) Create(ctx context.Context, mg resource.Managed) error {
	return nil // noop
}
