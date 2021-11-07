<script>
  import {onMount, onDestroy} from 'svelte'
  import {EditorState, EditorSelection} from '@codemirror/state'
	import {EditorView, keymap} from '@codemirror/view'
	import {defaultKeymap} from '@codemirror/commands'
	import {history, historyKeymap} from '@codemirror/history'
	import {lineNumbers, highlightActiveLineGutter, gutter} from '@codemirror/gutter'
	import {defaultHighlightStyle} from '@codemirror/highlight'
	import {html} from '@codemirror/lang-html'
	import {report, linesAndSelectors} from 'stores/report'
	import {
		validationErrorsMarker,
		validationErrorsEffect,
		validationErrorsState
	} from 'lib/coremirror-validation-errors'
  import {EVENT_LINE_TO_EDITOR, EVENT_LINE_TO_REPORT} from 'lib/constants'

	export let parserFunction

  let editorElement
  let editorView = null

  const eTheme = EditorView.baseTheme({
		'&.cm-editor': {
			fontSize: '0.9rem',
			height: '100%'
		}
	})

	const initialEditorState = EditorState.create({
		doc: '',
		extensions: [
			validationErrorsState,
			gutter({
				class: 'validation-error-gutter',
				markers: (v) => v.state.field(validationErrorsState),
				initialSpacer: () => validationErrorsMarker,
				domEventHandlers: {
					mouseover(view, edLine) {
						const line = view.state.doc.lineAt(edLine.from)
            if (line?.number && $linesAndSelectors[line.number]) {
							// tooltip?
							// console.log('Info', $linesAndSelectors[line.number])
						}
						return true
					},
					mousedown(view, edLine) {
            const line = view.state.doc.lineAt(edLine.from)
            if (line?.number) {
              window.dispatchEvent(
                new window.CustomEvent(
                  EVENT_LINE_TO_REPORT,
                  { detail: {line: line.number} }
                )
              )
            }
            return true
					}
				}
			}),
			lineNumbers(),
			highlightActiveLineGutter(),
			history(),
			defaultHighlightStyle,
			keymap.of([
				...defaultKeymap,
				...historyKeymap
			]),
			html(),
			eTheme
		]
	})

  const handleReportLineClickEvent = (e) => {
    if (!editorView || !e.detail?.line) {
      return
    }

    const {line} = e.detail
    const editorLine = editorView.state.doc.line(line)
		const selection = EditorSelection.cursor(editorLine.from)
		editorView.dispatch({effects: EditorView.centerOn.of(selection), selection})
		editorView.focus()
  }

  const onSubmitHtml = async () => {
		const html = editorView.state.doc.toString()

		try {
			const reportData = await parserFunction(html)
			report.update(reportData)
		} catch (err) {
			console.log('error', err)
		}
	}

  const unsubscribeLinesReport = linesAndSelectors.subscribe((linesSelector) => {
    if (!editorView) {
      return
    }

	  editorView.dispatch({
			effects: validationErrorsEffect.of({type: 'empty'})
		})
		Object.keys(linesSelector).forEach((line) => {
			const editorLine = editorView.state.doc.line(line)
      if (editorLine?.from) {
        editorView.dispatch({
          effects: validationErrorsEffect.of({pos: editorLine.from, selector: linesSelector[line]})
        })
      }
		})
  })

  onMount(() => {
		editorView = new EditorView({
			state: initialEditorState,
			parent: editorElement
		})

		return () => {
			if (editorView) {
				editorView.destroy()
				editorView = null
			}
		}
	})

  onMount(() => {
    window.addEventListener(EVENT_LINE_TO_EDITOR, handleReportLineClickEvent)
    return () => window.removeEventListener(EVENT_LINE_TO_EDITOR, handleReportLineClickEvent)
  })

  onDestroy(unsubscribeLinesReport)
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
		flex-shrink: 0;
		overflow: hidden;
    white-space: nowrap;
	}
</style>

<div class="editor-area">
  <div class="editor-area-edit" bind:this={editorElement}></div>
</div>
<div class="editor-footer">
  <button on:click|preventDefault={onSubmitHtml}>Check HTML</button>
</div>
