import './init'
import './hotwire'

import {onDomReady} from 'utils/dom'
import {getTheme, activateTheme} from 'utils/theme'

onDomReady(() => activateTheme(getTheme()))
