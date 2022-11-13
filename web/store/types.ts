import { AxiosError } from 'axios'
import { EntitiesUser } from '~/store/backend'

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

export type AppData = {
  url: string
  name: string
  environment: string
  documentationURL: string
  githubURL: string
}

export interface State {
  user: EntitiesUser | null
  authUser: AuthUser | null
  axiosError: AxiosError | null
  nextRoute: string | null
  authStateChanged: boolean
  notification: Notification | null
}
