import { AxiosError } from 'axios'
import {
  EntitiesProject,
  EntitiesUser,
  RequestsContentIntegrationCreateRequest,
  RequestsContentIntegrationUpdateRequest,
  RequestsLinkIntegrationCreateRequest,
  RequestsLinkIntegrationUpdateRequest,
  RequestsPhoneCallIntegrationCreateRequest,
  RequestsPhoneCallIntegrationUpdateRequest,
  RequestsProjectIntegrationsUpdateRequest,
  RequestsProjectUpdateRequest,
  RequestsWhatsappIntegrationCreateRequest,
  RequestsWhatsappIntegrationUpdateRequest,
} from '~/store/backend'
import { ErrorMessagesSerialized } from '~/plugins/errors'

export type AuthUser = {
  uid: string
  email: string
  photoURL?: string
  displayName: string | null
}

type NotificationType = 'error' | 'success' | 'info'

export interface Notification {
  message: string
  timeout: number
  active: boolean
  type: NotificationType
}

export interface NotificationRequest {
  message: string
  type: NotificationType
}

export interface UpdateProjectRequest extends RequestsProjectUpdateRequest {
  projectId: string
}
export interface AddWhatsappIntegrationRequest
  extends RequestsWhatsappIntegrationCreateRequest {
  projectId: string
}

export interface ProjectIntegrationIdRequest {
  projectId: string
  integrationId: string
}

export interface AddContentIntegrationRequest
  extends RequestsContentIntegrationCreateRequest {
  projectId: string
}

export interface AddPhoneCallIntegrationRequest
  extends RequestsPhoneCallIntegrationCreateRequest {
  projectId: string
}

export interface UpdatePhoneCallIntegrationRequest
  extends RequestsPhoneCallIntegrationUpdateRequest,
    ProjectIntegrationIdRequest {}

export interface UpdateContentIntegrationRequest
  extends RequestsContentIntegrationUpdateRequest,
    ProjectIntegrationIdRequest {}

export interface UpdateWhatsappIntegrationRequest
  extends RequestsWhatsappIntegrationUpdateRequest,
    ProjectIntegrationIdRequest {}

export interface AddLinkIntegrationRequest
  extends RequestsLinkIntegrationCreateRequest {
  projectId: string
}

export interface UpdateLinkIntegrationRequest
  extends RequestsLinkIntegrationUpdateRequest,
    ProjectIntegrationIdRequest {}

export interface UpdateProjectIntegrationsRequest
  extends RequestsProjectIntegrationsUpdateRequest {
  projectId: string
}

export type AppData = {
  url: string
  name: string
  environment: string
  documentationURL: string
  githubURL: string
}

export interface State {
  projects: Array<EntitiesProject>
  creatingProject: boolean
  activeProjectId: string | null
  errorMessages: ErrorMessagesSerialized
  user: EntitiesUser | null
  authUser: AuthUser | null
  axiosError: AxiosError | null
  nextRoute: string | null
  authStateChanged: boolean
  notification: Notification | null
}
