<script>
	import {onMount, onDestroy} from 'svelte'
	import {EditorState, EditorSelection, StateField, StateEffect} from '@codemirror/state'
	import {EditorView, keymap} from '@codemirror/view'
	import {defaultKeymap} from '@codemirror/commands'
	import {history, historyKeymap} from '@codemirror/history'
	import {
		lineNumbers,
		highlightActiveLineGutter,
		GutterMarker,
		gutter
	} from '@codemirror/gutter'
	import {defaultHighlightStyle} from '@codemirror/highlight'
	import {RangeSet} from '@codemirror/rangeset'
	import {html} from '@codemirror/lang-html'
	import {report} from 'stores/report'
	import ReportListComponent from './ReportList'

	let editorElement

	let editorView = null

	const eTheme = EditorView.baseTheme({
		'&.cm-editor': {
			fontSize: '0.9rem',
			height: '100%'
		}
	})

	// const eDarkTheme = EditorView.baseTheme({
	// 	"&": {
	// 		color: "white",
	// 		backgroundColor: "#034",
	// 		fontSize: '0.9rem',
	// 		height: '100%'
	// 	},
	// 	".cm-content": {
	// 		caretColor: "#0e9"
	// 	},
	// 	"&.cm-focused .cm-cursor": {
	// 		borderLeftColor: "#0e9"
	// 	},
	// 	"&.cm-focused .cm-selectionBackground, ::selection": {
	// 		backgroundColor: "#074"
	// 	},
	// 	".cm-gutters": {
	// 		backgroundColor: "#045",
	// 		color: "#ddd",
	// 		border: "none"
	// 	}
	// }, {dark: true})

	const breakpointEffect = StateEffect.define({
		map: (val, mapping) => ({pos: mapping.mapPos(val.pos), on: val.on})
	})

	const breakpointState = StateField.define({
		create() { return RangeSet.empty },
		update(set, transaction) {
			set = set.map(transaction.changes)
			for (let e of transaction.effects) {
				if (e.is(breakpointEffect)) {
					if (e.value.on)
						set = set.update({add: [breakpointMarker.range(e.value.pos)]})
					else
						set = set.update({filter: from => from != e.value.pos})
				}
			}
			return set
		}
	})

	const toggleBreakpoint = (view, pos) => {
		let breakpoints = view.state.field(breakpointState)
		let hasBreakpoint = false
		breakpoints.between(pos, pos, () => {hasBreakpoint = true})
		view.dispatch({
			effects: breakpointEffect.of({pos, on: !hasBreakpoint})
		})
	}

	const breakpointMarker = new class extends GutterMarker {
		toDOM() {
			const marker = document.createElement('div')
			marker.className = 'lint-marker-error'
			marker.innerText = 'err'
			marker.style = 'width: 10px'
			return marker
		}
	}

	const initialEditorState = EditorState.create({
		doc: '',
		extensions: [
			breakpointState,
			gutter({
				class: 'cm-breakpoint-gutter',
				markers: v => v.state.field(breakpointState),
				initialSpacer: () => breakpointMarker,
				domEventHandlers: {
					mousedown(view, line) {
						toggleBreakpoint(view, line.from)
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

	const unsubscribeReport = report.subscribe((rData) => {
		console.log('rData', rData)
	})

	const handleScroll = (line) => {
		const editorLine = editorView.state.doc.line(line)
		const selection = EditorSelection.cursor(editorLine.from)
		editorView.dispatch({effects: EditorView.centerOn.of(selection), selection})
		editorView.focus()
	}

	const onSubmitHtml = async () => {
		const html = editorView.state.doc.toString()

		try {
			const reportData = await window.VMail(html)
			report.update(reportData)
		} catch (err) {
			console.log('error', err)
		}
	}

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

	onDestroy(unsubscribeReport)
	onDestroy(() => report.reset())
</script>

<style>
	.parser-view {
		display: flex;
		position: absolute;
		bottom: 0;
    left: 0;
    right: 0;
    top: 0;
	}

	.parser-editor {
		display: flex;
		flex-direction: column;
		height: 100%;
		flex: 1;
		position: relative;
	}

	.parser-editor-header {
		flex-shrink: 0;
		overflow: hidden;
    white-space: nowrap;
	}

	.parser-editor-area {
		flex-grow: 1;
		position: relative;
	}

	.parser-editor-area-edit {
		height: 100%;
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
		overflow: scroll;
	}

	.parser-editor-header {
		flex-shrink: 0;
		overflow: hidden;
    white-space: nowrap;
	}

	.parser-resize {
		display: flex;
		width: 30px;
	}

	.parser-report {
		height: 100%;
		flex: 1;
		position: relative;
		display: flex;
		flex-direction: column;
	}

	.parser-report-header {
		flex-shrink: 0;
		overflow: hidden;
    white-space: nowrap;
	}

	.parser-report-area {
		flex-grow: 1;
		position: relative;
		overflow: scroll;
	}
</style>

<div class="parser-view">
	<div class="parser-editor">
		<div class="parser-editor-header">
			<a href="/">Home</a>
			<a href="/about.html">About</a>
		</div>
		<div class="parser-editor-area">
			<div class="parser-editor-area-edit" bind:this={editorElement}></div>
		</div>
		<div class="parser-editor-footer">
			<button on:click|preventDefault={onSubmitHtml}>Check</button>
		</div>
	</div>
	<div class="parser-resize"></div>
	<div class="parser-report">
		<div class="parser-report-header">
			<p>Header</p>
		</div>
		<div class="parser-report-area">
			{#if Object.keys($report).length > 0}
				<ReportListComponent handleScroll={handleScroll} />
			{:else}
				<p>Submit your HTML</p>
			{/if}
		</div>
	</div>
</div>
