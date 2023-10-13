package secrets

import (
	"context"
	"fmt"

	prv1 "github.com/krateoplatformops/provider-runtime/apis/common/v1"
	bwclient "github.com/matteogastaldello/bitwarden-provider/internal/clients"
	"github.com/matteogastaldello/bitwarden-provider/internal/clients/secrets"
	"github.com/pkg/errors"

	"k8s.io/client-go/tools/record"

	"github.com/krateoplatformops/provider-runtime/pkg/meta"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/krateoplatformops/provider-runtime/pkg/controller"
	"github.com/krateoplatformops/provider-runtime/pkg/event"
	"github.com/krateoplatformops/provider-runtime/pkg/logging"
	"github.com/krateoplatformops/provider-runtime/pkg/ratelimiter"
	"github.com/krateoplatformops/provider-runtime/pkg/reconciler"
	"github.com/krateoplatformops/provider-runtime/pkg/resource"

	"github.com/matteogastaldello/bitwarden-provider/internal/clients/unlocker"

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
	fmt.Print("\nOBSERVE\n")
	cr, ok := mg.(*bwv1.BitwardenSecret)
	if !ok {
		return reconciler.ExternalObservation{}, errors.New(errNotBitwardenSecret)
	}

	resp, err := unlocker.Unlock(ctx, e.bwCli, &e.kube, cr.Spec.ConnectorConfigRef)

	if err != nil {
		return reconciler.ExternalObservation{}, err
	}
	if !resp.Success {
		return reconciler.ExternalObservation{}, errors.New(errToUnlockVault)
	}

	status := cr.Status.DeepCopy()
	spec := cr.Spec.DeepCopy()

	res, err := secrets.Exists(ctx, e.bwCli, status.SecretId)
	if err != nil && res == nil {
		return reconciler.ExternalObservation{}, err
	}

	if *res {
		e.log.Debug("Secret already exists", "id", status.SecretId, "name", spec.Name)
		e.rec.Eventf(cr, corev1.EventTypeNormal, "AlredyExists", "Secret '%s/%s' already exists", status.SecretId, spec.Name)

		cr.SetConditions(prv1.Available())
		return reconciler.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: true,
		}, nil
	}

	e.log.Debug("Secret does not exists", "id", status.SecretId, "name", spec.Name)

	return reconciler.ExternalObservation{
		ResourceExists:   false,
		ResourceUpToDate: true,
	}, nil
}

func (e *external) Update(ctx context.Context, mg resource.Managed) error {
	return nil // noop
}

func (e *external) Delete(ctx context.Context, mg resource.Managed) error {
	fmt.Println("DELETE")
	cr, ok := mg.(*bwv1.BitwardenSecret)
	if !ok {
		return errors.New(errNotBitwardenSecret)
	}

	//test AnnotationKeyManagementPolicy - non so se sia l'utilizzo corretto di questa impostazione
	if condition := meta.IsActionAllowed(cr, meta.ActionDelete); !condition {
		return errors.New("Action not allowed")
	}

	cr.SetConditions(prv1.Deleting())

	status := cr.Status.DeepCopy()
	spec := cr.Spec.DeepCopy()
	ok, err := secrets.Delete(ctx, e.bwCli, status.SecretId)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("Something went wrong deleting secret on remote server")
	}

	e.log.Debug("Secret Deleted", "id", status.SecretId, "name", spec.Name)
	e.rec.Eventf(cr, corev1.EventTypeNormal, "SecretDeleted", "Secret '%s/%s' deleted", status.SecretId, spec.Name)

	return nil
}

func (e *external) Create(ctx context.Context, mg resource.Managed) error {
	fmt.Println("CREATE")
	cr, ok := mg.(*bwv1.BitwardenSecret)
	if !ok {
		return errors.New(errNotBitwardenSecret)
	}

	cr.SetConditions(prv1.Creating())

	spec := cr.Spec.DeepCopy()
	sec, err := secrets.Create(ctx, e.bwCli, spec.Secret)
	if err != nil {
		return err
	}
	cr.Status.SecretId = sec.Id
	e.kube.Status().Update(ctx, cr)

	e.log.Debug("Secret created", "id", sec.Id, "name", spec.Name)
	e.rec.Eventf(cr, corev1.EventTypeNormal, "SecretCreated", "Secret '%s/%s' created", spec.Id, spec.Name)

	return nil
}
