import Cookies from 'js-cookie'

const readmeAppCookie = 'readme_app'
const referrerCookie = 'originalReferrer'

export function consumeReferrerCookie() {
  let referrer = Cookies.get(referrerCookie)

  if (!referrer) {
    return 'unknown'
  }

  Cookies.remove(referrerCookie, {
    domain: '.tapglue.com',
    path: '/'
  })

  return referrer
}

export function setReadmeToken(token) {
  let content = {
    key: token
  }

  Cookies.set(readmeAppCookie, JSON.stringify(content), {
    domain: '.tapglue.com',
    path: '/'
  })
}
