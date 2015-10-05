import React, { Component, PropTypes, findDOMNode } from 'react'

import Clipboard from 'clipboard/dist/clipboard'

import AccountStore from '../stores/AccountStore'
import ApplicationStore from '../stores/ApplicationStore'
import { requestApps } from '../actions/ConsoleActionCreator'
import { requestAppCreate } from '../actions/ConsoleActionCreator'
import { requestAppDelete } from '../actions/ConsoleActionCreator'
import { requestAppUpdate } from '../actions/ConsoleActionCreator'

import CopyIcon from '../../icons/Apps_Icon_CopytoClipboard.svg?t=custom'

class AppForm extends Component {
  static propTypes = {
    description: PropTypes.string,
    name: PropTypes.string,
    onClose: PropTypes.func,
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

    let name = findDOMNode(this.refs.name).value
    let description = findDOMNode(this.refs.description).value

    this.props.onSubmit(name, description)
  }

  render() {
    let c = this.props.submitClass ?
      `btn-${this.props.submitClass}` : 'btn-default'
    let actions = (
      <div className='actions'>
        {( () => {
          if (this.props.onClose) {
            return (
              <button
                className='btn-secondary'
                type='button'
                onClick={this.handleClose}>
                Close
              </button>
            )
          }
        })()}
        <input
          className={c}
          type='submit'
          value={this.props.submitLabel}/>
      </div>
    )

    return (
      <form onSubmit={this.handleSubmit}>
        <div className='group'>
          <input
            defaultValue={this.props.name}
            id='app-name'
            placeholder='Name'
            ref='name'
            required
            type='text'/>
          <span className='bar'></span>
          <span className='help'>Between 2 and 40 characters</span>
          <label
            htmlFor='app-name'>
            Name
          </label>
        </div>
        <div className='group'>
          <input
            defaultValue={this.props.description}
            id='app-description'
            maxLength='100'
            placeholder='Description'
            ref='description'
            type='text'/>
          <span className='bar'></span>
          <span className='help'>Between 2 and 40 characters</span>
          <label
            htmlFor='app-description'>
            Description
          </label>
        </div>
        {actions}
      </form>
    )
  }
}

export class App extends Component {
  static propTypes = {
    actions: PropTypes.bool,
    app: PropTypes.object.isRequired
  }

  constructor() {
    super()

    this._showActions = false
    this._showDelete = false
    this._showEdit = false
    this._showToken = false
    this.state = this.getState()
  }

  componentDidUpdate() {
    if (!this.refs.token) {
      return
    }

    let token = findDOMNode(this.refs.token)

    token.setSelectionRange(0, token.value.length)
  }

  getState() {
    return {
      showActions: this._showActions,
      showDefault: !this._showDelete && !this._showEdit && !this._showToken,
      showDelete: this._showDelete,
      showEdit: this._showEdit,
      showToken: this._showToken
    }
  }

  handleDelete = (event) => {
    event.preventDefault()

    requestAppDelete(this.props.app.id, AccountStore.user)
  }

  handleEdit = (name, description) => {
    requestAppUpdate(this.props.app.id, name, description, AccountStore.user)
      .then(this.toggleEdit)
  }

  toggleActions = (event) => {
    event.preventDefault()

    if (this.state.showDefault) {
      switch (event.type) {
        case 'mouseenter':
          this._showActions = true

          break
        case 'mouseleave':
          this._showActions = false

          break
        default:
          // nothing to do
      }

      this.setState(this.getState())
    }
  }

  toggleDelete = (event) => {
    event.preventDefault()

    this._showDelete = !this._showDelete
    this.setState(this.getState())
  }

  toggleEdit = (event) => {
    if (event && event.preventDefault) {
      event.preventDefault()
    }

    this._showEdit = !this._showEdit
    this.setState(this.getState())
  }

  viewDefault() {
    let app = this.props.app
    let actionsClass = 'actions'
    if (!this.state.showActions) { actionsClass += ' hide' }
    let actions = (
      <div className='grid__col-3'>
        <div className={actionsClass}>
          <button
            className='btn-default outline small'
            onClick={this.toggleEdit}>
            <span className='glyphicon glyphicon-pencil'></span>
          </button>
          <button
            className='btn-default outline small'
            onClick={this.toggleDelete}>
            <span className='glyphicon glyphicon-trash'></span>
          </button>
        </div>
      </div>
    )

    return (
      <div
        className='app card grid__col-md-6'
        onMouseEnter={this.toggleActions}
        onMouseLeave={this.toggleActions}>
        <header className='grid grid--bleed'>
          <h2 className='grid__col-9'>
            {app.name}
          </h2>
          {(() => {
            if (this.props.actions) {
              return actions
            }
          })()}
        </header>
        <main>
          <div className='description'>
            <p>{app.description}</p>
          </div>
          <div>
            <h3>API Token: </h3>
            <button
              className='btn-default block outline copy'
              data-clipboard-text={app.token}>
              {CopyIcon}
              <span className='sub'>{app.token}</span>
            </button>
          </div>
        </main>
      </div>
    )
  }

  viewDelete() {
    return (
      <div className='grid__col-md-6 note note-block note-alert'>
        <h2
          className='alert-heading'>
          Do you really want to delete this app?
        </h2>
        <p>
          This will remove all the data assoicated
          with <strong>{this.props.app.name}</strong> and the operation is
          irreversible.
        </p>
        <div className='actions'>
          <button
            className='btn-secondary'
            href='#'
            onClick={this.toggleDelete}>
            Abort
          </button>
          <button
            className='btn-alert'
            href='#'
            onClick={this.handleDelete}>
            Delete
          </button>
        </div>
      </div>
    )
  }

  viewEdit() {
    return (
      <div className='grid__col-md-6 note note-block note-success'>
        <AppForm
          description={this.props.app.description}
          name={this.props.app.name}
          onClose={this.toggleEdit}
          onSubmit={this.handleEdit}
          submitClass='success'
          submitLabel='Save'/>
      </div>
    )
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

export default class Apps extends Component {
  constructor() {
    super()

    this.state = this.getState()

    this.clipboard = new Clipboard('.copy')
    this.clipboard.on('success', event => {
      let content = event.trigger.innerHTML

      event.trigger.innerHTML = '<span class="sub">Copied</sub>'

      setTimeout(() => {
        event.trigger.innerHTML = content
      }, 2000)
    })
  }

  componentDidMount() {
    ApplicationStore.addChangeListener(this.handleChange)

    requestApps(AccountStore.user)
  }

  componentWillUnmount() {
    ApplicationStore.removeChangeListener(this.handleChange)

    this.clipboard.destroy()
  }

  getState() {
    return {
      apps: ApplicationStore.apps
    }
  }

  handleChange = () => {
    this.setState(this.getState())
  }

  handleCreate = (name, description) => {
    requestAppCreate(name, description, AccountStore.user, true)
  }

  render() {
    let apps = this.state.apps
    let appRows = []
    let createKey = 'app-create'
    let len = Math.round(apps.length / 2)
    let createAppended = false
    let createApp = (
      <div className='card grid__col-sm-6' key={createKey}>
        <AppForm
          onSubmit={this.handleCreate}
          submitLabel='Create'/>
      </div>
    )

    for (let i = 0; i < len; i++) {
      let [ a, b ] = apps
      let rowKey = `app-row-${i}`
      let pair = [(
        <App actions={true} app={a} key={a.id}/>
      )]

      apps.shift()

      if (b) {
        pair.push(<App actions={true} app={b} key={b.id}/>)

        apps.shift()
      } else {
        pair.push(createApp)
        createAppended = true
      }

      appRows.push((
        <div className='grid' key={rowKey}>
          {pair}
        </div>
      ))
    }

    if (apps.length > 0) {
      createKey += `-${apps[apps.length - 1].id}`
    }

    return (
      <div className='apps'>
        {appRows}
        {( () => {
          if (!createAppended) {
            return (
              <div className='grid'>
                {createApp}
              </div>
            )
          }
        })()}
      </div>
    )
  }
}
