export const hasAuthCookie = (): boolean => {
  if (typeof window !== 'undefined') {
    return (
      (document.cookie
        .match('(^|;)\\s*' + 'auth' + '\\s*=\\s*([^;]+)')
        ?.pop() || '') === 'true'
    )
  }
  return false
}
