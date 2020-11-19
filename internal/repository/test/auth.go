package test

import (
	"errors"

	"github.com/porter-dev/porter/internal/repository"
	"gorm.io/gorm"

	ints "github.com/porter-dev/porter/internal/models/integrations"
)

// KubeIntegrationRepository implements repository.KubeIntegrationRepository
type KubeIntegrationRepository struct {
	canQuery         bool
	kubeIntegrations []*ints.KubeIntegration
}

// NewKubeIntegrationRepository will return errors if canQuery is false
func NewKubeIntegrationRepository(canQuery bool) repository.KubeIntegrationRepository {
	return &KubeIntegrationRepository{
		canQuery,
		[]*ints.KubeIntegration{},
	}
}

// CreateKubeIntegration creates a new kube auth mechanism
func (repo *KubeIntegrationRepository) CreateKubeIntegration(
	am *ints.KubeIntegration,
) (*ints.KubeIntegration, error) {
	if !repo.canQuery {
		return nil, errors.New("Cannot write database")
	}

	repo.kubeIntegrations = append(repo.kubeIntegrations, am)
	am.ID = uint(len(repo.kubeIntegrations))

	return am, nil
}

// ReadKubeIntegration finds a kube auth mechanism by id
func (repo *KubeIntegrationRepository) ReadKubeIntegration(
	id uint,
) (*ints.KubeIntegration, error) {
	if !repo.canQuery {
		return nil, errors.New("Cannot read from database")
	}

	if int(id-1) >= len(repo.kubeIntegrations) || repo.kubeIntegrations[id-1] == nil {
		return nil, gorm.ErrRecordNotFound
	}

	index := int(id - 1)
	return repo.kubeIntegrations[index], nil
}

// ListKubeIntegrationsByProjectID finds all kube auth mechanisms
// for a given project id
func (repo *KubeIntegrationRepository) ListKubeIntegrationsByProjectID(
	projectID uint,
) ([]*ints.KubeIntegration, error) {
	if !repo.canQuery {
		return nil, errors.New("Cannot read from database")
	}

	res := make([]*ints.KubeIntegration, 0)

	for _, kubeAM := range repo.kubeIntegrations {
		if kubeAM.ProjectID == projectID {
			res = append(res, kubeAM)
		}
	}

	return res, nil
}

// OIDCIntegrationRepository implements repository.OIDCIntegrationRepository
type OIDCIntegrationRepository struct {
	canQuery         bool
	oidcIntegrations []*ints.OIDCIntegration
}

// NewOIDCIntegrationRepository will return errors if canQuery is false
func NewOIDCIntegrationRepository(canQuery bool) repository.OIDCIntegrationRepository {
	return &OIDCIntegrationRepository{
		canQuery,
		[]*ints.OIDCIntegration{},
	}
}

// CreateOIDCIntegration creates a new oidc auth mechanism
func (repo *OIDCIntegrationRepository) CreateOIDCIntegration(
	am *ints.OIDCIntegration,
) (*ints.OIDCIntegration, error) {
	if !repo.canQuery {
		return nil, errors.New("Cannot write database")
	}

	repo.oidcIntegrations = append(repo.oidcIntegrations, am)
	am.ID = uint(len(repo.oidcIntegrations))

	return am, nil
}

// ReadOIDCIntegration finds a oidc auth mechanism by id
func (repo *OIDCIntegrationRepository) ReadOIDCIntegration(
	id uint,
) (*ints.OIDCIntegration, error) {
	if !repo.canQuery {
		return nil, errors.New("Cannot read from database")
	}

	if int(id-1) >= len(repo.oidcIntegrations) || repo.oidcIntegrations[id-1] == nil {
		return nil, gorm.ErrRecordNotFound
	}

	index := int(id - 1)
	return repo.oidcIntegrations[index], nil
}

// ListOIDCIntegrationsByProjectID finds all oidc auth mechanisms
// for a given project id
func (repo *OIDCIntegrationRepository) ListOIDCIntegrationsByProjectID(
	projectID uint,
) ([]*ints.OIDCIntegration, error) {
	if !repo.canQuery {
		return nil, errors.New("Cannot read from database")
	}

	res := make([]*ints.OIDCIntegration, 0)

	for _, oidcAM := range repo.oidcIntegrations {
		if oidcAM.ProjectID == projectID {
			res = append(res, oidcAM)
		}
	}

	return res, nil
}

// OIntegrationRepository implements repository.OIntegrationRepository
type OIntegrationRepository struct {
	canQuery      bool
	oIntegrations []*ints.OIntegration
}

// NewOIntegrationRepository will return errors if canQuery is false
func NewOIntegrationRepository(canQuery bool) repository.OIntegrationRepository {
	return &OIntegrationRepository{
		canQuery,
		[]*ints.OIntegration{},
	}
}

// CreateOIntegration creates a new o auth mechanism
func (repo *OIntegrationRepository) CreateOIntegration(
	am *ints.OIntegration,
) (*ints.OIntegration, error) {
	if !repo.canQuery {
		return nil, errors.New("Cannot write database")
	}

	repo.oIntegrations = append(repo.oIntegrations, am)
	am.ID = uint(len(repo.oIntegrations))

	return am, nil
}

// ReadOIntegration finds a o auth mechanism by id
func (repo *OIntegrationRepository) ReadOIntegration(
	id uint,
) (*ints.OIntegration, error) {
	if !repo.canQuery {
		return nil, errors.New("Cannot read from database")
	}

	if int(id-1) >= len(repo.oIntegrations) || repo.oIntegrations[id-1] == nil {
		return nil, gorm.ErrRecordNotFound
	}

	index := int(id - 1)
	return repo.oIntegrations[index], nil
}

// ListOIntegrationsByProjectID finds all o auth mechanisms
// for a given project id
func (repo *OIntegrationRepository) ListOIntegrationsByProjectID(
	projectID uint,
) ([]*ints.OIntegration, error) {
	if !repo.canQuery {
		return nil, errors.New("Cannot read from database")
	}

	res := make([]*ints.OIntegration, 0)

	for _, oAM := range repo.oIntegrations {
		if oAM.ProjectID == projectID {
			res = append(res, oAM)
		}
	}

	return res, nil
}

// AWSIntegrationRepository implements repository.AWSIntegrationRepository
type AWSIntegrationRepository struct {
	canQuery        bool
	awsIntegrations []*ints.AWSIntegration
}

// NewAWSIntegrationRepository will return errors if canQuery is false
func NewAWSIntegrationRepository(canQuery bool) repository.AWSIntegrationRepository {
	return &AWSIntegrationRepository{
		canQuery,
		[]*ints.AWSIntegration{},
	}
}

// CreateAWSIntegration creates a new aws auth mechanism
func (repo *AWSIntegrationRepository) CreateAWSIntegration(
	am *ints.AWSIntegration,
) (*ints.AWSIntegration, error) {
	if !repo.canQuery {
		return nil, errors.New("Cannot write database")
	}

	repo.awsIntegrations = append(repo.awsIntegrations, am)
	am.ID = uint(len(repo.awsIntegrations))

	return am, nil
}

// ReadAWSIntegration finds a aws auth mechanism by id
func (repo *AWSIntegrationRepository) ReadAWSIntegration(
	id uint,
) (*ints.AWSIntegration, error) {
	if !repo.canQuery {
		return nil, errors.New("Cannot read from database")
	}

	if int(id-1) >= len(repo.awsIntegrations) || repo.awsIntegrations[id-1] == nil {
		return nil, gorm.ErrRecordNotFound
	}

	index := int(id - 1)
	return repo.awsIntegrations[index], nil
}

// ListAWSIntegrationsByProjectID finds all aws auth mechanisms
// for a given project id
func (repo *AWSIntegrationRepository) ListAWSIntegrationsByProjectID(
	projectID uint,
) ([]*ints.AWSIntegration, error) {
	if !repo.canQuery {
		return nil, errors.New("Cannot read from database")
	}

	res := make([]*ints.AWSIntegration, 0)

	for _, awsAM := range repo.awsIntegrations {
		if awsAM.ProjectID == projectID {
			res = append(res, awsAM)
		}
	}

	return res, nil
}

// GCPIntegrationRepository implements repository.GCPIntegrationRepository
type GCPIntegrationRepository struct {
	canQuery        bool
	gcpIntegrations []*ints.GCPIntegration
}

// NewGCPIntegrationRepository will return errors if canQuery is false
func NewGCPIntegrationRepository(canQuery bool) repository.GCPIntegrationRepository {
	return &GCPIntegrationRepository{
		canQuery,
		[]*ints.GCPIntegration{},
	}
}

// CreateGCPIntegration creates a new gcp auth mechanism
func (repo *GCPIntegrationRepository) CreateGCPIntegration(
	am *ints.GCPIntegration,
) (*ints.GCPIntegration, error) {
	if !repo.canQuery {
		return nil, errors.New("Cannot write database")
	}

	repo.gcpIntegrations = append(repo.gcpIntegrations, am)
	am.ID = uint(len(repo.gcpIntegrations))

	return am, nil
}

// ReadGCPIntegration finds a gcp auth mechanism by id
func (repo *GCPIntegrationRepository) ReadGCPIntegration(
	id uint,
) (*ints.GCPIntegration, error) {
	if !repo.canQuery {
		return nil, errors.New("Cannot read from database")
	}

	if int(id-1) >= len(repo.gcpIntegrations) || repo.gcpIntegrations[id-1] == nil {
		return nil, gorm.ErrRecordNotFound
	}

	index := int(id - 1)
	return repo.gcpIntegrations[index], nil
}

// ListGCPIntegrationsByProjectID finds all gcp auth mechanisms
// for a given project id
func (repo *GCPIntegrationRepository) ListGCPIntegrationsByProjectID(
	projectID uint,
) ([]*ints.GCPIntegration, error) {
	if !repo.canQuery {
		return nil, errors.New("Cannot read from database")
	}

	res := make([]*ints.GCPIntegration, 0)

	for _, gcpAM := range repo.gcpIntegrations {
		if gcpAM.ProjectID == projectID {
			res = append(res, gcpAM)
		}
	}

	return res, nil
}
