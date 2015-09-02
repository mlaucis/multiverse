import React, { Component, PropTypes, findDOMNode } from 'react'

import AccountStore from '../stores/AccountStore'
import ApplicationStore from '../stores/ApplicationStore'
import { requestApps } from '../actions/ConsoleActionCreator'
import { requestMemberInvite } from '../actions/ConsoleActionCreator'

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
      <div className='card invite-developer'>
        <h3>Get your team on board</h3>
        <p>Invite a developer or others who should help integrate Tapglue
          into your app.</p>
        <form onSubmit={this.handleSubmit}>
          <div className='left'>
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
          <div className='right'>
            <input className='btn-default' type='submit' value='Invite'/>
          </div>
        </form>
      </div>
    )
  }

  viewSuccess() {
    return (
      <div className='card invite-developer note note-block note-success'>
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
    icon: PropTypes.string.isRequired,
    link: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired
  }

  render() {
    let iconClass = `glyphicon glyphicon-${this.props.icon}`

    return (
      <a className='btn-secondary outline resource' href={this.props.link}>
        <span className={iconClass}></span>
        <span>{this.props.name}</span>
      </a>
    )
  }
}

class ProductResource extends Component {
  static propTypes = {
    category: PropTypes.string.isRequired,
    link: PropTypes.string.isRequired,
    title: PropTypes.string.isRequired
  }

  render() {
    return (
      <div className='product-resource'>
        <h4>{this.props.category}</h4>
        <a href={this.props.link}>{this.props.title}</a>
      </div>
    )
  }
}

class TestingApp extends Component {
  static propTypes = {
    app: PropTypes.object.isRequired
  }

  render() {
    return (
      <div className='row testing-app'>
        <div className='col-md-6'>
          <div className='card'>
            <h2>Instant API Access</h2>
            <p>{this.props.app.description}</p>
          </div>
        </div>
        <div className='col-md-6'>
          <div className='card extra'>
            <p>
              <span>API TOKEN:</span>
              <span className='token'>{this.props.app.token}</span>
            </p>
          </div>
        </div>
      </div>
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
    let reasonHref = 'http://www.tapglue.com/blog/' +
      'why-your-mobile-app-needs-a-social-activity-feed/'
    let app = this.state.app ? (
      <TestingApp app={this.state.app}/>
    ) : (
      <div></div>
    )

    return (
      <div className='home'>
        <div className='row teaser'>
          <h1>Welcome to Tapglue, {this.state.user.firstName}!</h1>
          <p>Integrating Tapglue into your app is a matter of a few hours.</p>
          <p>Jump straight into it with these helpful resources.</p>
          <div className='resources'>
            <IntegrationResource
              icon='file'
              link='https://developers.tapglue.com'
              name='API Docs'/>
            <IntegrationResource
              icon='file'
              link='https://developers.tapglue.com/page/ios-guide'
              name='iOS Guide'/>
            <IntegrationResource
              icon='file'
              link='https://github.com/tapglue/ios_sdk#tapglue-ios-sdk'
              name='iOS SDK'/>
          </div>
        </div>
        {app}
        <div className='row'>
          <div className='col-md-6'>
            <InviteDeveloper/>
          </div>
          <div className='col-md-6'>
            <div className='card'>
              <h3>Product Resources</h3>
              <ProductResource
                category='guide'
                link='#'
                title='Launching a Social Activity Feed'/>
              <ProductResource
                category='blogpost'
                link='http://www.tapglue.com/blog/empty-room-problem/'
                title='How to solve the Empty Room Problem'/>
              <ProductResource
                category='blogpost'
                link={reasonHref}
                title='7 Reasons Why Your App Needs an Activity Feed'/>
            </div>
          </div>
        </div>
      </div>
    )
  }
}
