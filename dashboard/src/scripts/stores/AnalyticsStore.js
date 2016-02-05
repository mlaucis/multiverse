import moment from 'moment'

import AnalyticsConstants from '../constants/AnalyticsConstants'

import AccountStore from './AccountStore'

import { EventStore } from '../utils/StoreUtils'
import { register } from '../dispatcher/ConsoleDispatcher'

const FORMAT = AnalyticsConstants.TIME_FORMAT_RANGE

let labels = []
let sets = {}

class AnalyticsStore extends EventStore {
  constructor() {
    super()

    this.dispatchToken = register(this.handleAction)
  }

  get dimensions() {
    return sets
  }

  get labels() {
    return labels
  }

  handleAction = (action) => {
    switch (action.type) {
      case AnalyticsConstants.ANALYTICS_METRICS_SUCCESS:
        labels = generateLabels(action.start, action.end)
        sets = generateSets(labels, action.response)

        this.emitChange()
        break
      default:
      // nothing to do
    }
  };
}

export default new AnalyticsStore

function generateLabels(s, e) {
  let days = moment(e, FORMAT).diff(moment(s, FORMAT), 'days') + 1
  let labels = []

  for (let i = 0; i < days; i++) {
    let label = moment(s, FORMAT).add(i, 'days').format(FORMAT)

    labels.push(label)
  }

  return labels
}

function generateSets(labels, data) {
  let sets = {}
  let keys = Object.keys(data)

  keys.unshift(keys.pop())

  keys.forEach(d => {
    sets[d] = labels.map(l => {
      let ts = data[d]

      for (let i = 0; i < ts.length; i++) {
        if (l == ts[i].bucket) {
          return ts[i].value
        }
      }

      return 0
    })
  })

  return sets
}
