import { decamelizeKeys } from 'humps'

import AccountConstants from '../constants/AccountConstants'
import ApplicationConstants from '../constants/ApplicationConstants'
import MemberConstants from '../constants/MemberConstants'

import AccountStore from './AccountStore'

import { EventStore } from '../utils/StoreUtils'
import { register, waitFor } from '../dispatcher/ConsoleDispatcher'

class TrackingStore extends EventStore {
  constructor() {
    super()

    this.dispatchToken = register(this.handleAction)
  }

  identify() {
    if (AccountStore.account && AccountStore.user) {
      let account = AccountStore.account
      let user = AccountStore.user

      let payload = {
        company: {
          createdAt: Date.parse(account.createdAt),
          id: account.id,
          name: account.name
        },
        createdAt: Date.parse(user.createdAt),
        email: user.email,
        firstName: user.firstName,
        lastName: user.lastName
      }

      if (account.metadata && account.metadata.plan) {
        payload.company.plan = account.metadata.plan
      }
      if (user.metadata && user.metadata.originalReferrer) {
        payload.originalReferrer = user.metadata.originalReferrer
      }
      if (user.metadata && user.metadata.referrer) {
        payload.referrer = user.metadata.referrer
      }

      window.analytics.identify(user.id, payload)
    }
  }

  trackEvent(event, props) {
    this.identify()

    window.analytics.track(event, props)
  }

  trackPage() {
    this.identify()

    window.analytics.page()
  }

  handleAction = (action) => {
    waitFor([ AccountStore.dispatchToken ])

    let account = AccountStore.account
    let user = AccountStore.user

    switch (action.type) {
      case AccountConstants.ACCOUNTUSER_CREATE_FAILURE:
        this.trackEvent('Member signed up', {
          eventId: 1,
          email: action.email,
          firstName: action.firstName,
          lastName: action.lastName,
          organizationId: action.accountId,
          originalReferrer: action.originalReferrer,
          referrer: action.referrer,
          success: false
        })

        break
      case AccountConstants.ACCOUNTUSER_CREATE_SUCCESS:
        this.trackEvent('Member signed up', {
          eventId: 1,
          email: action.email,
          firstName: action.firstName,
          lastName: action.lastName,
          memberId: action.response.id,
          organizationId: action.accountId,
          originalReferrer: action.originalReferrer,
          referrer: action.referrer,
          success: true
        })
        if (typeof twttr !== 'undefined' && twttr.conversion) {
          twttr.conversion.trackPid('ntlyu', decamelizeKeys({
            twSaleAmount: 0,
            twOrderQuantity: 0
          }))
        }

        break
      case AccountConstants.LOGIN_FAILURE:
        this.trackEvent('Member logged in', {
          eventId: 2,
          success: false
        })

        break
      case AccountConstants.LOGIN_SUCCESS:
        this.trackEvent('Member logged in', {
          email: user.email,
          eventId: 2,
          memberId: user.id,
          success: true
        })

        break
      case AccountConstants.LOGOUT_FAILURE:
        this.trackEvent('Member logged out', {
          email: action.user.email,
          eventId: 3,
          memberId: action.user.id,
          success: false
        })

        break
      case AccountConstants.LOGOUT_SUCCESS:
        this.trackEvent('Member logged out', {
          email: action.user.email,
          eventId: 3,
          memberId: action.user.id,
          success: true
        })

        break
      case ApplicationConstants.APP_CREATE_FAILURE:
        this.trackEvent('Application created', {
          email: user.email,
          eventId: 7,
          appName: action.name,
          appDescription: action.description,
          manually: action.manual,
          success: false
        })

        break
      case ApplicationConstants.APP_CREATE_SUCCESS:
        this.trackEvent('Application created', {
          email: user.email,
          eventId: 7,
          appId: action.response.id,
          appName: action.name,
          appDescription: action.description,
          manually: action.manual,
          success: true
        })

        break
      case ApplicationConstants.APP_EDIT_FAILURE:
        this.trackEvent('Application edited', {
          email: user.email,
          eventId: 8,
          appId: action.id,
          appName: action.name,
          appDescription: action.description,
          manually: true,
          success: false
        })

        break
      case ApplicationConstants.APP_EDIT_SUCCESS:
        this.trackEvent('Application edited', {
          email: user.email,
          eventId: 8,
          appId: action.id,
          appName: action.name,
          appDescription: action.description,
          manually: true,
          success: true
        })

        break
      case ApplicationConstants.APP_DELETE_FAILURE:
        this.trackEvent('Application deleted', {
          email: user.email,
          eventId: 9,
          appId: action.id,
          success: false
        })

        break
      case ApplicationConstants.APP_DELETE_SUCCESS:
        this.trackEvent('Application deleted', {
          email: user.email,
          eventId: 9,
          appId: action.id,
          success: true
        })

        break
      case MemberConstants.MEMBER_CREATE_FAILURE:
        this.trackEvent('Member created', {
          email: user.email,
          eventId: 10,
          success: false
        })

        break
      case MemberConstants.MEMBER_CREATE_SUCCESS:
        let member = action.response

        this.trackEvent('Member created', {
          email: user.email,
          eventId: 10,
          memberId: member.id,
          success: true
        })

        break
      case MemberConstants.MEMBER_DELETE_FAILURE:
        this.trackEvent('Member deleted', {
          email: user.email,
          eventId: 11,
          memberId: action.id,
          success: false
        })

        break
      case MemberConstants.MEMBER_DELETE_SUCCESS:
        this.trackEvent('Member deleted', {
          email: user.email,
          eventId: 11,
          memberId: action.id,
          success: true
        })

        break
      case MemberConstants.MEMBER_INVITE_FAILURE:
        this.trackEvent('Member invited', {
          email: user.email,
          eventId: 12,
          invitee: {
            email: action.email,
            firstName: action.firstName,
            lastName: action.lastName
          },
          inviter: {
            email: _user.email,
            firstName: user.firstName,
            id: user.id,
            lastName: user.lastName
          },
          org: {
            id: account.id,
            name: account.name
          },
          success: false
        })

        break
      case MemberConstants.MEMBER_INVITE_SUCCESS:
        this.trackEvent('Member invited', {
          email: user.email,
          eventId: 12,
          invitee: {
            email: action.email,
            firstName: action.firstName,
            lastName: action.lastName
          },
          inviter: {
            email: user.email,
            firstName: user.firstName,
            id: user.id,
            lastName: user.lastName
          },
          org: {
            id: account.id,
            name: account.name
          },
          success: true
        })

        break
      default:
      // nothing to do
    }
  };
}

export default new TrackingStore
