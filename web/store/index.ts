import { GetterTree, ActionTree, MutationTree, ActionContext } from 'vuex'
import { AxiosError, AxiosResponse } from 'axios'
import {
  AddContentIntegrationRequest,
  AddLinkIntegrationRequest,
  AddPhoneCallIntegrationRequest,
  AddWhatsappIntegrationRequest,
  AppData,
  AuthUser,
  NotificationRequest,
  ProjectIntegrationIdRequest,
  State,
  UpdateContentIntegrationRequest,
  UpdateLinkIntegrationRequest,
  UpdatePhoneCallIntegrationRequest,
  UpdateProjectIntegrationsRequest,
  UpdateProjectRequest,
  UpdateWhatsappIntegrationRequest,
} from '~/store/types'
import {
  EntitiesContentIntegration,
  EntitiesLinkIntegration,
  EntitiesPhoneCallIntegration,
  EntitiesProject,
  EntitiesProjectIntegration,
  EntitiesProjectSettings,
  EntitiesUser,
  EntitiesWhatsappIntegration,
  ResponsesNoContent,
  ResponsesOkArrayEntitiesProject,
  ResponsesOkArrayEntitiesProjectIntegration,
  ResponsesOkEntitiesContentIntegration,
  ResponsesOkEntitiesLinkIntegration,
  ResponsesOkEntitiesPhoneCallIntegration,
  ResponsesOkEntitiesProject,
  ResponsesOkEntitiesProjectSettings,
  ResponsesOkEntitiesUser,
  ResponsesOkEntitiesWhatsappIntegration,
} from '~/store/backend'
import axios, { setAuthToken } from '~/plugins/axios'
import {
  ErrorMessages,
  ErrorMessagesSerialized,
  getErrorMessages,
} from '~/plugins/errors'

export const state = (): State => ({
  authUser: null,
  projects: [],
  errorMessages: {},
  creatingProject: false,
  nextRoute: null,
  activeProjectId: null,
  authStateChanged: false,
  notification: null,
  axiosError: null,
  user: null,
})

export type RootState = ReturnType<typeof state>

export const getters: GetterTree<RootState, RootState> = {
  authUser: (state) => state.authUser,
  authStateChanged: (state) => state.authStateChanged,
  notification: (state) => state.notification,
  projects: (state) => state.projects,
  creatingProject: (state) => state.creatingProject,
  hasProjects: (state) => state.projects.length > 0,
  activeProject: (state): EntitiesProject | null => {
    const project = state.projects.find((project) => {
      return project.id === state.activeProjectId
    })
    if (project) {
      return project
    }
    return null
  },
  activeProjectId: (state) => {
    const project = state.projects.find((project) => {
      return project.id === state.activeProjectId
    })
    if (project) {
      return project.id
    }
    if (state.projects.length) {
      return state.projects[0].id
    }
    return '0'
  },
  errorMessages: (state) =>
    ErrorMessages.fromObject<string>(state.errorMessages),
  app(): AppData {
    let url = process.env.APP_URL as string
    if (url.length > 0 && url[url.length - 1] === '/') {
      url = url.substring(0, url.length - 1)
    }
    return {
      url,
      environment: process.env.APP_ENV as string,
      documentationURL: process.env.APP_DOCUMENTATION_URL as string,
      githubURL: process.env.APP_GITHUB_URL as string,
      name: process.env.APP_NAME as string,
    }
  },
}

export const mutations: MutationTree<RootState> = {
  setAuthUser(state: RootState, payload: AuthUser | null) {
    state.authUser = payload
    state.authStateChanged = true
  },

  setUser(state: RootState, payload: EntitiesUser | null) {
    state.user = payload
  },

  setProjects(state: RootState, payload: Array<EntitiesProject>) {
    state.projects = payload
  },

  setCreatingProject(state: RootState, payload: boolean) {
    state.creatingProject = payload
  },

  setActiveProjectId(state: RootState, payload: string) {
    state.activeProjectId = payload
  },

  setErrorMessages(state: RootState, payload: ErrorMessagesSerialized) {
    state.errorMessages = payload
  },

  clearErrorMessages(state: RootState) {
    state.errorMessages = {}
  },

  setNextRoute(state: RootState, payload: string | null) {
    state.nextRoute = payload
  },

  setNotification(state: State, notification: NotificationRequest) {
    state.notification = {
      ...state.notification,
      active: true,
      message: notification.message,
      type: notification.type,
      timeout: Math.floor(Math.random() * 100) + 3000, // Reset the timeout
    }
  },

  disableNotification(state: State) {
    if (state.notification) {
      state.notification.active = false
    }
  },
}

export const actions: ActionTree<RootState, RootState> = {
  async onAuthStateChanged(
    context: ActionContext<RootState, RootState>,
    { authUser }
  ) {
    if (authUser == null) {
      context.commit('setAuthUser', null)
      context.commit('setUser', null)
      return
    }
    setAuthToken(await authUser.getIdToken())
    const { uid, email, photoURL, displayName } = authUser
    await context.commit('setAuthUser', { uid, email, photoURL, displayName })
  },

  setAuthUser: (
    context: ActionContext<RootState, RootState>,
    authUser: AuthUser | null
  ) => {
    context.commit('setAuthUser', authUser)
  },

  async loadUser(context: ActionContext<RootState, RootState>) {
    const response = await axios.get<ResponsesOkEntitiesUser>('/v1/users/me')
    context.commit('setUser', response.data.data)
  },

  async loadProjects(context: ActionContext<RootState, RootState>) {
    const response = await axios.get<ResponsesOkArrayEntitiesProject>(
      '/v1/projects'
    )

    const projects = response.data.data

    const activeProject = context.state.projects.find((project) => {
      return project.id === context.state.activeProjectId
    })
    if (activeProject === undefined && projects.length) {
      context.commit('setActiveProjectId', projects[0].id)
    }

    await context.commit('setProjects', projects)
  },

  async setActiveProjectId(
    context: ActionContext<RootState, RootState>,
    projectId: string
  ) {
    await context.commit('setActiveProjectId', projectId)
  },

  createProject(
    context: ActionContext<RootState, RootState>,
    { name, website }
  ) {
    return new Promise<EntitiesProject>((resolve, reject) => {
      context.commit('setCreatingProject', true)
      context.commit('clearErrorMessages')

      axios
        .post<ResponsesOkEntitiesProject>('/v1/projects', { name, website })
        .then(async (response: AxiosResponse<ResponsesOkEntitiesProject>) => {
          await Promise.all([
            context.dispatch('addNotification', {
              message: response.data.message ?? 'Project created successfully',
              type: 'success',
            }),
            context.dispatch('loadProjects'),
            context.commit('setCreatingProject', false),
          ])
          context.commit('setActiveProjectId', response.data.data.id)
          resolve(response.data.data)
        })
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.commit('setErrorMessages', getErrorMessages(error)),
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Validation errors while creating project',
              type: 'error',
            }),
            context.commit('setCreatingProject', false),
          ])
          reject(error)
        })
    })
  },

  updateProject(
    context: ActionContext<RootState, RootState>,
    payload: UpdateProjectRequest
  ) {
    return new Promise<EntitiesProject>((resolve, reject) => {
      context.commit('clearErrorMessages')
      axios
        .put<ResponsesOkEntitiesProject>(
          `/v1/projects/${payload.projectId}`,
          payload
        )
        .then(async (response: AxiosResponse<ResponsesOkEntitiesProject>) => {
          await Promise.all([
            context.dispatch('addNotification', {
              message: response.data.message ?? 'Project updated successfully',
              type: 'success',
            }),
            context.dispatch('loadProjects'),
          ])
          resolve(response.data.data)
        })
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.commit('setErrorMessages', getErrorMessages(error)),
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Validation errors while updating project',
              type: 'error',
            }),
          ])
          reject(error)
        })
    })
  },
  deleteProject(
    context: ActionContext<RootState, RootState>,
    projectId: string
  ) {
    return new Promise<boolean>((resolve, reject) => {
      axios
        .delete<ResponsesNoContent>(`/v1/projects/${projectId}`)
        .then(async (response: AxiosResponse<ResponsesNoContent>) => {
          await Promise.all([
            context.dispatch('addNotification', {
              message: response.data.message ?? 'Project deleted successfully',
              type: 'success',
            }),
            context.dispatch('loadProjects'),
          ])
          resolve(true)
        })
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Validation errors while deleting project',
              type: 'error',
            }),
          ])
          reject(error)
        })
    })
  },

  addWhatsappIntegration(
    context: ActionContext<RootState, RootState>,
    payload: AddWhatsappIntegrationRequest
  ) {
    return new Promise<EntitiesWhatsappIntegration>((resolve, reject) => {
      context.commit('clearErrorMessages')
      axios
        .post<ResponsesOkEntitiesWhatsappIntegration>(
          `/v1/projects/${payload.projectId}/whatsapp-integrations`,
          payload
        )
        .then(
          async (
            response: AxiosResponse<ResponsesOkEntitiesWhatsappIntegration>
          ) => {
            await Promise.all([
              context.dispatch('addNotification', {
                message:
                  response.data.message ??
                  'Whatsapp integration added successfully',
                type: 'success',
              }),
            ])
            resolve(response.data.data)
          }
        )
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.commit('setErrorMessages', getErrorMessages(error)),
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Validation errors while adding whatsapp integration',
              type: 'error',
            }),
          ])
          reject(error)
        })
    })
  },

  getProjectSettings(
    context: ActionContext<RootState, RootState>,
    projectId: string
  ) {
    return new Promise<EntitiesProjectSettings>((resolve, reject) => {
      axios
        .get<ResponsesOkEntitiesProjectSettings>(
          `/v1/settings/${context.state.authUser?.uid}/projects/${projectId}`
        )
        .then((response: AxiosResponse<ResponsesOkEntitiesProjectSettings>) => {
          resolve(response.data.data)
        })
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Error while fetching project settings',
              type: 'error',
            }),
          ])
          reject(error)
        })
    })
  },

  getWhatsappIntegration(
    context: ActionContext<RootState, RootState>,
    payload: ProjectIntegrationIdRequest
  ) {
    return new Promise<EntitiesWhatsappIntegration>((resolve, reject) => {
      axios
        .get<ResponsesOkEntitiesWhatsappIntegration>(
          `/v1/projects/${payload.projectId}/whatsapp-integrations/${payload.integrationId}`
        )
        .then(
          (response: AxiosResponse<ResponsesOkEntitiesWhatsappIntegration>) => {
            resolve(response.data.data)
          }
        )
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Error while fetching whatsapp integration integration',
              type: 'error',
            }),
          ])
          reject(error)
        })
    })
  },

  updateWhatsappIntegration(
    context: ActionContext<RootState, RootState>,
    payload: UpdateWhatsappIntegrationRequest
  ) {
    return new Promise<EntitiesWhatsappIntegration>((resolve, reject) => {
      context.commit('clearErrorMessages')
      axios
        .put<ResponsesOkEntitiesWhatsappIntegration>(
          `/v1/projects/${payload.projectId}/whatsapp-integrations/${payload.integrationId}`,
          payload
        )
        .then(
          async (
            response: AxiosResponse<ResponsesOkEntitiesWhatsappIntegration>
          ) => {
            await Promise.all([
              context.dispatch('addNotification', {
                message:
                  response.data.message ??
                  'Whatsapp integration updated successfully',
                type: 'success',
              }),
            ])
            resolve(response.data.data)
          }
        )
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.commit('setErrorMessages', getErrorMessages(error)),
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Validation errors while updating whatsapp integration',
              type: 'error',
            }),
          ])
          reject(error)
        })
    })
  },

  deleteWhatsappIntegration(
    context: ActionContext<RootState, RootState>,
    payload: UpdateWhatsappIntegrationRequest
  ) {
    return new Promise<boolean>((resolve, reject) => {
      axios
        .delete<ResponsesNoContent>(
          `/v1/projects/${payload.projectId}/whatsapp-integrations/${payload.integrationId}`
        )
        .then(async (response: AxiosResponse<ResponsesNoContent>) => {
          await Promise.all([
            context.dispatch('addNotification', {
              message:
                response.data.message ??
                'Whatsapp integration deleted successfully',
              type: 'success',
            }),
          ])
          resolve(true)
        })
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Validation errors while deleting whatsapp integration',
              type: 'error',
            }),
          ])
          reject(error)
        })
    })
  },

  updateContentIntegration(
    context: ActionContext<RootState, RootState>,
    payload: UpdateContentIntegrationRequest
  ) {
    return new Promise<EntitiesContentIntegration>((resolve, reject) => {
      context.commit('clearErrorMessages')
      axios
        .put<ResponsesOkEntitiesContentIntegration>(
          `/v1/projects/${payload.projectId}/content-integrations/${payload.integrationId}`,
          payload
        )
        .then(
          async (
            response: AxiosResponse<ResponsesOkEntitiesContentIntegration>
          ) => {
            await Promise.all([
              context.dispatch('addNotification', {
                message:
                  response.data.message ??
                  'Content integration updated successfully',
                type: 'success',
              }),
            ])
            resolve(response.data.data)
          }
        )
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.commit('setErrorMessages', getErrorMessages(error)),
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Validation errors while updating content integration',
              type: 'error',
            }),
          ])
          reject(error)
        })
    })
  },

  deleteContentIntegration(
    context: ActionContext<RootState, RootState>,
    payload: UpdateContentIntegrationRequest
  ) {
    return new Promise<boolean>((resolve, reject) => {
      axios
        .delete<ResponsesNoContent>(
          `/v1/projects/${payload.projectId}/content-integrations/${payload.integrationId}`
        )
        .then(async (response: AxiosResponse<ResponsesNoContent>) => {
          await Promise.all([
            context.dispatch('addNotification', {
              message:
                response.data.message ??
                'Content integration deleted successfully',
              type: 'success',
            }),
          ])
          resolve(true)
        })
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Validation errors while deleting content integration',
              type: 'error',
            }),
          ])
          reject(error)
        })
    })
  },

  addContentIntegration(
    context: ActionContext<RootState, RootState>,
    payload: AddContentIntegrationRequest
  ) {
    return new Promise<EntitiesContentIntegration>((resolve, reject) => {
      context.commit('clearErrorMessages')
      axios
        .post<ResponsesOkEntitiesContentIntegration>(
          `/v1/projects/${payload.projectId}/content-integrations`,
          payload
        )
        .then(
          async (
            response: AxiosResponse<ResponsesOkEntitiesContentIntegration>
          ) => {
            await Promise.all([
              context.dispatch('addNotification', {
                message:
                  response.data.message ??
                  'Content integration added successfully',
                type: 'success',
              }),
            ])
            resolve(response.data.data)
          }
        )
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.commit('setErrorMessages', getErrorMessages(error)),
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Validation errors while adding content integration',
              type: 'error',
            }),
          ])
          reject(error)
        })
    })
  },

  getContentIntegration(
    context: ActionContext<RootState, RootState>,
    payload: ProjectIntegrationIdRequest
  ) {
    return new Promise<EntitiesContentIntegration>((resolve, reject) => {
      axios
        .get<ResponsesOkEntitiesContentIntegration>(
          `/v1/projects/${payload.projectId}/content-integrations/${payload.integrationId}`
        )
        .then(
          (response: AxiosResponse<ResponsesOkEntitiesContentIntegration>) => {
            resolve(response.data.data)
          }
        )
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Error while fetching content integration',
              type: 'error',
            }),
          ])
          reject(error)
        })
    })
  },

  updatePhoneCallIntegration(
    context: ActionContext<RootState, RootState>,
    payload: UpdatePhoneCallIntegrationRequest
  ) {
    return new Promise<EntitiesPhoneCallIntegration>((resolve, reject) => {
      context.commit('clearErrorMessages')
      axios
        .put<ResponsesOkEntitiesPhoneCallIntegration>(
          `/v1/projects/${payload.projectId}/phone-call-integrations/${payload.integrationId}`,
          payload
        )
        .then(
          async (
            response: AxiosResponse<ResponsesOkEntitiesPhoneCallIntegration>
          ) => {
            await Promise.all([
              context.dispatch('addNotification', {
                message:
                  response.data.message ??
                  'Phone Call integration updated successfully',
                type: 'success',
              }),
            ])
            resolve(response.data.data)
          }
        )
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.commit('setErrorMessages', getErrorMessages(error)),
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Validation errors while updating phone call integration',
              type: 'error',
            }),
          ])
          reject(error)
        })
    })
  },

  deletePhoneCallIntegration(
    context: ActionContext<RootState, RootState>,
    payload: ProjectIntegrationIdRequest
  ) {
    return new Promise<boolean>((resolve, reject) => {
      axios
        .delete<ResponsesNoContent>(
          `/v1/projects/${payload.projectId}/phone-call-integrations/${payload.integrationId}`
        )
        .then(async (response: AxiosResponse<ResponsesNoContent>) => {
          await Promise.all([
            context.dispatch('addNotification', {
              message:
                response.data.message ??
                'Phone call integration deleted successfully',
              type: 'success',
            }),
          ])
          resolve(true)
        })
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Validation errors while deleting phone call integration',
              type: 'error',
            }),
          ])
          reject(error)
        })
    })
  },

  addPhoneCallIntegration(
    context: ActionContext<RootState, RootState>,
    payload: AddPhoneCallIntegrationRequest
  ) {
    return new Promise<EntitiesPhoneCallIntegration>((resolve, reject) => {
      context.commit('clearErrorMessages')
      axios
        .post<ResponsesOkEntitiesPhoneCallIntegration>(
          `/v1/projects/${payload.projectId}/phone-call-integrations`,
          payload
        )
        .then(
          async (
            response: AxiosResponse<ResponsesOkEntitiesPhoneCallIntegration>
          ) => {
            await Promise.all([
              context.dispatch('addNotification', {
                message:
                  response.data.message ??
                  'Phone call integration added successfully',
                type: 'success',
              }),
            ])
            resolve(response.data.data)
          }
        )
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.commit('setErrorMessages', getErrorMessages(error)),
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Validation errors while adding phone call integration',
              type: 'error',
            }),
          ])
          reject(error)
        })
    })
  },

  getPhoneCallIntegration(
    context: ActionContext<RootState, RootState>,
    payload: ProjectIntegrationIdRequest
  ) {
    return new Promise<EntitiesContentIntegration>((resolve, reject) => {
      axios
        .get<ResponsesOkEntitiesContentIntegration>(
          `/v1/projects/${payload.projectId}/phone-call-integrations/${payload.integrationId}`
        )
        .then(
          (response: AxiosResponse<ResponsesOkEntitiesContentIntegration>) => {
            resolve(response.data.data)
          }
        )
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Error while fetching phone call integration',
              type: 'error',
            }),
          ])
          reject(error)
        })
    })
  },

  getLinkIntegration(
    context: ActionContext<RootState, RootState>,
    payload: ProjectIntegrationIdRequest
  ) {
    return new Promise<EntitiesLinkIntegration>((resolve, reject) => {
      axios
        .get<ResponsesOkEntitiesLinkIntegration>(
          `/v1/projects/${payload.projectId}/link-integrations/${payload.integrationId}`
        )
        .then((response: AxiosResponse<ResponsesOkEntitiesLinkIntegration>) => {
          resolve(response.data.data)
        })
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Error while fetching link integration',
              type: 'error',
            }),
          ])
          reject(error)
        })
    })
  },

  updateLinkIntegration(
    context: ActionContext<RootState, RootState>,
    payload: UpdateLinkIntegrationRequest
  ) {
    return new Promise<EntitiesLinkIntegration>((resolve, reject) => {
      context.commit('clearErrorMessages')
      axios
        .put<ResponsesOkEntitiesLinkIntegration>(
          `/v1/projects/${payload.projectId}/link-integrations/${payload.integrationId}`,
          payload
        )
        .then(
          async (
            response: AxiosResponse<ResponsesOkEntitiesLinkIntegration>
          ) => {
            await Promise.all([
              context.dispatch('addNotification', {
                message:
                  response.data.message ??
                  'Link integration updated successfully',
                type: 'success',
              }),
            ])
            resolve(response.data.data)
          }
        )
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.commit('setErrorMessages', getErrorMessages(error)),
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Validation errors while updating link integration',
              type: 'error',
            }),
          ])
          reject(error)
        })
    })
  },

  deleteLinkIntegration(
    context: ActionContext<RootState, RootState>,
    payload: ProjectIntegrationIdRequest
  ) {
    return new Promise<boolean>((resolve, reject) => {
      axios
        .delete<ResponsesNoContent>(
          `/v1/projects/${payload.projectId}/link-integrations/${payload.integrationId}`
        )
        .then(async (response: AxiosResponse<ResponsesNoContent>) => {
          await Promise.all([
            context.dispatch('addNotification', {
              message:
                response.data.message ??
                'Link integration deleted successfully',
              type: 'success',
            }),
          ])
          resolve(true)
        })
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Validation errors while deleting link integration',
              type: 'error',
            }),
          ])
          reject(error)
        })
    })
  },

  addLinkIntegration(
    context: ActionContext<RootState, RootState>,
    payload: AddLinkIntegrationRequest
  ) {
    return new Promise<EntitiesLinkIntegration>((resolve, reject) => {
      context.commit('clearErrorMessages')
      axios
        .post<ResponsesOkEntitiesLinkIntegration>(
          `/v1/projects/${payload.projectId}/link-integrations`,
          payload
        )
        .then(
          async (
            response: AxiosResponse<ResponsesOkEntitiesLinkIntegration>
          ) => {
            await Promise.all([
              context.dispatch('addNotification', {
                message:
                  response.data.message ??
                  'Phone call integration added successfully',
                type: 'success',
              }),
            ])
            resolve(response.data.data)
          }
        )
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.commit('setErrorMessages', getErrorMessages(error)),
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Validation errors while adding link integration',
              type: 'error',
            }),
          ])
          reject(error)
        })
    })
  },

  getProjectIntegrations(
    context: ActionContext<RootState, RootState>,
    projectId: string
  ) {
    return new Promise<Array<EntitiesProjectIntegration>>((resolve, reject) => {
      axios
        .get<ResponsesOkArrayEntitiesProjectIntegration>(
          `/v1/projects/${projectId}/integrations`
        )
        .then(
          (
            response: AxiosResponse<ResponsesOkArrayEntitiesProjectIntegration>
          ) => {
            resolve(response.data.data)
          }
        )
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Error while fetching project integration',
              type: 'error',
            }),
          ])
          reject(error)
        })
    })
  },

  updateProjectIntegrations(
    context: ActionContext<RootState, RootState>,
    payload: UpdateProjectIntegrationsRequest
  ) {
    return new Promise<EntitiesLinkIntegration>((resolve, reject) => {
      context.commit('clearErrorMessages')
      axios
        .put<ResponsesOkEntitiesLinkIntegration>(
          `/v1/projects/${payload.projectId}/integrations/`,
          payload
        )
        .then(
          async (
            response: AxiosResponse<ResponsesOkEntitiesLinkIntegration>
          ) => {
            await Promise.all([
              context.dispatch('addNotification', {
                message:
                  response.data.message ?? 'Integrations updated successfully',
                type: 'success',
              }),
            ])
            resolve(response.data.data)
          }
        )
        .catch(async (error: AxiosError) => {
          await Promise.all([
            context.commit('setErrorMessages', getErrorMessages(error)),
            context.dispatch('addNotification', {
              message:
                error.response?.data?.message ??
                'Validation errors while updating integrations',
              type: 'error',
            }),
          ])
          reject(error)
        })
    })
  },
  setNextRoute: (
    context: ActionContext<RootState, RootState>,
    route: string | null
  ) => {
    context.commit('setNextRoute', route)
  },

  addNotification(
    context: ActionContext<RootState, RootState>,
    request: NotificationRequest
  ) {
    context.commit('setNotification', request)
  },

  disableNotification(context: ActionContext<RootState, RootState>) {
    context.commit('disableNotification')
  },
}
