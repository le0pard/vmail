<script>
  import {onMount} from 'svelte'
  import {report} from 'stores/report'
  import {EVENT_LINE_TO_EDITOR, EVENT_LINE_TO_REPORT} from 'lib/constants'

  const handleLineClick = (line) => {
    window.dispatchEvent(new window.CustomEvent(EVENT_LINE_TO_EDITOR, { detail: {line} }))
  }

  const handleEditorLineClickEvent = (e) => {
    if (!e.detail?.line) {
      return
    }

    const {line} = e.detail
    console.log('Scroll to line: ', line)
  }

  onMount(() => {
    window.addEventListener(EVENT_LINE_TO_REPORT, handleEditorLineClickEvent)
    return () => window.removeEventListener(EVENT_LINE_TO_REPORT, handleEditorLineClickEvent)
  })
</script>

<ul>
  {#if $report.html_tags}
    {#each Object.keys($report.html_tags).sort() as tagName (tagName)}
      {#each Object.keys($report.html_tags[tagName]).sort() as tagAttr (tagAttr)}
        <li>
          HTML Tag: {tagName}: {tagAttr}
          <div>
            {#each $report.html_tags[tagName][tagAttr].lines as line}
              <button on:click|preventDefault={() => handleLineClick(line)}>{line}</button>
            {/each}
          </div>
        </li>
      {/each}
    {/each}
  {/if}
  {#if $report.html_attributes}
    {#each Object.keys($report.html_attributes).sort() as attrName (attrName)}
      {#each Object.keys($report.html_attributes[attrName]).sort() as attrVal (attrVal)}
        <li>
          HTML Attribute: {attrName}: {attrVal}
          <div>
            {#each $report.html_attributes[attrName][attrVal].lines as line}
              <button on:click|preventDefault={() => handleLineClick(line)}>{line}</button>
            {/each}
          </div>
        </li>
      {/each}
    {/each}
  {/if}
  {#if $report.css_properties}
    {#each Object.keys($report.css_properties).sort() as propName (propName)}
      {#each Object.keys($report.css_properties[propName]).sort() as propVal (propVal)}
        <li>
          CSS Prop: {propName}: {propVal}
          <div>
            {#each $report.css_properties[propName][propVal].lines as line}
              <button on:click|preventDefault={() => handleLineClick(line)}>{line}</button>
            {/each}
          </div>
        </li>
      {/each}
    {/each}
  {/if}
  {#if $report.at_rule_css_statements}
    {#each Object.keys($report.at_rule_css_statements).sort() as atName (atName)}
      {#each Object.keys($report.at_rule_css_statements[atName]).sort() as atVal (atVal)}
        <li>
          AT Rule CSS Statement: {atName}: {atVal}
          <div>
            {#each $report.at_rule_css_statements[atName][atVal].lines as line}
              <button on:click|preventDefault={() => handleLineClick(line)}>{line}</button>
            {/each}
          </div>
        </li>
      {/each}
    {/each}
  {/if}
  {#if $report.css_selector_types}
    {#each Object.keys($report.css_selector_types).sort() as selectorType (selectorType)}
      <li>
        CSS SELECTOR: {selectorType}
      </li>
    {/each}
  {/if}
  {#if $report.css_dimentions}
    {#each Object.keys($report.css_dimentions).sort() as dimensionType (dimensionType)}
      <li>
        CSS Dimention: {dimensionType}
      </li>
    {/each}
  {/if}
  {#if $report.css_functions}
    {#each Object.keys($report.css_functions).sort() as functionType (functionType)}
      <li>
        CSS Function: {functionType}
      </li>
    {/each}
  {/if}
  {#if $report.css_pseudo_selectors}
    {#each Object.keys($report.css_pseudo_selectors).sort() as pseudoSelectorType (pseudoSelectorType)}
      <li>
        CSS Pseudo Selector: {pseudoSelectorType}
      </li>
    {/each}
  {/if}
  {#if $report.img_formats}
    {#each Object.keys($report.img_formats).sort() as imgType (imgType)}
      <li>
        Img formats: {imgType}
      </li>
    {/each}
  {/if}
  {#if $report.css_variables && $report.css_variables.length > 0}
    <li>
      CSS Variables: {JSON.stringify($report.css_variables)}
    </li>
  {/if}
</ul>
