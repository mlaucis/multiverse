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
    if (AccountStore.isAuthenticated) {
      let user = AccountStore.user

      window.analytics.identify(user.id, {
        firstName: user.firstName,
        lastName: user.lastName
      })
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

    switch (action.type) {
      case AccountConstants.ACCOUNTUSER_CREATE_FAILURE:
        this.trackEvent('Member signed up', {
          eventId: 1,
          email: action.email,
          firstName: action.firstName,
          lastName: action.lastName,
          organizationId: action.accountId,
          plan: action.plan,
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
          plan: action.plan,
          referrer: action.referrer,
          success: true
        })

        break
      case AccountConstants.LOGIN_FAILURE:
        this.trackEvent('Member logged in', {
          eventId: 2,
          success: false
        })

        break
      case AccountConstants.LOGIN_SUCCESS:
        let me = action.response

        this.trackEvent('Member logged in', {
          eventId: 2,
          memberId: me.id,
          success: true
        })

        break
      case AccountConstants.LOGOUT_FAILURE:
        this.trackEvent('Member logged out', {
          eventId: 3,
          memberId: action.user.id,
          success: false
        })

        break
      case AccountConstants.LOGOUT_SUCCESS:
        this.trackEvent('Member logged out', {
          eventId: 3,
          memberId: action.user.id,
          success: true
        })

        break
      case ApplicationConstants.APP_CREATE_FAILURE:
        this.trackEvent('Application created', {
          eventId: 7,
          appName: action.name,
          appDescription: action.description,
          manually: true,
          success: false
        })

        break
      case ApplicationConstants.APP_CREATE_SUCCESS:
        this.trackEvent('Application created', {
          eventId: 7,
          appId: action.response.id,
          appName: action.name,
          appDescription: action.description,
          manually: true,
          success: true
        })

        break
      case ApplicationConstants.APP_EDIT_FAILURE:
        this.trackEvent('Application edited', {
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
          eventId: 9,
          appId: action.id,
          success: false
        })

        break
      case ApplicationConstants.APP_DELETE_SUCCESS:
        this.trackEvent('Application deleted', {
          eventId: 9,
          appId: action.id,
          success: true
        })

        break
      case MemberConstants.MEMBER_CREATE_FAILURE:
        this.trackEvent('Member created', {
          eventId: 10,
          success: false
        })

        break
      case MemberConstants.MEMBER_CREATE_SUCCESS:
        let member = action.response

        this.trackEvent('Member created', {
          eventId: 10,
          memberId: member.id,
          success: true
        })

        break
      case MemberConstants.MEMBER_DELETE_FAILURE:
        this.trackEvent('Member deleted', {
          eventId: 11,
          memberId: action.id,
          success: false
        })

        break
      case MemberConstants.MEMBER_DELETE_SUCCESS:
        this.trackEvent('Member deleted', {
          eventId: 11,
          memberId: action.id,
          success: true
        })

        break
      case MemberConstants.MEMBER_INVITE_FAILURE:
        this.trackEvent('Member invited', {
          eventId: 12,
          email: action.email,
          success: false
        })

        break
      case MemberConstants.MEMBER_INVITE_SUCCESS:
        this.trackEvent('Member invited', {
          eventId: 12,
          email: action.email,
          success: true
        })

        break
      default:
      // nothing to do
    }
  }
}

export default new TrackingStore
