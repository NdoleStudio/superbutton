import { Context, Middleware } from '@nuxt/types'

const guestMiddleware: Middleware = (context: Context) => {
  if (context.store.getters.authUser !== null) {
    context.redirect('/')
  }
}

export default guestMiddleware
