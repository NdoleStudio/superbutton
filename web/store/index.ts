import { GetterTree, ActionTree, MutationTree, ActionContext } from 'vuex'
import { AppData, AuthUser, NotificationRequest, State } from '~/store/types'
import { EntitiesUser, ResponsesOkEntitiesUser } from '~/store/backend'
import axios, { setAuthToken } from '~/plugins/axios'

export const state = (): State => ({
  authUser: null,
  nextRoute: null,
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

  async loadUser({ commit }) {
    const user = await axios.get<ResponsesOkEntitiesUser>('/v1/users/me')
    commit('setUser', user)
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
