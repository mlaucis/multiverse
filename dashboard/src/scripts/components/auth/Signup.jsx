import React, { Component, PropTypes } from 'react'
import  { findDOMNode } from 'react-dom'

import { Link } from 'react-router'

import AccountStore from '../../stores/AccountStore'
import { requestAccountCreate } from '../../actions/ConsoleActionCreator'
import { requestAccountUserCreate } from '../../actions/ConsoleActionCreator'
import { requestAppCreate } from '../../actions/ConsoleActionCreator'
import { consumeReferrerCookie } from '../../utils/CookieUtils'

class SignupForm extends Component {
  static propTypes = {
    errors: PropTypes.array,
    submit: PropTypes.func.isRequired
  };

  constructor() {
    super()

    this.state = { isSending: false }
  }

  componentWillReceiveProps(next) {
    if (next.errors.length > 0) {
      this.setState({ isSending: false })
    }
  }

  handleSubmit = (event) => {
    event.preventDefault()

    if (this.state.isSending) {
      return
    }

    this.setState({ isSending: true })

    this.props.submit({
      accountName: findDOMNode(this.refs.accountName).value,
      email: findDOMNode(this.refs.email).value,
      password: findDOMNode(this.refs.password).value,
      firstName: findDOMNode(this.refs.firstName).value,
      lastName: findDOMNode(this.refs.lastName).value
    })
  };

  render() {
    let errors = this.props.errors.map( error => {
      return <p className='error' key={error.code}>{error.message}</p>
    })
    let action = this.state.isSending ? (
      <input
        className='btn-default block'
        disabled
        type='submit'
        value='Sending...'/>
    ) : (
      <input
        className='btn-default block'
        type='submit'
        value='Sign Up'/>
    )

    return (
      <form onSubmit={this.handleSubmit}>
        <div className='form-group'>
          <div className='group errors'>
            {errors}
          </div>
          <div className='grid grid--justify-space-between'>
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
          </div>
          <div className='group'>
            <input id='accountName'
                   name='accountName'
                   placeholder='Company'
                   ref='accountName'
                   type='text'
            />
            <span className='bar'></span>
            <span className='help'></span>
            <label hmtlFor='accountName'>Organisation</label>
          </div>
          <div className='group'>
            <input id='email'
                   name='email'
                   placeholder='Email'
                   ref='email'
                   type='email'
                   required
            />
            <span className='bar'></span>
            <span className='help'></span>
            <label hmtlFor='email'>Email</label>
          </div>
          <div className='group'>
            <input id='password'
                   name='password'
                   placeholder='Password'
                   ref='password'
                   type='password'
                   required
            />
            <span className='bar'></span>
            <span className='help'></span>
            <label hmtlFor='password'>Password</label>
          </div>
          <p>Already have an account? <Link to='AUTH_LOGIN'>Log In</Link></p>
        </div>
        <div className='actions'>
          {action}
        </div>
      </form>
    )
  }
}

export default class Signup extends Component {
  static contextTypes = {
    router: PropTypes.func.isRequired
  };

  static propTypes = {
    query: PropTypes.object
  };

  constructor() {
    super()

    this.state = this.getState()
  }

  getState() {
    return {
      errors: AccountStore.errors
    }
  }

  componentDidMount() {
    AccountStore.addChangeListener(this._onChange)
  }

  componentWillUnmount() {
    AccountStore.removeChangeListener(this._onChange)
  }

  _onChange = () => {
    this.setState(this.getState())
  };

  handleSubmit = (values) => {
    let plan = this.props.query.plan || 'free'
    let router = this.context.router
    let originalReferrer = consumeReferrerCookie()

    if (values.accountName === '') {
      values.accountName = values.email
    }

    requestAccountCreate(values, plan, originalReferrer)
    .then( account => {
      requestAccountUserCreate(
        values,
        account.id,
        originalReferrer,
        document.referrer
      ).then( () => {
        requestAppCreate(
          'Testing Application',
          'This is your first app. Use its API token for testing.',
          AccountStore.user,
          false
        ).then( () => {
          router.transitionTo('DASHBOARD')
        })
      })
    })
  };

  render() {
    return (
      <SignupForm
        errors={this.state.errors}
        submit={this.handleSubmit}/>
    )
  }
}
