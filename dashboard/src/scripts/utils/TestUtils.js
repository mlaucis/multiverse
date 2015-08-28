import React, { PropTypes } from 'react'
import objectAssign from 'object-assign'

export let stubRouterContext = (Component, props, stubs) => {
  let RouterStub = objectAssign({}, {
    makePath() {},
    makeHref() {},
    transitionTo() {},
    replaceWith() {},
    goBack() {},
    getCurrentPath() {},
    getCurrentRoutes() {},
    getCurrentPathname() {},
    getCurrentParams() {},
    getCurrentQuery() {},
    isActive() {},
    getRouteAtDepth() {},
    setRouteComponentAtDepth() {}
  }, stubs)

  return React.createClass({
    childContextTypes: {
      router: PropTypes.object,
      routeDepth: PropTypes.number
    },

    getChildContext() {
      return {
        router: RouterStub,
        routeDepth: 0
      }
    },

    render() {
      return <Component {...props} />
    }
  })
 }
