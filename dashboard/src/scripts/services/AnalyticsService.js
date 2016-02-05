import AccountStore from '../stores/AccountStore'

import { request } from '../utils/APIUtils'

export function metrics(app, start, end) {
  let account = AccountStore.account
  let where = {end: end, start: start}
  let param = encodeURIComponent(JSON.stringify(where))

  return request(
    'get',
    `/orgs/${account.id}/apps/${app}/analytics?where=${param}`,
  )
}
