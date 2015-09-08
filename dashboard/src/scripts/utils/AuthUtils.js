import Cookies from 'js-cookie'

const referrerCookie = 'originalReferrer'

export function consumeReferrerCookie() {
  let referrer = Cookies.get(referrerCookie)

  if (!referrer) {
    return 'unknown'
  }

  Cookies.remove(referrerCookie, { domain: '.tapglue.com', path: '/' })

  return referrer
}
