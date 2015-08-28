jest.dontMock('../OnboardingStore')

import OnboardingConstants from '../../constants/OnboardingConstants'

describe('OnboardingStore', () => {
  let OnboardingStore
  let ConsoleDispatcher
  let callback
  let persona = require('../OnboardingStore').personas.marketing

  beforeEach(() => {
    OnboardingStore = require('../OnboardingStore')
    ConsoleDispatcher = require('../../dispatcher/ConsoleDispatcher')
    callback = ConsoleDispatcher.register.mock.calls[0][0]

    localStorage.clear()
  })

  it('registers a callback with the dispatcher', () => {
    expect(ConsoleDispatcher.register).toBeCalled()
  })

  it('initializes without state', () => {
    expect(OnboardingStore.options).toBeUndefined()
    expect(OnboardingStore.persona).toBeUndefined()
  })

  it('returns all available personas', () => {
    expect(OnboardingStore.personas.marketing).toBeDefined()
    expect(OnboardingStore.personas.product).toBeDefined()
    expect(OnboardingStore.personas.techlead).toBeDefined()
    expect(OnboardingStore.personas.mobiledev).toBeDefined()
  })

  it('updates its persona', () => {
    callback({
      type: OnboardingConstants.SELECT_PERSONA,
      persona: persona
    })

    expect(OnboardingStore.persona).toEqual(persona)
  })

  it('updates its options', () => {
    let options = [ 'ios', 'others' ]

    callback({
      type: OnboardingConstants.SELECT_OPTIONS,
      options: options
    })

    expect(OnboardingStore.options).toEqual(options)
  })
})
