import { GetterTree, ActionTree, MutationTree, ActionContext } from 'vuex'
import { AxiosError, AxiosResponse } from 'axios'
import {
  AddWhatsappIntegrationRequest,
  AppData,
  AuthUser,
  NotificationRequest,
  State,
  UpdateProjectRequest,
} from '~/store/types'
import {
  EntitiesProject,
  EntitiesUser,
  ResponsesOkArrayEntitiesProject,
  ResponsesOkEntitiesProject,
  ResponsesOkEntitiesUser,
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
    context.commit('setAuthUser', { uid, email, photoURL, displayName })
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
    if (activeProject === undefined) {
      context.commit('setActiveProjectId', projects[0].id)
    }

    await context.commit('setProjects', projects)
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
            context.commit('setActiveProjectId', response.data.data.id),
            context.commit('setCreatingProject', false),
          ])
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

  addWhatsappIntegration(
    context: ActionContext<RootState, RootState>,
    payload: AddWhatsappIntegrationRequest
  ) {
    return new Promise<EntitiesProject>((resolve, reject) => {
      context.commit('clearErrorMessages')
      axios
        .post<ResponsesOkEntitiesProject>(
          `/v1/projects/${payload.projectId}/whatsapp-integrations`,
          payload
        )
        .then(async (response: AxiosResponse<ResponsesOkEntitiesProject>) => {
          await Promise.all([
            context.dispatch('addNotification', {
              message:
                response.data.message ??
                'Whatsapp integration added successfully',
              type: 'success',
            }),
          ])
          resolve(response.data.data)
        })
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
