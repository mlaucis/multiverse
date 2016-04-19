import React, { Component, PropTypes } from 'react'
import { findDOMNode } from 'react-dom'

import AccountStore from '../stores/AccountStore'
import ApplicationStore from '../stores/ApplicationStore'
import { requestApps } from '../actions/ConsoleActionCreator'
import { requestMemberInvite } from '../actions/ConsoleActionCreator'
import { setReadmeToken } from '../utils/CookieUtils'

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
    let firstName = findDOMNode(this.refs.firstName).value
    let lastName = findDOMNode(this.refs.lastName).value
    let toggle = this.toggleSuccess

    requestMemberInvite(email, firstName, lastName).then( () => {
      toggle()
    })
  };

  viewDefault() {
    return (
      <div className='card invite-developer grid__col-sm-6'>
        <div className='grid__col-12'>
          <h3>Get your team on board</h3>
          <p>Invite a developer or others who should help integrate Tapglue
            into your app.</p>
        </div>
        <form
          className='grid'
          onSubmit={this.handleSubmit}>
          <div className='group grid__col-sm-6'>
            <input id='firstName'
                   name='firstName'
                   placeholder='First Name'
                   ref='firstName'
                   type='text'
                   required
            />
            <span className='bar'></span>
            <span className='help'></span>
            <label hmtlFor='firstName'>First Name</label>
          </div>
          <div className='group grid__col-sm-6'>
            <input id='lastName'
                   name='lastName'
                   placeholder='Last Name'
                   ref='lastName'
                   type='text'
                   required
            />
            <span className='bar'></span>
            <span className='help'></span>
            <label hmtlFor='lastName'>Last Name</label>
          </div>
          <div className='group grid__col-12'>
            <input
              id='developer-email'
              placeholder='Email Address'
              ref='email'
              required
              type='email'/>
            <span className='bar'></span>
            <span className='help'>
              Valid email in the form of member@company.org
            </span>
            <label htmlFor='developer-email'>Email Address</label>
          </div>
          <div className='grid__col-12'>
            <input className='btn-default block' type='submit' value='Invite'/>
          </div>
        </form>
      </div>
    )
  }

  viewSuccess() {
    return (
      <div className='card grid__col-sm-6 note note-success '>
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
  };

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
  };

  render() {
    let iconClass = `fa fa-${this.props.icon}`

    return (
      <a
        className='btn-secondary outline resource grid__col-xs-2'
        href={this.props.link}>
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
  };

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
  };

  componentDidMount() {
    setReadmeToken(this.props.app.token)
  }

  render() {
    return (
      <div className='grid testing-app'>
        <div className='card grid__col-md-6'>
          <h2>Instant API Access</h2>
          <p>{this.props.app.description}</p>
        </div>
        <div className='card grid__col-md-6'>
          <div className='grid grid--align-center'>
            <div className='grid__col-lg-3'><p>API TOKEN:</p></div>
            <div className='token grid__col-lg-9'>{this.props.app.token}</div>
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
  };

  handleClick = (event) => {
    event.preventDefault()

    this.setState(this.getState())
  };
  render() {
    let reasonHref = 'http://www.tapglue.com/blog/' +
      'why-your-mobile-app-needs-a-social-activity-feed/'
    let app = this.state.app ? (
      <TestingApp app={this.state.app}/>
    ) : (
      <div></div>
    )

    return (
      <div className='home grid grid--align-content-start'>
        <div className='teaser grid__col-12'>
          <h1>Welcome to Tapglue, {this.state.user.firstName}!</h1>
          <p>Integrating Tapglue into your app is a matter of a few hours.</p>
          <p>Jump straight into it with these helpful resources.</p>
          <div className='resources grid grid--justify-center'>
            <IntegrationResource
              icon='file'
              link='https://developers.tapglue.com'
              name='Docs'/>
            <IntegrationResource
              icon='apple'
              link='https://developers.tapglue.com/docs/ios'
              name='iOS Guide'/>
            <IntegrationResource
              icon='android'
              link='https://developers.tapglue.com/docs/android'
              name='Android Guide'/>
          </div>
        </div>
        {app}
        <div className='grid'>
          <InviteDeveloper/>
          <div className='grid__col-md-6'>
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
