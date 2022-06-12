<svelte:options immutable="{true}" />

<script>
  import {onMount, onDestroy, getContext} from 'svelte'
  import {EditorState, EditorSelection} from '@codemirror/state'
  import {
    EditorView,
    keymap,
    lineNumbers,
    highlightActiveLineGutter,
    gutter
  } from '@codemirror/view'
  import {defaultKeymap, history, historyKeymap} from '@codemirror/commands'
  import {defaultHighlightStyle, syntaxHighlighting} from '@codemirror/language'
  import {html} from '@codemirror/lang-html'
  import {oneDarkTheme, oneDarkHighlightStyle} from 'lib/codemirrorDarkTheme'
  import {isDarkThemeON} from 'stores/theme'
  import {report, reportLoading, reportError, linesAndSelectors} from 'stores/report'
  import {inlinerLoading, inlinerError} from 'stores/inliner'
  import {splitState} from 'stores/split'
  import {
    validationErrorsMarker,
    validationErrorsEffect,
    validationErrorsState
  } from 'lib/codemirrorValidationErrors'
  import {
    EVENT_LINE_TO_EDITOR,
    EVENT_LINE_TO_REPORT,
    EVENT_SUBMIT_EXAMPLE,
    EVENT_INLINE_CSS
  } from 'lib/constants'
  import {loadSampleContent} from 'lib/sampleHelpers'
  import {getTooltipText} from 'lib/reportHelpers'

  const TOOLTIP_SHIFT_PX = 20

  let editorElement
  let editorView = null
  let tooltipElement = null
  let tooltipTextElement = null

  const {getWebWorker} = getContext('ww')

  const getEditorState = (doc = '') => {
    const [eTheme, eThemeHighLight] = (() => {
      if ($isDarkThemeON) {
        return [oneDarkTheme, syntaxHighlighting(oneDarkHighlightStyle)]
      }

      return [
        EditorView.baseTheme({
          '&.cm-editor': {
            fontSize: '0.9rem',
            height: '100%'
          }
        }),
        syntaxHighlighting(defaultHighlightStyle)
      ]
    })()

    const showTooltip = (view, edLine, event) => {
      if (!tooltipElement || !tooltipTextElement) {
        return
      }

      const line = view.state.doc.lineAt(edLine.from)
      if (
        line?.number !== null &&
        line.number >= 0 &&
        $linesAndSelectors[line.number] &&
        $linesAndSelectors[line.number].length > 0
      ) {
        tooltipTextElement.textContent = getTooltipText($linesAndSelectors[line.number])
        tooltipElement.style.top = `${event.clientY + TOOLTIP_SHIFT_PX}px`
        tooltipElement.style.left = `${event.clientX + TOOLTIP_SHIFT_PX}px`
        tooltipElement.style.opacity = '1'
        tooltipElement.style.display = 'block'
      }
    }

    const hideTooltip = () => {
      if (!tooltipElement || !tooltipTextElement) {
        return
      }

      tooltipElement.style.display = 'none'
      tooltipElement.style.opacity = '0'
      tooltipTextElement.textContent = ''
    }

    const clickMarker = (view, edLine) => {
      const line = view.state.doc.lineAt(edLine.from)
      if (line.number !== null && line?.number >= 0) {
        splitState.switchToRightOnMobile()
        window.dispatchEvent(
          new window.CustomEvent(EVENT_LINE_TO_REPORT, {
            detail: {line: line.number}
          })
        )
      }
    }

    return EditorState.create({
      doc,
      extensions: [
        validationErrorsState,
        gutter({
          class: 'validation-error-gutter',
          markers: (v) => v.state.field(validationErrorsState),
          initialSpacer: () => validationErrorsMarker,
          domEventHandlers: {
            mouseover(view, edLine, event) {
              showTooltip(view, edLine, event)
              return true
            },
            mousemove(_view, _edLine, event) {
              if (!tooltipElement || !tooltipTextElement) {
                return true
              }

              tooltipElement.style.top = `${event.clientY + TOOLTIP_SHIFT_PX}px`
              tooltipElement.style.left = `${event.clientX + TOOLTIP_SHIFT_PX}px`
              return true
            },
            mouseout() {
              hideTooltip()
              return true
            },
            click(view, edLine) {
              clickMarker(view, edLine)
              return true
            },
            keyup(view, edLine, event) {
              // not working - edLine not changing :(
              if (event.keyCode === 13 || event.key === 'Enter') {
                event.preventDefault()
                clickMarker(view, edLine)
                return true
              }
              return false
            }
          }
        }),
        lineNumbers(),
        highlightActiveLineGutter(),
        history(),
        keymap.of([...defaultKeymap, ...historyKeymap]),
        html(),
        eThemeHighLight,
        eTheme
      ]
    })
  }

  const runHtmlAnalyze = () => {
    const editorContent = editorView.state.doc.toString()
    if (!editorContent || !editorContent.length) {
      return Promise.resolve()
    }

    return new Promise((resolve, reject) => {
      splitState.switchToRightOnMobile()
      const webWorker = getWebWorker()
      if (webWorker?.processHTML) {
        webWorker
          .processHTML(editorContent)
          .then((data) => {
            report.set(data)
            resolve()
          })
          .catch(reject)
      } else {
        reject(new Error('Web Worker is not available'))
      }
    })
  }

  const handleReportLineClickEvent = (e) => {
    if (!editorView || !e.detail?.line) {
      return
    }

    const {line} = e.detail
    const editorLine = editorView.state.doc.line(line)
    const selection = EditorSelection.cursor(editorLine.from)
    editorView.dispatch({
      selection,
      effects: EditorView.scrollIntoView(selection, {
        y: 'center'
      })
    })
    editorView.focus()
  }

  const handleInlineCss = () => {
    const editorContent = editorView.state.doc.toString()
    if (!editorContent || !editorContent.length) {
      return
    }

    inlinerError.set(null)
    inlinerLoading.set(true)

    const webWorker = getWebWorker()
    if (webWorker?.inlineCSS) {
      webWorker
        .inlineCSS(editorContent)
        .then((data) => {
          const currentEditorValue = editorView.state.doc.toString()
          const endPosition = currentEditorValue.length

          editorView.dispatch({
            changes: {
              from: 0,
              to: endPosition,
              insert: data
            }
          })
          setTimeout(() => runHtmlAnalyze(), 0)
        })
        .catch((err) => {
          inlinerError.set(err)
        })
        .finally(() => {
          inlinerLoading.set(false)
        })
    } else {
      inlinerError.set(new Error('Web Worker is not available'))
      inlinerLoading.set(false)
    }
  }

  const handleLoadSample = () => {
    reportError.set(null)
    reportLoading.set(true)

    loadSampleContent()
      .then((sampleContent) => {
        if (!sampleContent) {
          return Promise.resolve()
        }

        const currentEditorValue = editorView.state.doc.toString()
        const endPosition = currentEditorValue.length

        editorView.dispatch({
          changes: {
            from: 0,
            to: endPosition,
            insert: sampleContent
          }
        })

        return runHtmlAnalyze()
      })
      .catch((err) => reportError.set(err))
      .finally(() => reportLoading.set(false))
  }

  const onSubmitHtml = () => {
    reportError.set(null)
    reportLoading.set(true)

    runHtmlAnalyze()
      .catch((err) => reportError.set(err))
      .finally(() => reportLoading.set(false))
  }

  const applyErrorGutters = (linesSelector) => {
    if (!editorView) {
      return
    }

    editorView.dispatch({
      effects: validationErrorsEffect.of({type: 'empty'})
    })
    Object.keys(linesSelector).forEach((line) => {
      const editorLine = editorView.state.doc.line(line)
      if (editorLine?.from !== null && editorLine.from >= 0) {
        editorView.dispatch({
          effects: validationErrorsEffect.of({
            pos: editorLine.from,
            selector: linesSelector[line]
          })
        })
      }
    })
  }

  const createEditor = (doc = '') => {
    editorView = new EditorView({
      state: getEditorState(doc),
      parent: editorElement
    })
  }

  const destroyEditor = () => {
    if (editorView) {
      editorView.destroy()
      editorView = null
    }
  }

  const unsubscribeLinesReport = linesAndSelectors.subscribe((linesSelector) => {
    applyErrorGutters(linesSelector)
  })

  const unsubscribeIsDarkTheme = isDarkThemeON.subscribe(() => {
    if (!editorView) {
      return
    }

    const htmlContent = editorView.state.doc.toString()
    destroyEditor()
    createEditor(htmlContent)
    applyErrorGutters($linesAndSelectors)
  })

  onMount(() => {
    createEditor()
    return destroyEditor
  })

  onMount(() => {
    window.addEventListener(EVENT_LINE_TO_EDITOR, handleReportLineClickEvent)
    window.addEventListener(EVENT_INLINE_CSS, handleInlineCss)
    window.addEventListener(EVENT_SUBMIT_EXAMPLE, handleLoadSample)
    return () => {
      window.removeEventListener(EVENT_LINE_TO_EDITOR, handleReportLineClickEvent)
      window.removeEventListener(EVENT_INLINE_CSS, handleInlineCss)
      window.removeEventListener(EVENT_SUBMIT_EXAMPLE, handleLoadSample)
    }
  })

  onDestroy(unsubscribeLinesReport)
  onDestroy(unsubscribeIsDarkTheme)
</script>

<style>
  .editor-area {
    flex-grow: 1;
    position: relative;
  }

  .editor-area-edit {
    height: 100%;
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    overflow: scroll;
  }

  .editor-footer {
    background-color: var(--bgColor);
    flex-shrink: 0;
    overflow: hidden;
    white-space: nowrap;
    padding: 0.3rem 0;
    box-shadow: 0 -10px 30px 0 rgb(82 63 105 / 8%);
  }

  .editor-area-btn {
    color: var(--buttonColor);
    display: inline-block;
    background-color: var(--buttonBgColor);
    line-height: 1.25;
    text-align: center;
    white-space: nowrap;
    vertical-align: middle;
    cursor: pointer;
    user-select: none;
    border: 1px solid transparent;
    padding: 0.8rem 1rem;
    font-size: 1rem;
    width: 100%;
  }

  .editor-area-btn:hover {
    background-color: var(--buttonBgHoverColor);
  }

  .editor-area-btn:active {
    background-color: var(--buttonBgActiveColor);
  }

  .editor-tooltip {
    display: none;
    opacity: 0;
    position: fixed;
    z-index: 100;
    max-width: 25rem;
    white-space: pre-wrap;
    overflow: hidden;
    border-radius: 0.4rem;
    padding: 0.4rem;
    box-shadow: 0 3px 18px rgb(0 0 0 / 0%);
    border: 2px solid var(--splitBorderColor);
    color: var(--baseColor);
    background-color: var(--cardBgColor);
    font-size: 0.8rem;
  }

  .editor-tooltip-message {
    background-position: top left;
    background-repeat: no-repeat;
  }
</style>

<div class="editor-area">
  <div class="editor-area-edit" bind:this="{editorElement}"></div>
</div>
<div class="editor-footer">
  <button class="editor-area-btn" on:click|preventDefault="{onSubmitHtml}"
    >Check email HTML and CSS</button
  >
</div>
<!-- tooltip -->
<div class="editor-tooltip" bind:this="{tooltipElement}">
  <div class="editor-tooltip-message" bind:this="{tooltipTextElement}"></div>
</div>
