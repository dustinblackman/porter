package release

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/porter-dev/porter/internal/kubernetes"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

type KNativeResourceRequest struct {
	api           string
	name          string
	namespace     string
	resource      string
	labelSelector string
}

type KNativeYamls struct {
	sync.Mutex

	client *rest.RESTClient
	Wg     *sync.WaitGroup
	YAMLs  []map[string]interface{}
}

func (k *KNativeYamls) AppendYaml(yaml map[string]interface{}) {
	k.Lock()
	defer k.Unlock()

	k.YAMLs = append(k.YAMLs, yaml)
}

func (k *KNativeYamls) RequestYaml(resourceRequest KNativeResourceRequest) {
	defer k.Wg.Done()

	res, err := makeKNativeResourceRequest(k.client, resourceRequest)
	if err != nil {
		fmt.Errorf("knative request failed", err)
		return
	}

	if res == nil {
		return
	}

	k.AppendYaml(res)
}

func makeKNativeResourceRequest(client *rest.RESTClient, resourceRequest KNativeResourceRequest) (map[string]interface{}, error) {
	req := client.Get().AbsPath(resourceRequest.api).
		Namespace(resourceRequest.namespace).
		Resource(resourceRequest.resource)

	if resourceRequest.labelSelector != "" {
		req = req.Param("labelSelector", resourceRequest.labelSelector)
	}

	res, err := req.DoRaw(context.TODO())
	if err != nil {
		return nil, err
	}

	var jsonRes map[string]interface{}
	json.Unmarshal(res, &jsonRes)

	items := jsonRes["items"].([]interface{})
	if len(items) == 0 {
		return nil, nil
	}
	firstEntry := items[0].(map[string]interface{})
	metadata := firstEntry["metadata"].(map[string]interface{})
	name := metadata["name"].(string)

	req = client.Get().AbsPath(resourceRequest.api).
		Namespace(resourceRequest.namespace).
		Resource(resourceRequest.resource).
		Name(name)

	res, err = req.DoRaw(context.TODO())
	if err != nil {
		return nil, err
	}

	var jsonResSingle map[string]interface{}
	json.Unmarshal(res, &jsonResSingle)

	return jsonResSingle, nil
}

func getKnativeYAMLs(agent *kubernetes.Agent, namespace string, controllersOnly bool) ([]map[string]interface{}, error) {
	restConf, err := agent.RESTClientGetter.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	restConf.NegotiatedSerializer = runtime.NewSimpleNegotiatedSerializer(runtime.SerializerInfo{})
	restConf.GroupVersion = &schema.GroupVersion{
		Group:   "api",
		Version: "v1",
	}
	restClient, err := rest.RESTClientFor(restConf)
	if err != nil {
		return nil, err
	}

	ksvc, err := makeKNativeResourceRequest(restClient, KNativeResourceRequest{
		api:           "/apis/serving.knative.dev/v1",
		namespace:     namespace,
		resource:      "services",
		labelSelector: "",
	})
	if err != nil {
		panic(err)
	}

	knativeYamlsReq := KNativeYamls{
		client: restClient,
		Wg:     &sync.WaitGroup{},
		YAMLs:  []map[string]interface{}{},
	}
	knativeYamlsReq.AppendYaml(ksvc)

	resourceRequests := []KNativeResourceRequest{}
	if !controllersOnly {
		resourceRequests = append(resourceRequests, []KNativeResourceRequest{
			{
				api:           "/apis/serving.knative.dev/v1",
				name:          namespace,
				namespace:     namespace,
				resource:      "routes",
				labelSelector: "",
			},
			{
				api:           "/apis/serving.knative.dev/v1",
				name:          namespace,
				namespace:     namespace,
				resource:      "configurations",
				labelSelector: "",
			},
			{
				api:           "/apis/projectcontour.io/v1",
				name:          namespace,
				namespace:     namespace,
				resource:      "httpproxies",
				labelSelector: "",
			},
		}...)
	}

	revisionStatus := ksvc["status"].(map[string]interface{})
	latestCreatedRevisionName := revisionStatus["latestCreatedRevisionName"].(string)
	latestReadyRevisionName := revisionStatus["latestReadyRevisionName"].(string)
	revisions := []string{latestReadyRevisionName}
	if latestCreatedRevisionName != latestCreatedRevisionName {
		revisions = append(revisions, latestCreatedRevisionName)
	}

	for _, revisionName := range revisions {
		resourceRequests = append(resourceRequests, KNativeResourceRequest{
			api:           "/apis/apps/v1",
			name:          namespace,
			namespace:     namespace,
			resource:      "deployments",
			labelSelector: fmt.Sprintf("serving.knative.dev/revision=%s", revisionName),
		})

		// TODO need to add configurations to show failed deployments more clearly.
		revisionResourceRequests := []KNativeResourceRequest{
			{
				api:           "/api/v1",
				name:          revisionName,
				namespace:     namespace,
				resource:      "services",
				labelSelector: "networking.internal.knative.dev/serviceType=Public",
			},
			{
				api:           "/api/v1",
				name:          fmt.Sprintf("%s-private", revisionName),
				namespace:     namespace,
				resource:      "services",
				labelSelector: "networking.internal.knative.dev/serviceType=Private",
			},
			{
				api:           "/apis/autoscaling/v2",
				name:          revisionName,
				namespace:     namespace,
				resource:      "horizontalpodautoscalers",
				labelSelector: "",
			},
			{
				api:           "/apis/autoscaling.internal.knative.dev/v1alpha1",
				name:          revisionName,
				namespace:     namespace,
				resource:      "podautoscalers",
				labelSelector: "",
			},
		}

		if !controllersOnly {
			resourceRequests = append(resourceRequests, revisionResourceRequests...)
		}
	}

	domainName, ok := ksvc["metadata"].(map[string]interface{})["annotations"].(map[string]interface{})["external-dns.alpha.kubernetes.io/hostname"].(string)
	if ok && !controllersOnly {
		resourceRequests = append(resourceRequests, KNativeResourceRequest{
			api:       "/apis/cert-manager.io/v1",
			name:      domainName,
			namespace: namespace,
			resource:  "certificates",
		})
	}

	for _, resourceRequest := range resourceRequests {
		knativeYamlsReq.Wg.Add(1)
		go knativeYamlsReq.RequestYaml(resourceRequest)
	}

	knativeYamlsReq.Wg.Wait()
	return knativeYamlsReq.YAMLs, nil
}

func addKnativeYamls(manifest string, namespace string, yamls []map[string]interface{}, agent *kubernetes.Agent, controllersOnly bool) []map[string]interface{} {
	if agent != nil && strings.Contains(manifest, "serving.knative.dev/v1") {
		knativeYamls, err := getKnativeYAMLs(agent, namespace, controllersOnly)
		if err != nil {
			fmt.Println("[KNATIVE ERROR] ", err)
		} else {
			yamls = append(yamls, knativeYamls...)
		}
	}

	return yamls
}
