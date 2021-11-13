import {expose} from 'comlink'
import {reportStats, clientsList} from 'lib/report-helpers'

const getStatsAndClients = (rules) => (
  [reportStats(rules), clientsList(rules)]
)

expose({getStatsAndClients})
