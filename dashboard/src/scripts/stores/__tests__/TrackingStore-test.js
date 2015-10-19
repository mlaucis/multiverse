import { decamelizeKeys } from 'humps'

jest.dontMock('../TrackingStore')

import AccountConstants from '../../constants/AccountConstants'
import ApplicationConstants from '../../constants/ApplicationConstants'
import MemberConstants from '../../constants/MemberConstants'

describe('TrackingStore', () => {
  let ConsoleDispatcher
  let callback
  let identifyMock = jest.genMockFunction()
  let trackMock = jest.genMockFunction()
  let trackPidMock = jest.genMockFunction()
  let app = {
    description: 'instant',
    id: 54321,
    name: 'Yo'
  }
  let invitee = {
    email: 'wdg@cat.biz',
    firstName: 'Werner',
    lastName: 'Del Garda'
  }
  let meta = {
    originalReferrer: 'twitter.com',
    referrer: 'lolcat.biz'
  }
  let org = {
    id: 1234,
    name: 'Cat Biz'
  }
  let user = {
    accountId: 123,
    email: 'nyn@cat.biz',
    id: 1234,
    firstName: 'cat',
    lastName: 'nyn',
    password: '4321'
  }

  Object.defineProperty(window, 'analytics', {
    value: {
      identify: identifyMock,
      track: trackMock
    }
  })
  Object.defineProperty(window, 'twttr', {
    value: {
      conversion: {
        trackPid: trackPidMock
      }
    }
  })

  beforeEach( () => {
    require('../TrackingStore')

    ConsoleDispatcher = require('../../dispatcher/ConsoleDispatcher')
    callback = ConsoleDispatcher.register.mock.calls[0][0]

    let AccountStore = require('../AccountStore')

    AccountStore.account = org
    AccountStore.user = user
  })

  it('registers a callback with the dispatcher', () => {
    expect(ConsoleDispatcher.register).toBeCalled()
  })

  it('tracks ACCOUNTUSER_CREATE_FAILURE', () => {
    callback({
      type: AccountConstants.ACCOUNTUSER_CREATE_FAILURE,
      error: { errors: [] },
      originalReferrer: meta.originalReferrer,
      referrer: meta.referrer,
      ...user
    })

    expect(trackMock).lastCalledWith('Member signed up', {
      eventId: 1,
      email: user.email,
      firstName: user.firstName,
      lastName: user.lastName,
      organizationId: user.accountId,
      originalReferrer: meta.originalReferrer,
      referrer: meta.referrer,
      success: false
    })
  })

  it('tracks ACCOUNTUSER_CREATE_SUCCESS', () => {
    callback({
      type: AccountConstants.ACCOUNTUSER_CREATE_SUCCESS,
      response: { id: user.id },
      originalReferrer: meta.originalReferrer,
      referrer: meta.referrer,
      ...user
    })

    expect(trackMock).lastCalledWith('Member signed up', {
      eventId: 1,
      email: user.email,
      firstName: user.firstName,
      lastName: user.lastName,
      memberId: user.id,
      organizationId: user.accountId,
      originalReferrer: meta.originalReferrer,
      referrer: meta.referrer,
      success: true
    })
    expect(trackPidMock).lastCalledWith('ntlyu', decamelizeKeys({
      twSaleAmount: 0,
      twOrderQuantity: 0
    }))
  })

  it('tracks LOGIN_FAILURE', () => {
    callback({
      error: { errors: [] },
      type: AccountConstants.LOGIN_FAILURE
    })

    expect(trackMock).lastCalledWith('Member logged in', {
      eventId: 2,
      success: false
    })
  })

  it('tracks LOGIN_SUCCESS', () => {
    callback({
      type: AccountConstants.LOGIN_SUCCESS,
      response: { id: user.id }
    })

    expect(trackMock).lastCalledWith('Member logged in', {
      eventId: 2,
      memberId: user.id,
      success: true
    })
  })

  it('tracks LOGOUT_FAILURE', () => {
    callback({
      error: { errors: [] },
      type: AccountConstants.LOGOUT_FAILURE,
      user: { id: user.id }
    })

    expect(trackMock).lastCalledWith('Member logged out', {
      eventId: 3,
      memberId: user.id,
      success: false
    })
  })

  it('tracks LOGOUT_SUCCESS', () => {
    callback({
      type: AccountConstants.LOGOUT_SUCCESS,
      user: { id: user.id }
    })

    expect(trackMock).lastCalledWith('Member logged out', {
      eventId: 3,
      memberId: user.id,
      success: true
    })
  })

  it('tracks APP_CREATE_FAILURE', () => {
    callback({
      type: ApplicationConstants.APP_CREATE_FAILURE,
      error: { errors: [] },
      name: app.name,
      description: app.description,
      manual: true
    })

    expect(trackMock).lastCalledWith('Application created', {
      eventId: 7,
      appName: app.name,
      appDescription: app.description,
      manually: true,
      success: false
    })
  })

  it('tracks APP_CREATE_SUCCESS', () => {
    callback({
      type: ApplicationConstants.APP_CREATE_SUCCESS,
      response: { id: app.id },
      name: app.name,
      description: app.description,
      manual: true
    })

    expect(trackMock).lastCalledWith('Application created', {
      eventId: 7,
      appId: app.id,
      appName: app.name,
      appDescription: app.description,
      manually: true,
      success: true
    })
  })

  it('tracks APP_EDIT_FAILURE', () => {
    callback({
      type: ApplicationConstants.APP_EDIT_FAILURE,
      description: app.description,
      error: { errors: [] },
      id: app.id,
      name: app.name,
      manual: true
    })

    expect(trackMock).lastCalledWith('Application edited', {
      eventId: 8,
      appId: app.id,
      appName: app.name,
      appDescription: app.description,
      manually: true,
      success: false
    })
  })

  it('tracks APP_EDIT_SUCCESS', () => {
    callback({
      type: ApplicationConstants.APP_EDIT_SUCCESS,
      description: app.description,
      error: { errors: [] },
      id: app.id,
      name: app.name,
      manual: true
    })

    expect(trackMock).lastCalledWith('Application edited', {
      eventId: 8,
      appId: app.id,
      appName: app.name,
      appDescription: app.description,
      manually: true,
      success: true
    })
  })

  it('tracks APP_DELETE_FAILURE', () => {
    callback({
      type: ApplicationConstants.APP_DELETE_FAILURE,
      error: { errors: [] },
      id: app.id
    })

    expect(trackMock).lastCalledWith('Application deleted', {
      eventId: 9,
      appId: app.id,
      success: false
    })
  })

  it('tracks APP_DELETE_SUCCESS', () => {
    callback({
      type: ApplicationConstants.APP_DELETE_SUCCESS,
      id: app.id
    })

    expect(trackMock).lastCalledWith('Application deleted', {
      eventId: 9,
      appId: app.id,
      success: true
    })
  })

  it('tracks MEMBER_CREATE_FAILURE', () => {
    callback({
      error: { errors: [] },
      type: MemberConstants.MEMBER_CREATE_FAILURE
    })

    expect(trackMock).lastCalledWith('Member created', {
      eventId: 10,
      success: false
    })
  })

  it('tracks MEMBER_CREATE_SUCCESS', () => {
    callback({
      type: MemberConstants.MEMBER_CREATE_SUCCESS,
      response: { id: user.id }
    })

    expect(trackMock).lastCalledWith('Member created', {
      eventId: 10,
      memberId: user.id,
      success: true
    })
  })

  it('tracks MEMBER_DELETE_FAILURE', () => {
    callback({
      type: MemberConstants.MEMBER_DELETE_FAILURE,
      error: { errors: [] },
      id: user.id
    })

    expect(trackMock).lastCalledWith('Member deleted', {
      eventId: 11,
      memberId: user.id,
      success: false
    })
  })

  it('tracks MEMBER_DELETE_SUCCESS', () => {
    callback({
      type: MemberConstants.MEMBER_DELETE_SUCCESS,
      id: user.id
    })

    expect(trackMock).lastCalledWith('Member deleted', {
      eventId: 11,
      memberId: user.id,
      success: true
    })
  })

  it('tracks MEMBER_INVITE_FAILURE', () => {
    callback({
      type: MemberConstants.MEMBER_INVITE_FAILURE,
      error: { errors: [] },
      email: invitee.email,
      firstName: invitee.firstName,
      lastName: invitee.lastName
    })

    expect(trackMock).lastCalledWith('Member invited', {
      eventId: 12,
      invitee: invitee,
      inviter: {
        email: user.email,
        firstName: user.firstName,
        id: user.id,
        lastName: user.lastName
      },
      org: org,
      success: false
    })
  })

  it('tracks MEMBER_INVITE_SUCCESS', () => {
    callback({
      type: MemberConstants.MEMBER_INVITE_SUCCESS,
      email: invitee.email,
      firstName: invitee.firstName,
      lastName: invitee.lastName
    })

    expect(trackMock).lastCalledWith('Member invited', {
      eventId: 12,
      invitee: invitee,
      inviter: {
        email: user.email,
        firstName: user.firstName,
        id: user.id,
        lastName: user.lastName
      },
      org: org,
      success: true
    })
  })
})
