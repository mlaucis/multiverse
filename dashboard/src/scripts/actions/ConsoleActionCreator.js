import AccountConstants from '../constants/AccountConstants'
import ApplicationConstants from '../constants/ApplicationConstants'
import MemberConstants from '../constants/MemberConstants'
import OnboardingConstants from '../constants/OnboardingConstants'

import { dispatch, dispatchAsync } from '../dispatcher/ConsoleDispatcher'

import { accountCreate, accountGet } from '../services/AccountService'
import { accountUserCreate } from '../services/AccountService'
import { app, apps } from '../services/ApplicationService'
import { appCreate, appDelete, appUpdate } from '../services/ApplicationService'
import { login, logout } from '../services/AccountService'
import { member, memberDelete, memberUpdate } from '../services/MemberService'
import { members } from '../services/MemberService'

export function requestAccount(user) {
  return dispatchAsync(accountGet(user), {
    request: AccountConstants.ACCOUNT_REQUEST,
    success: AccountConstants.ACCOUNT_SUCCESS,
    failure: AccountConstants.ACCOUNT_FAILURE
  }, { user: user })
}

export function requestAccountCreate(vals) {
  return dispatchAsync(accountCreate(
    vals.accountName,
    vals.accountDescription
  ), {
    request: AccountConstants.ACCOUNT_CREATE_REQUEST,
    success: AccountConstants.ACCOUNT_CREATE_SUCCESS,
    failure: AccountConstants.ACCOUNT_CREATE_FAILURE
  })
}

export function requestAccountUserCreate(vals, accountID) {
  return dispatchAsync(accountUserCreate(
    vals.email,
    vals.password,
    vals.firstName,
    vals.lastName,
    accountID
  ), {
    request: AccountConstants.ACCOUNTUSER_CREATE_REQUEST,
    success: AccountConstants.ACCOUNTUSER_CREATE_SUCCESS,
    failure: AccountConstants.ACCOUNTUSER_CREATE_FAILURE
  }, {
    accountId: accountID,
    email: vals.email,
    firstName: vals.firstName,
    lastName: vals.lastName,
    password: vals.password
  })
}

export function requestApp(id, user) {
  dispatchAsync(app(id, user), {
    request: ApplicationConstants.APP_REQUEST,
    success: ApplicationConstants.APP_SUCCESS,
    failure: ApplicationConstants.APP_FAILURE
  })
}

export function requestApps(user) {
  dispatchAsync(apps(user), {
    request: ApplicationConstants.APPS_REQUEST,
    success: ApplicationConstants.APPS_SUCCESS,
    failure: ApplicationConstants.APPS_FAILURE
  })
}

export function requestAppCreate(name, description, user) {
  return dispatchAsync(appCreate(name, description, user), {
    request: ApplicationConstants.APP_CREATE_REQUEST,
    success: ApplicationConstants.APP_CREATE_SUCCESS,
    failure: ApplicationConstants.APP_CREATE_FAILURE
  }, { name: name, description: description })
}

export function requestAppDelete(id, user) {
  return dispatchAsync(appDelete(id, user), {
    request: ApplicationConstants.APP_DELETE_REQUEST,
    success: ApplicationConstants.APP_DELETE_SUCCESS,
    failure: ApplicationConstants.APP_DELETE_FAILURE
  }, { id: id })
}

export function requestAppUpdate(id, name, description, user) {
  return dispatchAsync(appUpdate(id, name, description, user), {
    request: ApplicationConstants.APP_EDIT_REQUEST,
    success: ApplicationConstants.APP_EDIT_SUCCESS,
    failure: ApplicationConstants.APP_EDIT_FAILURE
  }, { id: id, name: name, description: description })
}

export function requestLogin(email, password) {
  return dispatchAsync(login(email, password), {
    request: AccountConstants.LOGIN_REQUEST,
    success: AccountConstants.LOGIN_SUCCESS,
    failure: AccountConstants.LOGIN_FAILURE
  }, { email: email, password: password })
}

export function requestLogout(user) {
  return dispatchAsync(logout(user), {
    request: AccountConstants.LOGOUT_REQUEST,
    success: AccountConstants.LOGOUT_SUCCESS,
    failure: AccountConstants.LOGOUT_FAILURE
  }, { user: user })
}

export function requestMember(id, user) {
  return dispatchAsync(member(id, user), {
    request: MemberConstants.MEMBER_REQUEST,
    success: MemberConstants.MEMBER_SUCCESS,
    failure: MemberConstants.MEMBER_FAILURE
  })
}

export function requestMemberCreate(vals, accountID) {
  return dispatchAsync(accountUserCreate(
    vals.email,
    vals.password,
    vals.firstName,
    vals.lastName,
    accountID
  ), {
    request: MemberConstants.MEMBER_CREATE_REQUEST,
    success: MemberConstants.MEMBER_CREATE_SUCCESS,
    failure: MemberConstants.MEMBER_CREATE_FAILURE
  })
}

export function requestMemberDelete(id, user) {
  return dispatchAsync(memberDelete(id, user), {
    request: MemberConstants.MEMBER_DELETE_REQUEST,
    success: MemberConstants.MEMBER_DELETE_SUCCESS,
    failure: MemberConstants.MEMBER_DELETE_FAILURE
  }, { id: id })
}

export function requestMemberInvite(email) {
  // TODO(xla): Needs to properly hooked-up with calls to the backend.
  return new Promise( resolve => {
    dispatch(MemberConstants.MEMBER_INVITE_SUCCESS, { email: email })
    resolve(email)
  })
}

export function requestMemberUpdate(vals, id, accountID) {
  return dispatchAsync(memberUpdate(
    vals.email,
    vals.password,
    vals.firstName,
    vals.lastName,
    id,
    accountID
  ), {
    request: MemberConstants.MEMBER_UPDATE_REQUEST,
    success: MemberConstants.MEMBER_UPDATEE_SUCCESS,
    failure: MemberConstants.MEMBER_UPDATE_FAILURE
  })
}

export function requestMembers(user) {
  return dispatchAsync(members(user), {
    request: MemberConstants.MEMBERS_REQUEST,
    success: MemberConstants.MEMBERS_SUCCESS,
    failure: MemberConstants.MEMBERS_FAILURE
  })
}

export function selectOptions(options) {
  dispatch(OnboardingConstants.SELECT_OPTIONS, { options: options })
}

export function selectPersona(persona) {
  dispatch(OnboardingConstants.SELECT_PERSONA, { persona: persona })
}
