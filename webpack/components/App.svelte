<script>
	import {onMount} from 'svelte'
	import {EditorState} from '@codemirror/state'
	import {EditorView, keymap} from '@codemirror/view'
	import {defaultKeymap} from '@codemirror/commands'
	import {history, historyKeymap} from '@codemirror/history'
	import {lineNumbers, highlightActiveLineGutter} from '@codemirror/gutter'
	import {defaultHighlightStyle} from '@codemirror/highlight'
	import {html} from '@codemirror/lang-html'

	let editorElement

	let editorView = null

	const eTheme = EditorView.baseTheme({
		'&.cm-editor': {
			fontSize: '0.9rem',
			height: '100%'
		}
	})

	const eState = EditorState.create({
		doc: '',
		extensions: [
			lineNumbers(),
			highlightActiveLineGutter(),
			history(),
			defaultHighlightStyle.fallback,
			keymap.of([
				...defaultKeymap,
				...historyKeymap
			]),
			html(),
			eTheme
		]
	})

	onMount(() => {
		editorView = new EditorView({
			state: eState,
			parent: editorElement
		})

		return () => {
			if (editorView) {
				editorView.destroy()
				editorView = null
			}
		}
	})

	const onSubmitHtml = async () => {
		const html = editorView.state.doc.toString()
		try {
			const report = await window.VMail(html)
			console.log('Report', report)
		} catch (err) {
			console.log('error', err)
		}
	}
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
			<p>Header</p>
		</div>
		<div class="parser-report-area">
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
			<p>report</p>
		</div>
	</div>
</div>
