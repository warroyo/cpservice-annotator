package hook

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/prometheus/common/log"
	corev1 "k8s.io/api/core/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/mutate,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod.kb.io

type ServiceMutate struct {
	Name    string
	Client  client.Client
	decoder *admission.Decoder
}

type Config struct {
	Containers []corev1.Container `yaml:"containers"`
}

func shoudInject(service *corev1.Service, k8client client.Client, ctx context.Context) (bool, error, string) {
	var shouldInject bool
	ownerRefs := service.GetOwnerReferences()
	var cluster clusterv1.Cluster
	if len(ownerRefs) > 0 && ownerRefs[0].Kind == "Cluster" {
		clusterName := ownerRefs[0].Name
		if err := k8client.Get(ctx, client.ObjectKey{
			Namespace: service.Namespace,
			Name:      clusterName}, &cluster); err != nil {
			log.Error(err, "unable to fetch Cluster")

			// we'll ignore not-found errors, since they can't be fixed by an immediate
			// requeue (we'll need to wait for a new notification), and we can get them
			// on deleted requests.
			return false, err, ""
		}
	} else {
		shouldInject = false
	}
	aviinfrasetting, ok := cluster.Annotations["cpservicemutate.field.vmware.com/aviinfrasetting"]
	shouldInject = ok

	if shouldInject {
		alreadyUpdated, err := strconv.ParseBool(service.Annotations["cpservicemutate.field.vmware.com/status"])

		if err == nil && alreadyUpdated {
			shouldInject = false
		}
	}

	log.Info("Should Inject: ", shouldInject)

	return shouldInject, nil, aviinfrasetting
}

func (si *ServiceMutate) Handle(ctx context.Context, req admission.Request) admission.Response {
	service := &corev1.Service{}

	err := si.decoder.Decode(req, service)
	if err != nil {
		log.Info("service-mutator: cannot decode")
		return admission.Errored(http.StatusBadRequest, err)
	}

	if service.Annotations == nil {
		service.Annotations = map[string]string{}
	}

	shoudInject, err, aviinfrasetting := shoudInject(service, si.Client, ctx)

	if shoudInject {
		log.Info("Injecting annotations...")

		service.Annotations["cpservicemutate.field.vmware.com/status"] = "true"
		service.Annotations["aviinfrasetting.ako.vmware.com/name"] = aviinfrasetting

		log.Info("ServiceMutate ", si.Name, " injected.")
	} else {
		log.Info("Inject not needed.")
	}

	marshaledPod, err := json.Marshal(service)

	if err != nil {
		log.Info("Service-Mutator: cannot marshal")
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

// InjectDecoder injects the decoder.
func (si *ServiceMutate) InjectDecoder(d *admission.Decoder) error {
	si.decoder = d
	return nil
}
