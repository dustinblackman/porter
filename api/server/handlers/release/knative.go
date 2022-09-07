package release

import (
	"context"
	"encoding/json"
	"fmt"
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

func (k *KNativeYamls) appendYaml(yaml map[string]interface{}) {
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

	k.appendYaml(res)
}

func makeKNativeResourceRequest(client *rest.RESTClient, resourceRequest KNativeResourceRequest) (map[string]interface{}, error) {
	req := client.Get().AbsPath("/apis/" + resourceRequest.api).
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

	req = client.Get().AbsPath("/apis/" + resourceRequest.api).
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

func getKnativeYAMLs(agent *kubernetes.Agent, namespace string) ([]map[string]interface{}, error) {
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
		api:           "serving.knative.dev/v1",
		namespace:     namespace,
		resource:      "services",
		labelSelector: "",
	})
	if err != nil {
		panic(err)
	}

	revisionName := ksvc["status"].(map[string]interface{})["latestCreatedRevisionName"].(string)

	knativeYamlsReq := KNativeYamls{
		client: restClient,
		Wg:     &sync.WaitGroup{},
		YAMLs:  []map[string]interface{}{},
	}

	resourceRequests := []KNativeResourceRequest{
		{
			api:           "projectcontour.io/v1",
			name:          namespace,
			namespace:     namespace,
			resource:      "httpproxies",
			labelSelector: "",
		},
		{
			api:           "apps/v1",
			name:          namespace,
			namespace:     namespace,
			resource:      "deployments",
			labelSelector: fmt.Sprintf("serving.knative.dev/revision=%s", revisionName),
		},
	}

	for _, resourceRequest := range resourceRequests {
		knativeYamlsReq.Wg.Add(1)
		go knativeYamlsReq.RequestYaml(resourceRequest)
	}

	knativeYamlsReq.Wg.Wait()
	return knativeYamlsReq.YAMLs, nil
}
