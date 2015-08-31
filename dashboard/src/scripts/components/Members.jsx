import React, { Component, PropTypes, findDOMNode } from 'react'

import AccountStore from '../stores/AccountStore'
import MemberStore from '../stores/MemberStore'
import { requestMembers } from '../actions/ConsoleActionCreator'
import { requestMemberCreate } from '../actions/ConsoleActionCreator'
import { requestMemberDelete } from '../actions/ConsoleActionCreator'
import { requestMemberUpdate } from '../actions/ConsoleActionCreator'

class MemberForm extends Component {
  static propTypes = {
    email: PropTypes.string,
    firstName: PropTypes.string,
    lastName: PropTypes.string,
    errors: PropTypes.array,
    onClose: PropTypes.func.isRequired,
    onSubmit: PropTypes.func.isRequired,
    submitClass: PropTypes.string,
    submitLabel: PropTypes.string.isRequired
  }

  handleClose = (event) => {
    event.preventDefault()

    this.props.onClose()
  }

  handleSubmit = (event) => {
    event.preventDefault()

    let vals = {}

    for (let name in this.refs) {
      vals[name] = findDOMNode(this.refs[name]).value
    }

    this.props.onSubmit(vals)
  }

  render() {
    let submitClass = this.props.submitClass ?
      `btn-${this.props.submitClass}` : 'btn-default'
    let errors = this.props.errors.map( error => {
      return (
        <div className='alert alert-danger' key={error.code}>
          <strong>{error.code}: </strong>
          {error.message}
        </div>
      )
    })

    return (
      <form onSubmit={this.handleSubmit}>
        {errors}
        <div className='group'>
          <input
            defaultValue={this.props.email}
            id='meber-email'
            placeholder='Email'
            ref='email'
            type='email'/>
          <span className='bar'></span>
          <span className='help'>
            A valid email address in the form of member@org.com
          </span>
          <label htmlFor='member-email'>Email</label>
        </div>
        <div className='group'>
          <input
            defaultValue={this.props.firstName}
            id='meber-firstname'
            placeholder='First Name'
            ref='firstName'
            type='text'/>
          <span className='bar'></span>
          <span className='help'>The first name of the Member</span>
          <label htmlFor='member-firstname'>First Name</label>
        </div>
        <div className='group'>
          <input
            defaultValue={this.props.lastName}
            id='meber-lastname'
            placeholder='Last Name'
            ref='lastName'
            type='text'/>
          <span className='bar'></span>
          <span className='help'>The last name of the Member</span>
          <label htmlFor='member-lastname'>Last Name</label>
        </div>
        <div className='group'>
          <input
            id='meber-password'
            placeholder='Password'
            ref='password'
            type='password'/>
          <span className='bar'></span>
          <span className='help'>Members password to login</span>
          <label htmlFor='member-password'>Password</label>
        </div>
        <div className='actions'>
          <button
            className='btn-secondary'
            onClick={this.handleClose}
            type='button'>
            Close
          </button>
          <input
            className={submitClass}
            type='submit'
            value={this.props.submitLabel}/>
        </div>
      </form>
    )
  }
}

class Member extends Component {
  static propTypes = {
    member: PropTypes.object.isRequired
  }

  constructor() {
    super()

    this._showDelete = false
    this._showEdit = false
    this.state = this.getState()
  }

  getState() {
    return {
      showDelete: this._showDelete,
      showEdit: this._showEdit
    }
  }

  handleDelete = (event) => {
    event.preventDefault()

    requestMemberDelete(this.props.member.id, AccountStore.user)
  }

  handleEdit = (vals) => {
    requestMemberUpdate(vals, this.props.member.id, AccountStore.account.id)
  }

  viewDefault() {
    let member = this.props.member
    let status = member.enabled ? (
      <button className='btn-success small' disabled>
        <span className='glyphicon glyphicon-ok'></span>
      </button>
    ) : (
      <button className='btn-alert small' disabled>
        <span className='glyphicon glyphicon-remove'></span>
      </button>
    )

    return (
      <tr key={member.id}>
        <td className='status'>{status}</td>
        <td><strong>{member.email}</strong></td>
        <td>{member.firstName}</td>
        <td>{member.lastName}</td>
        <td className='actions'>
          <button
            className='btn-default outline small'
            onClick={this.toggleEdit}>
            Edit
          </button>
          <button
            className='btn-default outline small'
            onClick={this.toggleDelete}>
            Delete
          </button>
        </td>
      </tr>
    )
  }

  viewDelete() {
    let member = this.props.member

    return (
      <tr key={member.id}>
        <td className='note note-block note-alert' colSpan='5'>
          <h2
            className='alert-heading'>
            Do you really want to delete this member?
          </h2>
          <p>
            This will remove all the data assoicated with
            <strong>{this.props.member.email}</strong> and the operation is
            irreversible.
          </p>
          <div className='actions'>
            <button
              className='btn-secondary'
              onClick={this.toggleDelete}>
              Abort
            </button>
            <button
              className='btn-alert'
              onClick={this.handleDelete}>
              Delete
            </button>
          </div>
        </td>
      </tr>
    )
  }

  viewEdit() {
    let member = this.props.member

    return (
      <tr key={member.id}>
        <td className='note note-block note-success' colSpan='5'>
          <MemberForm
            email={member.email}
            errors={[]}
            firstName={member.firstName}
            lastName={member.lastName}
            onClose={this.toggleEdit}
            onSubmit={this.handleEdit}
            submitClass='success'
            submitLabel='Save'/>
        </td>
      </tr>
    )
  }

  toggleEdit = () => {
    this._showEdit = !this._showEdit
    this.setState(this.getState())
  }

  toggleDelete = () => {
    this._showDelete = !this._showDelete
    this.setState(this.getState())
  }

  render() {
    let view = this.viewDefault()

    if (this.state.showDelete) {
      view = this.viewDelete()
    }

    if (this.state.showEdit) {
      view = this.viewEdit()
    }

    return view
  }
}

export default class Members extends Component {
  constructor() {
    super()

    this._showCreae = false
    this.state = this.getState()
  }

  componentDidMount() {
    MemberStore.addChangeListener(this.handleChange)

    requestMembers(AccountStore.user)
  }

  componentWillUnmount() {
    MemberStore.removeChangeListener(this.handleChange)
  }

  getState() {
    return {
      errors: MemberStore.errors,
      members: MemberStore.members,
      showCreate: this._showCreate
    }
  }

  handleChange = () => {
    this.setState(this.getState())
  }

  handleCreate = (vals) => {
    requestMemberCreate(vals, AccountStore.account.id).then(this.toggleCreate)
  }

  toggleCreate = () => {
    this._showCreate = !this._showCreate

    this.setState(this.getState())
  }

  render() {
    let head = this.state.showCreate ? (
      <MemberForm
        errors={this.state.errors}
        onClose={this.toggleCreate}
        onSubmit={this.handleCreate}
        submitLabel='Create'/>
    ) : (
      <div className='btn-group'>
        <button
          className='btn-default'
          onClick={this.toggleCreate}>
          Create
        </button>
      </div>
    )
    let members = this.state.members.map( member => {
      return <Member key={member.id} member={member}/>
    })

    return (
      <div>
        <div className='row'>
          <div className='col-md-12'>
            <div className='portlet light'>
              <div className='portlet-body'>
                {head}
                <table>
                  <thead>
                    <th>Status</th>
                    <th>Email</th>
                    <th>First Name</th>
                    <th>Last Name</th>
                    <th></th>
                  </thead>
                  <tbody>
                    {members}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
      </div>
    )
  }
}
