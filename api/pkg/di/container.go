package di

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	lemonsqueezy "github.com/NdoleStudio/lemonsqueezy-go"

	"github.com/NdoleStudio/superbutton/pkg/listeners"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"github.com/NdoleStudio/superbutton/pkg/entities"
	"github.com/NdoleStudio/superbutton/pkg/handlers"
	"github.com/NdoleStudio/superbutton/pkg/middlewares"
	"github.com/NdoleStudio/superbutton/pkg/queue"
	"github.com/NdoleStudio/superbutton/pkg/repositories"
	"github.com/NdoleStudio/superbutton/pkg/services"
	"github.com/NdoleStudio/superbutton/pkg/telemetry"
	"github.com/NdoleStudio/superbutton/pkg/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/hirosassa/zerodriver"
	"github.com/jinzhu/now"
	"github.com/palantir/stacktrace"
	"github.com/rs/zerolog"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"google.golang.org/api/option"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// Container is used to resolve services at runtime
type Container struct {
	projectID       string
	db              *gorm.DB
	app             *fiber.App
	eventDispatcher *services.EventDispatcher
	logger          telemetry.Logger
}

// NewContainer creates a new dependency injection container
func NewContainer(version string, projectID string) (container *Container) {
	// Set location to UTC
	now.DefaultConfig = &now.Config{
		TimeLocation: time.UTC,
	}

	container = &Container{
		projectID: projectID,
		logger:    logger(3).WithService(fmt.Sprintf("%T", container)),
	}

	container.InitializeTraceProvider(version, os.Getenv("GCP_PROJECT_ID"))

	container.RegisterMarketingListeners()
	container.RegisterUserListeners()

	container.RegisterUserRoutes()
	container.RegisterEventRoutes()
	container.RegisterProjectRoutes()
	container.RegisterWhatsappIntegrationRoutes()
	container.ProjectIntegrationRoutes()
	container.RegisterContentIntegrationRoutes()
	container.RegisterPhoneCallIntegrationRoutes()
	container.RegisterLinkIntegrationRoutes()

	// UnAuthenticated routes
	container.RegisterProjectSettingsRoutes()
	container.RegisterLemonsqueezyRoutes()

	// this has to be last since it registers the /* route
	container.RegisterSwaggerRoutes()

	return container
}

// RegisterSwaggerRoutes registers routes for swagger
func (container *Container) RegisterSwaggerRoutes() {
	container.logger.Debug("registering swagger routes")
	container.App().Get("/*", swagger.HandlerDefault)
}

// AuthenticatedMiddleware creates a new instance of middlewares.Authenticated
func (container *Container) AuthenticatedMiddleware() fiber.Handler {
	container.logger.Debug("creating middlewares.Authenticated")
	return middlewares.Authenticated(container.Tracer())
}

// GoogleAuthMiddlewares creates router for authenticated requests
func (container *Container) GoogleAuthMiddlewares(audience string, subject string) []fiber.Handler {
	container.logger.Debug("creating GoogleAuthMiddlewares")
	return []fiber.Handler{
		middlewares.GoogleAuth(container.Logger(), container.Tracer(), audience, subject),
		container.AuthenticatedMiddleware(),
	}
}

// FirebaseAuthMiddlewares creates router for authenticated requests
func (container *Container) FirebaseAuthMiddlewares() []fiber.Handler {
	container.logger.Debug("creating FirebaseAuthRouter")
	return []fiber.Handler{
		middlewares.FirebaseAuth(container.Logger(), container.Tracer(), container.FirebaseAuthClient()),
		container.AuthenticatedMiddleware(),
	}
}

// RegisterMarketingListeners registers the marketing handlers to events
func (container *Container) RegisterMarketingListeners() {
	container.logger.Debug(fmt.Sprintf("registering %T listeners", &listeners.MarketingListener{}))
	routes := listeners.MarketingListeners(
		container.Tracer(),
		container.Logger(),
		container.MarketingService(),
	)
	for event, listener := range routes {
		container.EventDispatcher().Subscribe(event, listener)
	}
}

// RegisterUserListeners registers the user handlers to events
func (container *Container) RegisterUserListeners() {
	container.logger.Debug(fmt.Sprintf("registering %T listeners", &listeners.UserListener{}))
	routes := listeners.UserListeners(
		container.Tracer(),
		container.Logger(),
		container.UserService(),
	)
	for event, listener := range routes {
		container.EventDispatcher().Subscribe(event, listener)
	}
}

// RegisterUserRoutes registers routes for the /users prefix
func (container *Container) RegisterUserRoutes() {
	container.logger.Debug(fmt.Sprintf("registering %T routes", &handlers.UserHandler{}))
	container.UserHandler().RegisterRoutes(container.App(), container.FirebaseAuthMiddlewares())
}

// RegisterProjectSettingsRoutes registers routes for the /project-settings prefix
func (container *Container) RegisterProjectSettingsRoutes() {
	container.logger.Debug(fmt.Sprintf("registering %T routes", &handlers.ProjectSettingsHandler{}))
	container.ProjectSettingsHandler().RegisterRoutes(container.App())
}

// RegisterLemonsqueezyRoutes registers routes for the /project-settings prefix
func (container *Container) RegisterLemonsqueezyRoutes() {
	container.logger.Debug(fmt.Sprintf("registering %T routes", &handlers.LemonsqueezyHandler{}))
	container.LemonsqueezyHandler().RegisterRoutes(container.App())
}

// RegisterProjectRoutes registers routes for the /projects prefix
func (container *Container) RegisterProjectRoutes() {
	container.logger.Debug(fmt.Sprintf("registering %T routes", &handlers.ProjectHandler{}))
	container.ProjectHandler().RegisterRoutes(container.App(), container.FirebaseAuthMiddlewares())
}

// RegisterWhatsappIntegrationRoutes registers routes for the /projects/:projectID/whatsapp-integrations prefix
func (container *Container) RegisterWhatsappIntegrationRoutes() {
	container.logger.Debug(fmt.Sprintf("registering %T routes", &handlers.WhatsappIntegrationHandler{}))
	container.WhatsappIntegrationHandler().RegisterRoutes(container.App(), container.FirebaseAuthMiddlewares())
}

// RegisterContentIntegrationRoutes registers routes for the /projects/:projectID/content-integrations prefix
func (container *Container) RegisterContentIntegrationRoutes() {
	container.logger.Debug(fmt.Sprintf("registering %T routes", &handlers.ContentIntegrationHandler{}))
	container.ContentIntegrationHandler().RegisterRoutes(container.App(), container.FirebaseAuthMiddlewares())
}

// RegisterPhoneCallIntegrationRoutes registers routes for the /projects/:projectID/phone-call-integrations prefix
func (container *Container) RegisterPhoneCallIntegrationRoutes() {
	container.logger.Debug(fmt.Sprintf("registering %T routes", &handlers.ContentIntegrationHandler{}))
	container.PhoneCallIntegrationHandler().RegisterRoutes(container.App(), container.FirebaseAuthMiddlewares())
}

// RegisterLinkIntegrationRoutes registers routes for the /projects/:projectID/link-integrations prefix
func (container *Container) RegisterLinkIntegrationRoutes() {
	container.logger.Debug(fmt.Sprintf("registering %T routes", &handlers.ContentIntegrationHandler{}))
	container.LinkIntegrationHandler().RegisterRoutes(container.App(), container.FirebaseAuthMiddlewares())
}

// ProjectIntegrationRoutes registers routes for the /projects/:projectID/integrations prefix
func (container *Container) ProjectIntegrationRoutes() {
	container.logger.Debug(fmt.Sprintf("registering %T routes", &handlers.ProjectIntegrationHandler{}))
	container.ProjectIntegrationHandler().RegisterRoutes(container.App(), container.FirebaseAuthMiddlewares())
}

// UserHandlerValidator creates a new instance of validators.UserHandlerValidator
func (container *Container) UserHandlerValidator() (validator *validators.UserHandlerValidator) {
	container.logger.Debug(fmt.Sprintf("creating %T", validator))
	return validators.NewUserHandlerValidator(
		container.Logger(),
		container.Tracer(),
	)
}

// LemonsqueezyHandlerValidator creates a new instance of validators.LemonsqueezyHandlerValidator
func (container *Container) LemonsqueezyHandlerValidator() (validator *validators.LemonsqueezyHandlerValidator) {
	container.logger.Debug(fmt.Sprintf("creating %T", validator))
	return validators.NewLemonsqueezyHandlerValidator(
		container.Logger(),
		container.Tracer(),
		container.LemonsqueezyClient(),
	)
}

// LinkIntegrationHandlerValidator creates a new instance of validators.LinkIntegrationHandlerValidator
func (container *Container) LinkIntegrationHandlerValidator() (validator *validators.LinkIntegrationHandlerValidator) {
	container.logger.Debug(fmt.Sprintf("creating %T", validator))
	return validators.NewLinkIntegrationHandlerValidator(
		container.Logger(),
		container.Tracer(),
	)
}

// ContentIntegrationHandlerValidator creates a new instance of validators.ContentIntegrationHandlerValidator
func (container *Container) ContentIntegrationHandlerValidator() (validator *validators.ContentIntegrationHandlerValidator) {
	container.logger.Debug(fmt.Sprintf("creating %T", validator))
	return validators.NewContentIntegrationHandlerValidator(
		container.Logger(),
		container.Tracer(),
	)
}

// PhoneCallIntegrationHandlerValidator creates a new instance of validators.PhoneCallIntegrationHandlerValidator
func (container *Container) PhoneCallIntegrationHandlerValidator() (validator *validators.PhoneCallIntegrationHandlerValidator) {
	container.logger.Debug(fmt.Sprintf("creating %T", validator))
	return validators.NewPhoneCallIntegrationHandlerValidator(
		container.Logger(),
		container.Tracer(),
	)
}

// WhatsappIntegrationHandlerValidator creates a new instance of validators.WhatsappIntegrationHandlerValidator
func (container *Container) WhatsappIntegrationHandlerValidator() (validator *validators.WhatsappIntegrationHandlerValidator) {
	container.logger.Debug(fmt.Sprintf("creating %T", validator))
	return validators.NewWhatsappIntegrationHandlerValidator(
		container.Logger(),
		container.Tracer(),
	)
}

// ProjectHandlerValidator creates a new instance of validators.ProjectHandlerValidator
func (container *Container) ProjectHandlerValidator() (validator *validators.ProjectHandlerValidator) {
	container.logger.Debug(fmt.Sprintf("creating %T", validator))
	return validators.NewProjectHandlerValidator(
		container.Logger(),
		container.Tracer(),
	)
}

// UserHandler creates a new instance of handlers.UserHandler
func (container *Container) UserHandler() (handler *handlers.UserHandler) {
	container.logger.Debug(fmt.Sprintf("creating %T", handler))
	return handlers.NewUserHandler(
		container.Logger(),
		container.Tracer(),
		container.UserHandlerValidator(),
		container.UserService(),
	)
}

// WhatsappIntegrationHandler creates a new instance of handlers.WhatsappIntegrationHandler
func (container *Container) WhatsappIntegrationHandler() (handler *handlers.WhatsappIntegrationHandler) {
	container.logger.Debug(fmt.Sprintf("creating %T", handler))
	return handlers.NewWhatsappIntegrationHandler(
		container.Logger(),
		container.Tracer(),
		container.WhatsappIntegrationHandlerValidator(),
		container.WhatsappIntegrationService(),
	)
}

// ContentIntegrationHandler creates a new instance of handlers.ContentIntegrationHandler
func (container *Container) ContentIntegrationHandler() (handler *handlers.ContentIntegrationHandler) {
	container.logger.Debug(fmt.Sprintf("creating %T", handler))
	return handlers.NewContentIntegrationHandler(
		container.Logger(),
		container.Tracer(),
		container.ContentIntegrationHandlerValidator(),
		container.ContentIntegrationService(),
	)
}

// ProjectIntegrationHandler creates a new instance of handlers.ProjectIntegrationHandler
func (container *Container) ProjectIntegrationHandler() (handler *handlers.ProjectIntegrationHandler) {
	container.logger.Debug(fmt.Sprintf("creating %T", handler))
	return handlers.NewIntegrationHandler(
		container.Logger(),
		container.Tracer(),
		container.ProjectIntegrationService(),
	)
}

// PhoneCallIntegrationHandler creates a new instance of handlers.PhoneCallIntegrationHandler
func (container *Container) PhoneCallIntegrationHandler() (handler *handlers.PhoneCallIntegrationHandler) {
	container.logger.Debug(fmt.Sprintf("creating %T", handler))
	return handlers.NewPhoneCallIntegrationHandler(
		container.Logger(),
		container.Tracer(),
		container.PhoneCallIntegrationHandlerValidator(),
		container.PhoneCallIntegrationService(),
	)
}

// LinkIntegrationHandler creates a new instance of handlers.LinkIntegrationHandler
func (container *Container) LinkIntegrationHandler() (handler *handlers.LinkIntegrationHandler) {
	container.logger.Debug(fmt.Sprintf("creating %T", handler))
	return handlers.NewLinkIntegrationHandler(
		container.Logger(),
		container.Tracer(),
		container.LinkIntegrationHandlerValidator(),
		container.LinkIntegrationService(),
	)
}

// ProjectHandler creates a new instance of handlers.ProjectHandler
func (container *Container) ProjectHandler() (handler *handlers.ProjectHandler) {
	container.logger.Debug(fmt.Sprintf("creating %T", handler))
	return handlers.NewProjectHandler(
		container.Logger(),
		container.Tracer(),
		container.ProjectHandlerValidator(),
		container.ProjectService(),
	)
}

// ProjectSettingsHandler creates a new instance of handlers.ProjectSettingsHandler
func (container *Container) ProjectSettingsHandler() (handler *handlers.ProjectSettingsHandler) {
	container.logger.Debug(fmt.Sprintf("creating %T", handler))
	return handlers.NewProjectSettingsHandler(
		container.Logger(),
		container.Tracer(),
		container.ProjectSettingService(),
	)
}

// MarketingService creates a new instance of services.MarketingService
func (container *Container) MarketingService() (service *services.MarketingService) {
	container.logger.Debug(fmt.Sprintf("creating %T", service))
	return services.NewMarketingService(
		container.Logger(),
		container.Tracer(),
		container.FirebaseAuthClient(),
		os.Getenv("SENDGRID_API_KEY"),
		os.Getenv("SENDGRID_LIST_ID"),
	)
}

// LemonsqueezyService creates a new instance of services.LemonsqueezyService
func (container *Container) LemonsqueezyService() (service *services.LemonsqueezyService) {
	container.logger.Debug(fmt.Sprintf("creating %T", service))
	return services.NewLemonsqueezyService(
		container.Logger(),
		container.Tracer(),
		container.UserRepository(),
		container.EventDispatcher(),
	)
}

// UserService creates a new instance of services.UserService
func (container *Container) UserService() (service *services.UserService) {
	container.logger.Debug(fmt.Sprintf("creating %T", service))
	return services.NewUserService(
		container.Logger(),
		container.Tracer(),
		container.EventDispatcher(),
		container.UserRepository(),
		container.LemonsqueezyClient(),
	)
}

// ProjectIntegrationService creates a new instance of services.ProjectIntegrationService
func (container *Container) ProjectIntegrationService() (service *services.ProjectIntegrationService) {
	container.logger.Debug(fmt.Sprintf("creating %T", service))
	return services.NewProjectIntegrationService(
		container.Logger(),
		container.Tracer(),
		container.EventDispatcher(),
		container.ProjectIntegrationRepository(),
	)
}

// ProjectSettingService creates a new instance of services.ProjectSettingsService
func (container *Container) ProjectSettingService() (service *services.ProjectSettingsService) {
	container.logger.Debug(fmt.Sprintf("creating %T", service))
	return services.NewProjectSettingsService(
		container.Logger(),
		container.Tracer(),
		container.ProjectRepository(),
		container.WhatsappIntegrationRepository(),
		container.ContentIntegrationRepository(),
		container.ProjectIntegrationRepository(),
		container.PhoneCallIntegrationRepository(),
		container.LinkIntegrationRepository(),
	)
}

// WhatsappIntegrationService creates a new instance of services.WhatsappIntegrationService
func (container *Container) WhatsappIntegrationService() (service *services.WhatsappIntegrationService) {
	container.logger.Debug(fmt.Sprintf("creating %T", service))
	return services.NewWhatsappIntegrationService(
		container.Logger(),
		container.Tracer(),
		container.EventDispatcher(),
		container.ProjectRepository(),
		container.WhatsappIntegrationRepository(),
	)
}

// ContentIntegrationService creates a new instance of services.ContentIntegrationService
func (container *Container) ContentIntegrationService() (service *services.ContentIntegrationService) {
	container.logger.Debug(fmt.Sprintf("creating %T", service))
	return services.NewContentIntegrationService(
		container.Logger(),
		container.Tracer(),
		container.EventDispatcher(),
		container.ProjectRepository(),
		container.ContentIntegrationRepository(),
	)
}

// PhoneCallIntegrationService creates a new instance of services.PhoneCallIntegrationService
func (container *Container) PhoneCallIntegrationService() (service *services.PhoneCallIntegrationService) {
	container.logger.Debug(fmt.Sprintf("creating %T", service))
	return services.NewPhoneCallIntegrationService(
		container.Logger(),
		container.Tracer(),
		container.EventDispatcher(),
		container.ProjectRepository(),
		container.PhoneCallIntegrationRepository(),
	)
}

// LinkIntegrationService creates a new instance of services.LinkIntegrationService
func (container *Container) LinkIntegrationService() (service *services.LinkIntegrationService) {
	container.logger.Debug(fmt.Sprintf("creating %T", service))
	return services.NewLinkIntegrationService(
		container.Logger(),
		container.Tracer(),
		container.EventDispatcher(),
		container.ProjectRepository(),
		container.LinkIntegrationRepository(),
	)
}

// ProjectService creates a new instance of services.ProjectService
func (container *Container) ProjectService() (service *services.ProjectService) {
	container.logger.Debug(fmt.Sprintf("creating %T", service))
	return services.NewProjectService(
		container.Logger(),
		container.Tracer(),
		container.EventDispatcher(),
		container.ProjectRepository(),
	)
}

// EventDispatcher creates a new instance of services.EventDispatcher
func (container *Container) EventDispatcher() (dispatcher *services.EventDispatcher) {
	if container.eventDispatcher != nil {
		return container.eventDispatcher
	}

	container.logger.Debug(fmt.Sprintf("creating %T", dispatcher))
	dispatcher = services.NewEventDispatcher(
		container.Logger(),
		container.Tracer(),
		container.EventRepository(),
		container.EventsQueue(),
		os.Getenv("QUEUE_URL_EVENTS"),
	)

	container.eventDispatcher = dispatcher
	return dispatcher
}

// EventsQueue creates a new instance of services.PushQueue
func (container *Container) EventsQueue() queue.Client {
	container.logger.Debug("creating queue.Client")

	return queue.NewGooglePushQueue(
		container.Logger(),
		container.Tracer(),
		container.CloudTasksClient(),
		os.Getenv("QUEUE_NAME_EVENTS"),
		os.Getenv("QUEUE_AUTH_EMAIL"),
	)
}

// CloudTasksClient creates a new instance of cloudtasks.Client
func (container *Container) CloudTasksClient() (client *cloudtasks.Client) {
	container.logger.Debug(fmt.Sprintf("creating %T", client))

	client, err := cloudtasks.NewClient(
		context.Background(),
		option.WithCredentialsJSON(container.FirebaseCredentials()),
	)
	if err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, "cannot initialize cloud tasks client"))
	}

	return client
}

// EventRepository creates a new instance of repositories.EventRepository
func (container *Container) EventRepository() (repository repositories.EventRepository) {
	container.logger.Debug("creating GORM repositories.EventRepository")
	return repositories.NewGormEventRepository(
		container.Logger(),
		container.Tracer(),
		container.DB(),
	)
}

// RegisterEventRoutes registers routes for the /events prefix
func (container *Container) RegisterEventRoutes() {
	container.logger.Debug(fmt.Sprintf("registering %T routes", &handlers.EventsHandler{}))
	container.
		EventsHandler().
		RegisterRoutes(
			container.App(),
			container.GoogleAuthMiddlewares(
				os.Getenv("QUEUE_URL_EVENTS"),
				os.Getenv("QUEUE_AUTH_SUBJECT"),
			),
		)
}

// LemonsqueezyHandler creates a new instance of handlers.LemonsqueezyHandler
func (container *Container) LemonsqueezyHandler() (handler *handlers.LemonsqueezyHandler) {
	container.logger.Debug(fmt.Sprintf("creating %T", handler))

	return handlers.NewLemonsqueezyHandlerHandler(
		container.Logger(),
		container.Tracer(),
		container.LemonsqueezyService(),
		container.LemonsqueezyHandlerValidator(),
	)
}

// EventsHandler creates a new instance of handlers.EventsHandler
func (container *Container) EventsHandler() (handler *handlers.EventsHandler) {
	container.logger.Debug(fmt.Sprintf("creating %T", handler))

	return handlers.NewEventsHandler(
		container.Logger(),
		container.Tracer(),
		container.EventDispatcher(),
	)
}

// App creates a new instance of fiber.App
func (container *Container) App() (app *fiber.App) {
	if container.app != nil {
		return container.app
	}

	container.logger.Debug(fmt.Sprintf("creating %T", app))

	app = fiber.New()

	if isLocal() {
		app.Use(fiberLogger.New())
	}
	app.Use(
		middlewares.OtelTraceContext(
			container.Tracer(),
			container.Logger(),
			"X-Cloud-Trace-Context",
			os.Getenv("GCP_PROJECT_ID"),
		),
	)
	app.Use(cors.New())

	container.app = app

	return app
}

// InitializeOtelResources initializes open telemetry resources
func (container *Container) InitializeOtelResources(version string, namespace string) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(namespace),
		semconv.ServiceNamespaceKey.String(namespace),
		semconv.ServiceVersionKey.String(version),
		semconv.ServiceInstanceIDKey.String(hostName()),
		attribute.String("service.environment", os.Getenv("ENV")),
	)
}

// UserRepository registers a new instance of repositories.UserRepository
func (container *Container) UserRepository() repositories.UserRepository {
	container.logger.Debug("creating GORM repositories.UserRepository")
	return repositories.NewGormUserRepository(
		container.Logger(),
		container.Tracer(),
		container.DB(),
	)
}

// ProjectRepository registers a new instance of repositories.ProjectRepository
func (container *Container) ProjectRepository() repositories.ProjectRepository {
	container.logger.Debug("creating GORM repositories.ProjectRepository")
	return repositories.NewGormProjectRepository(
		container.Logger(),
		container.Tracer(),
		container.DB(),
	)
}

// WhatsappIntegrationRepository registers a new instance of repositories.WhatsappIntegrationRepository
func (container *Container) WhatsappIntegrationRepository() repositories.WhatsappIntegrationRepository {
	container.logger.Debug("creating GORM repositories.WhatsappIntegrationRepository")
	return repositories.NewGormWhatsappIntegrationRepository(
		container.Logger(),
		container.Tracer(),
		container.DB(),
	)
}

// ContentIntegrationRepository registers a new instance of repositories.ContentIntegrationRepository
func (container *Container) ContentIntegrationRepository() repositories.ContentIntegrationRepository {
	container.logger.Debug("creating GORM repositories.ContentIntegrationRepository")
	return repositories.NewGormContentIntegrationRepository(
		container.Logger(),
		container.Tracer(),
		container.DB(),
	)
}

// PhoneCallIntegrationRepository registers a new instance of repositories.PhoneCallIntegrationRepository
func (container *Container) PhoneCallIntegrationRepository() repositories.PhoneCallIntegrationRepository {
	container.logger.Debug("creating GORM repositories.PhoneCallIntegrationRepository")
	return repositories.NewGormPhoneCallIntegrationRepository(
		container.Logger(),
		container.Tracer(),
		container.DB(),
	)
}

// LinkIntegrationRepository registers a new instance of repositories.LinkIntegrationRepository
func (container *Container) LinkIntegrationRepository() repositories.LinkIntegrationRepository {
	container.logger.Debug("creating GORM repositories.LinkIntegrationRepository")
	return repositories.NewGormLinkIntegrationRepository(
		container.Logger(),
		container.Tracer(),
		container.DB(),
	)
}

// ProjectIntegrationRepository registers a new instance of repositories.ProjectIntegrationRepository
func (container *Container) ProjectIntegrationRepository() repositories.ProjectIntegrationRepository {
	container.logger.Debug("creating GORM repositories.UserRepository")
	return repositories.NewGormProjectIntegrationRepository(
		container.Logger(),
		container.Tracer(),
		container.DB(),
	)
}

// Logger creates a new instance of telemetry.Logger
func (container *Container) Logger(skipFrameCount ...int) telemetry.Logger {
	container.logger.Debug("creating telemetry.Logger")
	if len(skipFrameCount) > 0 {
		return logger(skipFrameCount[0])
	}
	return logger(3)
}

// GormLogger creates a new instance of gormLogger.Interface
func (container *Container) GormLogger() gormLogger.Interface {
	container.logger.Debug("creating gormLogger.Interface")
	return telemetry.NewGormLogger(
		container.Tracer(),
		container.Logger(6),
	)
}

// DB creates an instance of gorm.DB if it has not been created already
func (container *Container) DB() (db *gorm.DB) {
	if container.db != nil {
		return container.db
	}

	container.logger.Debug(fmt.Sprintf("creating %T", db))

	config := &gorm.Config{}
	if isLocal() {
		config = &gorm.Config{Logger: container.GormLogger()}
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_DSN")), config)
	if err != nil {
		container.logger.Fatal(err)
	}
	container.db = db

	container.logger.Debug(fmt.Sprintf("Running migrations for %T", db))

	if err = db.AutoMigrate(&entities.User{}); err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, fmt.Sprintf("cannot migrate %T", &entities.User{})))
	}
	if err = db.AutoMigrate(&repositories.GormEvent{}); err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, fmt.Sprintf("cannot migrate %T", &repositories.GormEvent{})))
	}
	if err = db.AutoMigrate(&entities.Project{}); err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, fmt.Sprintf("cannot migrate %T", &entities.Project{})))
	}
	if err = db.AutoMigrate(&entities.ProjectIntegration{}); err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, fmt.Sprintf("cannot migrate %T", &entities.ProjectIntegration{})))
	}
	if err = db.AutoMigrate(&entities.WhatsappIntegration{}); err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, fmt.Sprintf("cannot migrate %T", &entities.WhatsappIntegration{})))
	}
	if err = db.AutoMigrate(&entities.ContentIntegration{}); err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, fmt.Sprintf("cannot migrate %T", &entities.ContentIntegration{})))
	}
	if err = db.AutoMigrate(&entities.PhoneCallIntegration{}); err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, fmt.Sprintf("cannot migrate %T", &entities.PhoneCallIntegration{})))
	}
	if err = db.AutoMigrate(&entities.LinkIntegration{}); err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, fmt.Sprintf("cannot migrate %T", &entities.LinkIntegration{})))
	}

	return container.db
}

// FirebaseCredentials returns firebase credentials as bytes.
func (container *Container) FirebaseCredentials() []byte {
	container.logger.Debug("creating firebase credentials")
	return []byte(os.Getenv("FIREBASE_CREDENTIALS"))
}

// Tracer creates a new instance of telemetry.Tracer
func (container *Container) Tracer() (t telemetry.Tracer) {
	container.logger.Debug("creating telemetry.Tracer")
	return telemetry.NewOtelLogger(
		container.projectID,
		container.Logger(),
	)
}

// FirebaseApp creates a new instance of firebase.App
func (container *Container) FirebaseApp() (app *firebase.App) {
	container.logger.Debug(fmt.Sprintf("creating %T", app))
	app, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsJSON(container.FirebaseCredentials()))
	if err != nil {
		msg := "cannot initialize firebase application"
		container.logger.Fatal(stacktrace.Propagate(err, msg))
	}
	return app
}

// LemonsqueezyClient creates a new instance of lemonsqueezy.Client
func (container *Container) LemonsqueezyClient() (client *lemonsqueezy.Client) {
	container.logger.Debug(fmt.Sprintf("creating %T", client))
	return lemonsqueezy.New(
		lemonsqueezy.WithAPIKey(os.Getenv("LEMONSQUEEZY_API_KEY")),
		lemonsqueezy.WithSigningSecret(os.Getenv("LEMONSQUEEZY_SIGNING_SECRET")),
	)
}

// FirebaseAuthClient creates a new instance of auth.Client
func (container *Container) FirebaseAuthClient() (client *auth.Client) {
	container.logger.Debug(fmt.Sprintf("creating %T", client))
	authClient, err := container.FirebaseApp().Auth(context.Background())
	if err != nil {
		msg := "cannot initialize firebase auth client"
		container.logger.Fatal(stacktrace.Propagate(err, msg))
	}
	return authClient
}

// InitializeTraceProvider initializes the open telemetry trace provider
func (container *Container) InitializeTraceProvider(version string, namespace string) func() {
	if isLocal() {
		return container.initializeUptraceProvider(version, namespace)
	}
	return container.initializeGoogleTraceProvider(version, namespace)
}

func (container *Container) initializeGoogleTraceProvider(version string, namespace string) func() {
	container.logger.Debug("initializing google trace provider")

	exporter, err := cloudtrace.New(cloudtrace.WithProjectID(os.Getenv("GCP_PROJECT_ID")))
	if err != nil {
		container.logger.Fatal(stacktrace.Propagate(err, "cannot create cloud trace exporter"))
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(container.InitializeOtelResources(version, namespace)),
	)

	otel.SetTracerProvider(tp)

	return func() {
		_ = exporter.Shutdown(context.Background())
	}
}

func (container *Container) initializeUptraceProvider(version string, namespace string) (flush func()) {
	container.logger.Debug("initializing uptrace provider")
	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(os.Getenv("UPTRACE_DSN")),
		uptrace.WithServiceName(namespace),
		uptrace.WithServiceVersion(version),
	)

	// Send buffered spans and free resources.
	return func() {
		err := uptrace.Shutdown(context.Background())
		if err != nil {
			container.logger.Error(err)
		}
	}
}

func logger(skipFrameCount int) telemetry.Logger {
	fields := map[string]string{
		"pid":      strconv.Itoa(os.Getpid()),
		"hostname": hostName(),
	}

	return telemetry.NewZerologLogger(
		os.Getenv("GCP_PROJECT_ID"),
		fields,
		logDriver(skipFrameCount),
		nil,
	)
}

func logDriver(skipFrameCount int) *zerodriver.Logger {
	if isLocal() {
		return consoleLogger(skipFrameCount)
	}
	return jsonLogger(skipFrameCount)
}

func jsonLogger(skipFrameCount int) *zerodriver.Logger {
	logLevel := zerolog.DebugLevel
	zerolog.SetGlobalLevel(logLevel)

	// See: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogSeverity
	logLevelSeverity := map[zerolog.Level]string{
		zerolog.TraceLevel: "DEFAULT",
		zerolog.DebugLevel: "DEBUG",
		zerolog.InfoLevel:  "INFO",
		zerolog.WarnLevel:  "WARNING",
		zerolog.ErrorLevel: "ERROR",
		zerolog.PanicLevel: "CRITICAL",
		zerolog.FatalLevel: "CRITICAL",
	}

	zerolog.LevelFieldName = "severity"
	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		return logLevelSeverity[l]
	}
	zerolog.TimestampFieldName = "time"
	zerolog.TimeFieldFormat = time.RFC3339Nano

	zl := zerolog.New(os.Stderr).With().Timestamp().CallerWithSkipFrameCount(skipFrameCount).Logger()
	return &zerodriver.Logger{Logger: &zl}
}

func hostName() string {
	h, err := os.Hostname()
	if err != nil {
		h = strconv.Itoa(os.Getpid())
	}
	return h
}

func consoleLogger(skipFrameCount int) *zerodriver.Logger {
	l := zerolog.New(
		zerolog.ConsoleWriter{
			Out: os.Stderr,
		}).With().Timestamp().CallerWithSkipFrameCount(skipFrameCount).Logger()
	return &zerodriver.Logger{
		Logger: &l,
	}
}

func isLocal() bool {
	return os.Getenv("ENV") == "local"
}
