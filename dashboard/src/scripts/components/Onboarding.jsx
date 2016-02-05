import React, { Component, PropTypes } from 'react'
import { findDOMNode } from 'react-dom'
import { Link } from 'react-router'

import OnboardingStore from '../stores/OnboardingStore'
import { selectOptions, selectPersona } from '../actions/ConsoleActionCreator'

class OptionsSelect extends Component {
  static propTypes = {
    avatar: PropTypes.string.isRequired,
    nextStep: PropTypes.func.isRequired,
    options: PropTypes.object.isRequired,
    question: PropTypes.string.isRequired
  };

  handleNext = (event) => {
    event.preventDefault()

    let selected = []

    for (let name in this.refs) {
      let ref = this.refs[name]

      if (ref.props.type === 'submit') {
        continue
      }

      if (findDOMNode(ref).checked) {
        selected.push(ref.props.value)
      }
    }

    this.props.nextStep()
    selectOptions(selected)
  };

  render() {
    let avatarURL = `/${this.props.avatar}`
    let os = Object.keys(this.props.options).map( name => {
      let option = this.props.options[name]

      return (
        <label key={name}>
          <input
            ref={name}
            type='checkbox'
            value={option.display}/>
          {option.display}
        </label>
      )
    })

    return (
      <div id='options-select'>
        <h2><img src={avatarURL} width='70'/>{this.props.question}</h2>
        <form onSubmit={this.handleNext}>
        {os}
        <input type='submit' ref='submit' value='next'/>
        </form>
      </div>
    )
  }
}

class PersonaSelect extends Component {
  static propTypes = {
    nextStep: PropTypes.func.isRequired,
    personas: PropTypes.object.isRequired
  };

  handleSelect = (name) => {
    let persona = this.props.personas[name]
    let next = this.props.nextStep

    return function () {
      next()
      selectPersona(persona)
    }
  };

  render() {
    let ps = Object.keys(this.props.personas).map( name => {
      let persona = this.props.personas[name]
      let imgURL = `/${persona.avatar}`

      return (
        <li key={name}>
          <button onClick={this.handleSelect(name)}>
            <img width='140' src={imgURL}/>
            <br/>
            {persona.display}
          </button>
        </li>
      )
    })

    return (
      <div id='persona-select'>
        <h1>Which person describes you best?</h1>
        <ul>
        {ps}
        </ul>
      </div>
    )
  }
}

class Promotion extends Component {
  render() {
    return (
      <h2>Promotion</h2>
    )
  }
}

export default class Onboarding extends Component {
  constructor() {
    super()

    this._step = 1
    this.state = this.getState()
  }

  componentDidMount() {
    OnboardingStore.addChangeListener(this.handleChange)
  }

  componentWillUnmount() {
    OnboardingStore.removeChangeListener(this.handleChange)
  }

  getState() {
    return {
      options: OnboardingStore.options,
      persona: OnboardingStore.persona,
      personas: OnboardingStore.personas,
      step: this._step
    }
  }

  handleChange = () => {
    this.setState(this.getState())
  };

  nextStep = () => {
    this._step += 1
  };

  showStep() {
    switch(this.state.step) {
      case 1:
        return (
          <PersonaSelect
            nextStep={this.nextStep}
            personas={this.state.personas}/>
        )
      case 2:
        let persona = this.state.persona

        return (
          <OptionsSelect
            avatar={persona.avatar}
            nextStep={this.nextStep}
            options={persona.options}
            question={persona.question}/>
        )
      case 3:
        return <Promotion/>
    }
  }

  render() {
    let linkText = (this.state.step < 3) ? 'Skip' : 'Start'

    return (
      <div id='onboarding'>
        {this.showStep()}
        <Link to='DASHBOARD'>{linkText}</Link>
      </div>
    )
  }
}
