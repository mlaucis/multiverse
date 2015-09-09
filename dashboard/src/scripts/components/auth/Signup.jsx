import React, { Component, PropTypes, findDOMNode } from 'react'

import { Link } from 'react-router'

import AccountStore from '../../stores/AccountStore'
import { requestAccount } from '../../actions/ConsoleActionCreator'
import { requestAccountCreate } from '../../actions/ConsoleActionCreator'
import { requestAccountUserCreate } from '../../actions/ConsoleActionCreator'
import { requestAppCreate } from '../../actions/ConsoleActionCreator'
import { requestLogin } from '../../actions/ConsoleActionCreator'
import { consumeReferrerCookie } from '../../utils/AuthUtils'

class SignupForm extends Component {
  static propTypes = {
    errors: PropTypes.array,
    submit: PropTypes.func.isRequired
  }

  handleSubmit = (event) => {
    event.preventDefault()

    this.props.submit({
      accountName: findDOMNode(this.refs.accountName).value,
      email: findDOMNode(this.refs.email).value,
      password: findDOMNode(this.refs.password).value,
      firstName: findDOMNode(this.refs.firstName).value,
      lastName: findDOMNode(this.refs.lastName).value
    })
  }

  render() {
    let errors = this.props.errors.map( error => {
      return <p key={error.code}>{error.message}</p>
    })

    return (
      <form onSubmit={this.handleSubmit}>
        {errors}
        <div className='form-group'>
          <div className='group'>
            <input id='accountName'
                   name='accountName'
                   placeholder='Organisation'
                   ref='accountName'
                   type='text'
                   required
            />
            <span className='bar'></span>
            <span className='help'></span>
            <label hmtlFor='accountName'>Organisation</label>
          </div>
          <div className='group'>
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
          <div className='group'>
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
          <input
            className='btn-default block'
            type='submit'
            value='Sign Up'/>
        </div>
      </form>
    )
  }
}

export default class Signup extends Component {
  static contextTypes = {
    router: PropTypes.func.isRequired
  }

  static propTypes = {
    query: PropTypes.object
  }

  constructor() {
    super()

    this.state = this.getState()
  }

  getState() {
    return {
      errors: AccountStore.errors
    }
  }

  handleSubmit = (values) => {
    let plan = this.props.query.plan || 'free'
    let router = this.context.router
    let originalReferrer = consumeReferrerCookie()

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
  }

  componentDidMount() {
    AccountStore.addChangeListener(this._onChange)
  }

  componentWillUnmount() {
    AccountStore.removeChangeListener(this._onChange)
  }

  _onChange = () => {
    this.setState(this.getState())
  }

  render() {
    return (
      <section className='signup'>
        <SignupForm errors={this.state.errors} submit={this.handleSubmit}/>
      </section>
    )
  }
}
