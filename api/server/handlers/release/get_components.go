package release

import (
	"net/http"

	"github.com/porter-dev/porter/api/server/authz"
	"github.com/porter-dev/porter/api/server/handlers"
	"github.com/porter-dev/porter/api/server/shared"
	"github.com/porter-dev/porter/api/server/shared/config"
	"github.com/porter-dev/porter/api/types"
	"github.com/porter-dev/porter/internal/helm/grapher"
	"github.com/porter-dev/porter/internal/models"
	"github.com/stefanmcshane/helm/pkg/release"
)

type GetComponentsHandler struct {
	handlers.PorterHandlerReadWriter
	authz.KubernetesAgentGetter
}

func NewGetComponentsHandler(
	config *config.Config,
	writer shared.ResultWriter,
) *GetComponentsHandler {
	return &GetComponentsHandler{
		PorterHandlerReadWriter: handlers.NewDefaultPorterHandler(config, nil, writer),
		KubernetesAgentGetter:   authz.NewOutOfClusterAgentGetter(config),
	}
}

func (c *GetComponentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	helmRelease, _ := r.Context().Value(types.ReleaseScope).(*release.Release)

	yamlArr := grapher.ImportMultiDocYAML([]byte(helmRelease.Manifest))
	cluster, _ := r.Context().Value(types.ClusterScope).(*models.Cluster)
	agent, _ := c.GetAgent(r, cluster, "")
	yamlArr = addKnativeYamls(helmRelease.Manifest, helmRelease.Namespace, yamlArr, agent, true)

	objects := grapher.ParseObjs(yamlArr, helmRelease.Namespace)

	parsed := grapher.ParsedObjs{
		Objects: objects,
	}

	parsed.GetControlRel()
	parsed.GetLabelRel()
	parsed.GetSpecRel()

	c.WriteResult(w, r, parsed)
}
