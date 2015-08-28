import React, { Component, PropTypes, findDOMNode } from 'react'

import AccountStore from '../stores/AccountStore'
import ApplicationStore from '../stores/ApplicationStore'
import { requestApps } from '../actions/ConsoleActionCreator'
import { requestMemberInvite } from '../actions/ConsoleActionCreator'

import { App } from './Apps'

class InviteDeveloper extends Component {
  constructor() {
    super()

    this._showSuccess = false
    this.state = this.getState()
  }

  getState() {
    return {
      showSuccess: this._showSuccess
    }
  }

  handleSubmit = (event) => {
    event.preventDefault()

    let email = findDOMNode(this.refs.email).value
    let toggle = this.toggleSuccess

    requestMemberInvite(email).then( () => {
      toggle()
    })
  }

  viewDefault() {
    return (
      <div className='portlet light'>
        <h3>
            Send an Invitation Email
        </h3>
        <div className='portlet-body'>
          <form onSubmit={this.handleSubmit}>
            <div className='group'>
              <input
                id='developer-email'
                placeholder='Email Address'
                ref='email'
                required
                type='email'/>
              <span className='bar'></span>
              <span
                className='help'>
                Valid email in the form of member@company.org
              </span>
              <label htmlFor='developer-email'>Email Address</label>
            </div>
            <div className='actions'>
              <input className='btn-default' type='submit' value='Invite'/>
            </div>
          </form>
        </div>
      </div>
    )
  }

  viewSuccess() {
    return (
      <div className='note note-block note-success'>
        <h2>Your invite is out!</h2>
        <p>
          We sent an invite to your team member with detailed instructions how
          to proceed.
        </p>
        <div className='actions'>
          <button
            className='btn-secondary'
            href='#'
            onClick={this.toggleSuccess}>
            Close
          </button>
        </div>
      </div>
    )
  }

  toggleSuccess = () => {
    this._showSuccess = !this._showSuccess
    this.setState(this.getState())
  }

  render() {
    let view = this.viewDefault()

    if (this.state.showSuccess) {
      view = this.viewSuccess()
    }

    return view
  }
}

class IntegrationResource extends Component {
  static propTypes = {
    double: PropTypes.bool,
    icon: PropTypes.string.isRequired,
    link: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired
  }

  render() {
    let c = 'tile tpgl-green'
    let iconClass = `fa fa-${this.props.icon}`

    if (this.props.double) {
      c += ' double'
    }

    return (
      <a className={c} href={this.props.link}>
        <div className='tile-body'>
          <span className={iconClass}></span>
        </div>
        <div className='tile-object'>
          {this.props.name}
        </div>
      </a>
    )
  }
}

class ProductResource extends Component {
  static propTypes = {
    children: PropTypes.node.isRequired,
    icon: PropTypes.string.isRequired,
    link: PropTypes.string.isRequired,
    title: PropTypes.string.isRequired
  }

  render() {
    let iconClass = `glyphicon glyphicon-${this.props.icon}`

    return (
      <a className='dashboard-stat tpgl-green' href={this.props.link}>
        <div className='visual'>
          <i className={iconClass}/>
        </div>
        <div className='details'>
          <div className='number'>{this.props.title}</div>
          <div className='desc'>{this.props.children}</div>
        </div>
      </a>
    )
  }
}

export default class Dashboard extends Component {
  constructor() {
    super()

    this.state = this.getState()
  }

  componentDidMount() {
    AccountStore.addChangeListener(this.handleChange)
    ApplicationStore.addChangeListener(this.handleChange)

    requestApps(AccountStore.user)
  }

  componentWillUnmount() {
    AccountStore.removeChangeListener(this.handleChange)
    ApplicationStore.removeChangeListener(this.handleChange)
  }

  getState() {
    return {
      app: ApplicationStore.apps[0],
      user: AccountStore.user
    }
  }

  handleChange = () => {
    this.setState(this.getState())
  }

  handleClick = (event) => {
    event.preventDefault()

    this.setState(this.getState())
  }
  render() {
    let app = this.state.app ? (
      <App actions={false} app={this.state.app}/>
    ) : (
      <div></div>
    )

    return (
      <div>
        <h1>Welcome to Tapglue, {this.state.user.firstName}</h1>
        <div className='row'>
          <div className='col-md-6'>
            <h2>Manage your App</h2>
            {app}
          </div>
          <div className='col-md-6'>
            <h2>Invite your Team</h2>
            <InviteDeveloper/>
          </div>
        </div>
        <div className='row'>
          <div className='col-md-12'>
            <div className='portlet light'>
              <h2>Product</h2>
              <div className='row'>
                <div className='col-md-4'>
                  <ProductResource
                    icon='check'
                    link='#'
                    title='Product'>
                    <strong>Success</strong> Guide
                  </ProductResource>
                </div>
                <div className='col-md-4'>
                  <ProductResource
                    icon='repeat'
                    link='#'
                    title='Solving'>
                    the <strong>Empty Room</strong> Problem
                  </ProductResource>
                </div>
                <div className='col-md-4'>
                  <ProductResource
                    icon='list-alt'
                    link='#'
                    title='Seven Ways'>
                    <strong>Activity Feeds</strong> will boost your App
                  </ProductResource>
                </div>
              </div>
              <h2>Integration</h2>
              <div className='tiles'>
                <IntegrationResource
                  icon='rocket'
                  link='#'
                  name='Getting Started'/>
                <IntegrationResource
                  icon='bars'
                  link='#'
                  name='API Reference'/>
                <IntegrationResource
                  double={true}
                  icon='apple'
                  link='#'
                  name='iOS Integration Guide & SDK'/>
                <IntegrationResource
                  icon='code'
                  link='#'
                  name='Examples'/>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }
}
