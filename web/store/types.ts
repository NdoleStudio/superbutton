import { AxiosError } from 'axios'
import {
  EntitiesProject,
  EntitiesUser,
  RequestsProjectUpdateRequest,
  RequestsWhatsappIntegrationCreateRequest,
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
