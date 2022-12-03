package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/repositories"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/palantir/stacktrace"
	"golang.org/x/sync/errgroup"
)

type ProjectSettingsService struct {
	tracer                         telemetry.Tracer
	logger                         telemetry.Logger
	projectRepository              repositories.ProjectRepository
	whatsappIntegrationRepository  repositories.WhatsappIntegrationRepository
	contentIntegrationRepository   repositories.ContentIntegrationRepository
	projectIntegrationRepository   repositories.ProjectIntegrationRepository
	phoneCallIntegrationRepository repositories.PhoneCallIntegrationRepository
	linkIntegrationRepository      repositories.LinkIntegrationRepository
}

// NewProjectSettingsService creates a new ProjectSettingsService
func NewProjectSettingsService(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	projectRepository repositories.ProjectRepository,
	whatsappIntegrationRepository repositories.WhatsappIntegrationRepository,
	contentIntegrationRepository repositories.ContentIntegrationRepository,
	projectIntegrationRepository repositories.ProjectIntegrationRepository,
	phoneCallIntegrationRepository repositories.PhoneCallIntegrationRepository,
	linkIntegrationRepository repositories.LinkIntegrationRepository,
) (s *ProjectSettingsService) {
	return &ProjectSettingsService{
		logger:                         logger.WithService(fmt.Sprintf("%T", s)),
		tracer:                         tracer,
		projectRepository:              projectRepository,
		whatsappIntegrationRepository:  whatsappIntegrationRepository,
		contentIntegrationRepository:   contentIntegrationRepository,
		projectIntegrationRepository:   projectIntegrationRepository,
		phoneCallIntegrationRepository: phoneCallIntegrationRepository,
		linkIntegrationRepository:      linkIntegrationRepository,
	}
}

// Get returns the entities.ProjectSettings an entities.Project
func (service *ProjectSettingsService) Get(ctx context.Context, userID entities.UserID, projectID uuid.UUID) (*entities.ProjectSettings, error) {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	project, err := service.projectRepository.Load(ctx, userID, projectID)
	if err != nil {
		msg := fmt.Sprintf("cannot load project [%s] for user ID [%s]", projectID, userID)
		return nil, stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	integrations, err := service.projectIntegrationRepository.Fetch(ctx, userID, projectID)
	if err != nil {
		msg := fmt.Sprintf("cannot load project integrations [%s] for user ID [%s]", projectID, userID)
		return nil, stacktrace.PropagateWithCode(err, stacktrace.GetCode(err), msg)
	}

	ctxLogger.Info(fmt.Sprintf("found [%d] integrations for project [%s] and user [%s]", len(integrations), projectID, userID))
	if len(integrations) == 0 {
		return service.createWithoutIntegrations(project), nil
	}

	settings, err := service.fetchInParallel(ctx, userID, service.groupByType(integrations))
	if err != nil {
		msg := fmt.Sprintf("cannot fetch integrations for project [%s] and user [%s]", projectID, userID)
		return nil, stacktrace.Propagate(err, msg)
	}

	return service.creatSortedIntegrations(project, settings, integrations), nil
}

func (service *ProjectSettingsService) fetchInParallel(ctx context.Context, userID entities.UserID, integrationGroups map[entities.IntegrationType][]uuid.UUID) (map[uuid.UUID]*entities.ProjectSettingsIntegration, error) {
	ctx, span, _ := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	lock := sync.RWMutex{}
	integrationSettings := map[uuid.UUID]*entities.ProjectSettingsIntegration{}
	errGroup, ctx := errgroup.WithContext(ctx)
	for key, value := range integrationGroups {
		ids := value
		integrationType := key
		switch key {
		case entities.IntegrationTypeWhatsapp:
			errGroup.Go(func() error {
				integrations, err := service.whatsappIntegrationRepository.FetchMultiple(ctx, userID, ids)
				if err != nil {
					return stacktrace.Propagate(err, fmt.Sprintf("cannot fetch whatsapp integraions for userID [%s] and projects [%+#v]", userID, ids))
				}
				lock.Lock()
				defer lock.Unlock()

				for _, integration := range integrations {
					integrationSettings[integration.ID] = &entities.ProjectSettingsIntegration{
						Type:     integrationType,
						ID:       integration.ID,
						Settings: integration,
					}
				}
				return nil
			})
		case entities.IntegrationTypeContent:
			errGroup.Go(func() error {
				integrations, err := service.contentIntegrationRepository.FetchMultiple(ctx, userID, ids)
				if err != nil {
					return stacktrace.Propagate(err, fmt.Sprintf("cannot fetch [%s] integraions for userID [%s] and projects [%+#v]", integrationType, userID, ids))
				}
				lock.Lock()
				defer lock.Unlock()

				for _, integration := range integrations {
					integrationSettings[integration.ID] = &entities.ProjectSettingsIntegration{
						Type:     integrationType,
						ID:       integration.ID,
						Settings: integration,
					}
				}
				return nil
			})
		case entities.IntegrationTypePhoneCall:
			errGroup.Go(func() error {
				integrations, err := service.phoneCallIntegrationRepository.FetchMultiple(ctx, userID, ids)
				if err != nil {
					return stacktrace.Propagate(err, fmt.Sprintf("cannot fetch [%s] integraions for userID [%s] and projects [%+#v]", integrationType, userID, ids))
				}
				lock.Lock()
				defer lock.Unlock()

				for _, integration := range integrations {
					integrationSettings[integration.ID] = &entities.ProjectSettingsIntegration{
						Type:     integrationType,
						ID:       integration.ID,
						Settings: integration,
					}
				}
				return nil
			})
		case entities.IntegrationTypeLink:
			errGroup.Go(func() error {
				integrations, err := service.linkIntegrationRepository.FetchMultiple(ctx, userID, ids)
				if err != nil {
					return stacktrace.Propagate(err, fmt.Sprintf("cannot fetch [%s] integraions for userID [%s] and projects [%+#v]", integrationType, userID, ids))
				}
				lock.Lock()
				defer lock.Unlock()

				for _, integration := range integrations {
					integrationSettings[integration.ID] = &entities.ProjectSettingsIntegration{
						Type:     integrationType,
						ID:       integration.ID,
						Settings: integration,
					}
				}
				return nil
			})
		}
	}
	if err := errGroup.Wait(); err != nil {
		msg := fmt.Sprintf("cannot fetch project settings for integraiont groups [%s]", spew.Sdump(integrationGroups))
		return integrationSettings, service.tracer.WrapErrorSpan(span, stacktrace.Propagate(err, msg))
	}
	return integrationSettings, nil
}

func (service *ProjectSettingsService) createWithoutIntegrations(project *entities.Project) *entities.ProjectSettings {
	return &entities.ProjectSettings{
		Project:      project,
		Integrations: []*entities.ProjectSettingsIntegration{},
	}
}

func (service *ProjectSettingsService) creatSortedIntegrations(project *entities.Project, integrations map[uuid.UUID]*entities.ProjectSettingsIntegration, integrationOrder []*entities.ProjectIntegration) *entities.ProjectSettings {
	result := service.createWithoutIntegrations(project)
	for _, integration := range integrationOrder {
		if _, ok := integrations[integration.IntegrationID]; ok {
			result.Integrations = append(result.Integrations, integrations[integration.IntegrationID])
		}
	}
	return result
}

func (service *ProjectSettingsService) groupByType(integrations []*entities.ProjectIntegration) map[entities.IntegrationType][]uuid.UUID {
	result := map[entities.IntegrationType][]uuid.UUID{}
	for _, integration := range integrations {
		result[integration.Type] = append(result[integration.Type], integration.IntegrationID)
	}
	return result
}
