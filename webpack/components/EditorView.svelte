<script>
  import {onMount, onDestroy} from 'svelte'
  import {EditorState, EditorSelection} from '@codemirror/state'
	import {EditorView, keymap} from '@codemirror/view'
	import {defaultKeymap} from '@codemirror/commands'
	import {history, historyKeymap} from '@codemirror/history'
	import {lineNumbers, highlightActiveLineGutter, gutter} from '@codemirror/gutter'
	import {defaultHighlightStyle} from '@codemirror/highlight'
	import {html} from '@codemirror/lang-html'
	import {oneDarkTheme, oneDarkHighlightStyle} from 'lib/codemirrorDarkTheme'
	import {isDarkThemeON} from 'stores/theme'
	import {report, reportLoading, reportError, linesAndSelectors} from 'stores/report'
	import {splitState} from 'stores/split'
	import {
		validationErrorsMarker,
		validationErrorsEffect,
		validationErrorsState
	} from 'lib/codemirrorValidationErrors'
  import {EVENT_LINE_TO_EDITOR, EVENT_LINE_TO_REPORT} from 'lib/constants'

	export let parserFunction

  let editorElement
  let editorView = null

	const getEditorState = (doc = '') => {
		const [eTheme, eThemeHighLight] = (() => {
			if ($isDarkThemeON) {
				return [
					oneDarkTheme,
					oneDarkHighlightStyle
				]
			}

			return [
				EditorView.baseTheme({
					'&.cm-editor': {
						fontSize: '0.9rem',
						height: '100%'
					}
				}),
				defaultHighlightStyle
			]
		})()

		return EditorState.create({
			doc,
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
				eThemeHighLight,
				keymap.of([
					...defaultKeymap,
					...historyKeymap
				]),
				html(),
				eTheme
			]
		})
	}

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
		reportError.set(null)
		reportLoading.set(true)

		try {
			const html = editorView.state.doc.toString()
			const reportData = await parserFunction(html)
			report.set(reportData)
			splitState.switchToRightOnMobile()
		} catch (err) {
			reportError.set(err)
		}

		reportLoading.set(false)
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

	const unsubscribeIsDarkTheme = isDarkThemeON.subscribe(() => {
		if (!editorView) {
			return
		}

		const html = editorView.state.doc.toString()
		destroyEditor()
		createEditor(html)
	})

  onMount(() => {
		createEditor()
		return destroyEditor
	})

  onMount(() => {
    window.addEventListener(EVENT_LINE_TO_EDITOR, handleReportLineClickEvent)
    return () => window.removeEventListener(EVENT_LINE_TO_EDITOR, handleReportLineClickEvent)
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
		box-shadow: inset 0 -10rem 0 rgb(158 158 158 / 20%);
	}

	.editor-area-btn:active {
		background-color: var(--buttonBgActiveColor);
	}
</style>

<div class="editor-area">
  <div class="editor-area-edit" bind:this={editorElement}></div>
</div>
<div class="editor-footer">
  <button class="editor-area-btn" on:click|preventDefault={onSubmitHtml}>Check email HTML and CSS</button>
</div>
