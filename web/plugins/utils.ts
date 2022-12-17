export const setAuthCookie = () => {
  setCookie('auth', 'true')
}

export const removeAuthCookie = () => {
  eraseCookie('auth')
}

const getRootDomain = () => {
  const domain = window.location.origin
    .replace('https://', '')
    .replace('http://', '')
    .split('.')
  if (domain.length === 1) {
    return domain
  }
  return domain[domain.length - 2] + '.' + domain[domain.length - 1]
}

const setCookie = (key: string, value: string) => {
  const expirationDate = new Date()
  expirationDate.setMonth(expirationDate.getMonth() + 1)
  document.cookie = `${key}=${value};expires=${expirationDate};domain=.${getRootDomain()};path=/`
}

const eraseCookie = (key: string) => {
  document.cookie = `${key}=;domain=.${getRootDomain()}expires=Thu, 01 Jan 1970 00:00:00 GMT`
}
