import '../styles/theme.less'

import React, { Component } from 'react'
import ReactDOM from 'react-dom'
import Router from 'react-router'

let { DefaultRoute, Route, RouteHandler } = Router

import Analytics from './components/Analytics'
import Apps from './components/Apps'
import Auth from './components/Auth'
import AuthLogin from './components/auth/Login'
import AuthLogout from './components/auth/Logout'
import AuthSignup from './components/auth/Signup'
import Console from './components/Console'
import Dashboard from './components/Dashboard'
import Members from './components/Members'
import Onboarding from './components/Onboarding'

import RouteConstants from './constants/RouteConstants'

import AccountStore from './stores/AccountStore'
import TrackingStore from './stores/TrackingStore'

window.React = React

function requireAnonymous(Element) {
  return class Anonymous extends Component {
    static willTransitionTo(transition) {
      if (AccountStore.isAuthenticated) {
        transition.redirect('CONSOLE')
      }
    }

    render() {
      return <Element {...this.props}/>
    }
  }
}

function requireAuth(Element) {
  return class Authenticated extends Component {
    static willTransitionTo(transition) {
      if (!AccountStore.isAuthenticated) {
        transition.redirect('AUTH_LOGIN', {}, { 'nextPath': transition.path })
      }
    }

    render() {
      return <Element {...this.props}/>
    }
  }
}

// Initialisation is necessary to restore state from localStorage.
AccountStore.init()

class Wrapper extends Component {
  render() {
    return (
      <RouteHandler/>
    )
  }
}

let routes = (
  <Route name='WRAPPER' handler={Wrapper}>
    <Route name='Auth' path='/auth' handler={Auth}>
      <Route
        name='AUTH_LOGIN'
        path={RouteConstants.AUTH_LOGIN}
        handler={requireAnonymous(AuthLogin)}/>
      <Route
        name='AUTH_LOGOUT'
        path={RouteConstants.AUTH_LOGOUT}
        handler={requireAuth(AuthLogout)}/>

      <Route
        name='ONBOARDING'
        path={RouteConstants.ONBOARDING}
        handler={requireAuth(Onboarding)}/>
    </Route>

    <Route name='CONSOLE' path='/' handler={Console}>
      <DefaultRoute name='DASHBOARD'
        handler={requireAuth(Dashboard)}/>

      <Route
        name='APPS'
        path={RouteConstants.APPS}
        handler={Apps}/>

      <Route
        name='MEMBERS'
        path={RouteConstants.MEMBERS}
        handler={requireAuth(Members)}/>

      <Route
        name='ANALYTICS'
        path={RouteConstants.ANALYTICS}
        handler={requireAuth(Analytics)}/>
    </Route>
  </Route>
)

Router.run(routes, Router.HistoryLocation, (Root) => {
  ReactDOM.render(<Root/>, document.getElementById('container'))
  TrackingStore.trackPage()
})
