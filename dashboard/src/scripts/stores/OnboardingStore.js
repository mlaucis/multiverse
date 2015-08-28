import OnboardingConstants from '../constants/OnboardingConstants'
import { EventStore } from '../utils/StoreUtils'
import { register } from '../dispatcher/ConsoleDispatcher'

const personas = {
  'marketing': {
    avatar: 'avatar.png',
    display: 'Marketing Genius',
    options: {
      'retention': { display: 'Retention' },
      'growth': { display: 'User Growth' },
      'community': { display: 'Community Building' },
      'brand': { display: 'Brand Building' }
    },
    question: 'What do you expect from adding a social layer to your app?'
  },
  'product': {
    avatar: 'avatar.png',
    display: 'Product Star',
    options: {},
    question: ''
  },
  'techlead': {
    avatar: 'avatar.png',
    display: 'Technical Lead',
    options: {
      'web': { display: 'Web' },
      'ios': { display: 'iOS' },
      'android': { display: 'Android' },
      'other': { display: 'Others' }
    },
    question: 'Which platforms are you interested in?'
  },
  'mobiledev': {
    avatar: 'avatar.png',
    name: 'mobiledev',
    display: 'Mobile Developer',
    options: {},
    question: ''
  }
}

class OnboardingStore extends EventStore {
  constructor() {
    super()

    this.dispatchToken = register(this.handleAction)
    this._options = undefined
    this._persona = undefined
  }

  get options() {
    return this._options
  }

  get persona() {
    return this._persona
  }

  get personas() {
    return personas
  }

  handleAction = (action) => {
    switch (action.type) {
      case OnboardingConstants.SELECT_OPTIONS:
        this._options = action.options

        this.emitChange()
        break
      case OnboardingConstants.SELECT_PERSONA:
        this._persona = action.persona

        this.emitChange()
        break
      default:
      // nothing to do
    }
  }
}

export default new OnboardingStore
