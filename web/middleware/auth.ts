import { Context, Middleware } from '@nuxt/types'

const authMiddleware: Middleware = (context: Context) => {
  if (context.store.getters.authUser === null) {
    context.redirect('/login')
  }
}

export default authMiddleware
