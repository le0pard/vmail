# frozen_string_literal: true

require 'faraday'
require 'faraday_middleware'
require 'active_support/all'

class CaniuseGenerator # rubocop:disable Metrics/ClassLength

  SINGLE_KEY_MAP = %w[
    css-variables
    css-important
    html-doctype
    css-nesting
  ].freeze

  HTML_TAGS_MAPS = {
    'amp' => [['html', '⚡4email'], %w[html amp4email]],
    'html-abbr' => [['abbr', '']],
    'html-acronym' => [['acronym', '']],
    'html-address' => [['address', '']],
    'html-audio' => [['audio', '']],
    'html-base' => [['base', '']],
    'html-bdi' => [['bdi', '']],
    'html-body' => [['body', '']],
    'html-blockquote' => [['blockquote', '']],
    'html-button-reset' => [['button', 'type||reset']],
    'html-button-submit' => [['button', 'type||submit']],
    'html-code' => [['code', '']],
    'html-dfn' => [['dfn', '']],
    'html-del' => [['del', '']],
    'html-dialog' => [['dialog', '']],
    'html-dir' => [['dir', '']],
    'html-div' => [['div', '']],
    'html-img' => [['img', '']],
    'html-form' => [['form', '']],
    'html-hr' => [['hr', '']],
    'html-h1-h6' => [['h1', ''], ['h2', ''], ['h3', ''], ['h4', ''], ['h5', ''], ['h6', '']],
    'html-image-maps' => [%w[img usemap]],
    'html-input-checkbox' => [['input', 'type||checkbox']],
    'html-input-hidden' => [['input', 'type||hidden']],
    'html-input-radio' => [['input', 'type||radio']],
    'html-input-reset' => [['input', 'type||reset']],
    'html-input-submit' => [['input', 'type||submit']],
    'html-input-text' => [['input', 'type||text']],
    'html-link' => [['link', '']],
    'html-lists' => [['ul', ''], ['ol', ''], ['dl', '']],
    'html-marquee' => [['marquee', '']],
    'html-meta-color-scheme' => [['meta', 'name||theme-color']],
    'html-meter' => [['meter', '']],
    'html-object' => [['object', '']],
    'html-p' => [['p', '']],
    'html-picture' => [['picture', '']],
    'html-pre' => [['pre', '']],
    'html-progress' => [['progress', '']],
    'html-rp' => [['rp', '']],
    'html-rt' => [['rt', '']],
    'html-ruby' => [['ruby', '']],
    'html-select' => [['select', '']],
    'html-semantics' => [
      ['article', ''], ['aside', ''], ['details', ''],
      ['figcaption', ''], ['figure', ''], ['footer', ''],
      ['header', ''], ['main', ''], ['mark', ''],
      ['nav', ''], ['section', ''], ['summary', ''],
      ['time', '']
    ],
    'html-small' => [['small', '']],
    'html-span' => [['span', '']],
    'html-strike' => [['strike', '']],
    'html-strong' => [['strong', '']],
    'html-srcset' => [
      %w[img srcset], %w[img sizes],
      %w[source srcset], %w[source sizes]
    ],
    'html-style' => [['style', '']],
    'html-svg' => [['svg', '']],
    'html-table' => [['table', '']],
    'html-target' => [%w[a target]],
    'html-textarea' => [['textarea', '']],
    'html-video' => [['video', '']],
    'html-wbr' => [['wbr', '']]
  }.freeze

  HTML_ATTRIBUTES_MAPS = {
    'html-align' => [['align', '']],
    'html-aria-describedby' => [['aria-describedby', '']],
    'html-aria-hidden' => [['aria-hidden', '']],
    'html-aria-label' => [['aria-label', '']],
    'html-aria-labelledby' => [['aria-labelledby', '']],
    'html-aria-live' => [['aria-live', '']],
    'html-background' => [['background', '']],
    'html-height' => [['height', '']],
    'html-hidden' => [['hidden', '']],
    'html-lang' => [['lang', '']],
    'html-loading-attribute' => [['loading', '']],
    'html-popover' => [['popover', ''], ['popovertarget', '']],
    'html-required' => [['required', '']],
    'html-role' => [['role', '']],
    'html-valign' => [['valign', '']],
    'html-width' => [['width', '']]
  }.freeze

  CSS_SELECTOR_TYPES_MAPS = {
    'css-selector-adjacent-sibling' => '0',
    'css-selector-attribute' => '1',
    'css-selector-chaining' => '2',
    'css-selector-child' => '3',
    'css-selector-class' => '4',
    'css-selector-descendant' => '5',
    'css-selector-general-sibling' => '6',
    'css-selector-grouping' => '7',
    'css-selector-id' => '8',
    'css-selector-type' => '9',
    'css-selector-universal' => '10'
  }.freeze

  CSS_PROPERTIES_MAPS = {
    'css-accent-color' => [['accent-color', '']],
    'css-align-items' => [['align-items', '']],
    'css-animation' => [['animation', '']],
    'css-aspect-ratio' => [['aspect-ratio', '']],
    'css-background-blend-mode' => [['background-blend-mode', '']],
    'css-background-clip' => [['background-clip', '']],
    'css-background-color' => [['background-color', '']],
    'css-background-image' => [['background-image', '']],
    'css-background-origin' => [['background-origin', '']],
    'css-background-position' => [['background-position', '']],
    'css-background-repeat' => [['background-repeat', '']],
    'css-background-size' => [['background-size', '']],
    'css-background' => [['background', '']],
    'css-backdrop-filter' => [['backdrop-filter', '']],
    'css-conic-gradient' => [%w[background conic-gradient]],
    'css-border-collapse' => [['border-collapse', '']],
    'css-border-radius-logical' => [['border-start-start-radius', ''], ['border-start-end-radius', ''],
                                    ['border-end-start-radius', ''], ['border-end-end-radius', '']],
    'css-border-image' => [['border-image', '']],
    'css-border-radius' => [['border-radius', '']],
    'css-border-spacing' => [['border-spacing', '']],
    'css-border' => [['border', '']],
    'css-box-shadow' => [['box-shadow', '']],
    'css-box-sizing' => [['box-sizing', '']],
    'css-block-inline-size' => [['block-size', ''], ['inline-size', '']],
    'css-border-inline-block-individual' => [['border-inline', ''], ['border-block', '']],
    'css-border-inline-block-longhand' => [['border-inline-start', ''], ['border-inline-end', ''],
                                           ['border-block-start', ''], ['border-block-end', '']],
    'css-border-inline-block' => [['border-inline', ''], ['border-block', '']],
    'css-caption-side' => [['caption-side', '']],
    'css-clip-path' => [['clip-path', '']],
    'css-color-scheme' => [['color-scheme', '']],
    'css-column-layout-properties' => [['columns', ''], ['column-fill', ''], ['column-rule', ''], ['column-gap', ''],
                                       ['column-span', '']],
    'css-column-count' => [['column-count', '']],
    'css-direction' => [['direction', '']],
    'css-display' => [['display', '']],
    'css-display-flex' => [%w[display flex]],
    'css-display-grid' => [%w[display grid]],
    'css-display-none' => [%w[display none]],
    'css-empty-cells' => [['empty-cells', '']],
    'css-filter' => [['filter', '']],
    'css-flex-direction' => [['flex-direction', '']],
    'css-flex-wrap' => [['flex-wrap', '']],
    'css-float' => [['float', '']],
    'css-font-size' => [['font-size', '']],
    'css-font-stretch' => [['font-stretch', '']],
    'css-font-weight' => [['font-weight', '']],
    'css-font' => [['font', '']],
    'css-sytem-ui' => [
      %w[font-family system-ui], %w[font-family ui-serif], %w[font-family ui-sans-serif],
      %w[font-family ui-monospace], %w[font-family ui-rounded]
    ],
    'css-font-kerning' => [['font-kerning', '']],
    'css-gap' => [['gap', '']],
    'css-grid-template' => [['grid-template', ''], ['grid-template-areas', ''], ['grid-template-columns', ''],
                            ['grid-template-rows', '']],
    'css-height' => [['height', '']],
    'css-hyphens' => [['hyphens', '']],
    'css-hyphenate-character' => [['hyphenate-character', '']],
    'css-hyphenate-limit-chars' => [['hyphenate-limit-chars', '']],
    'css-inline-size' => [['inline-size', '']],
    'css-inset' => [['inset', '']],
    'css-justify-content' => [['justify-content', '']],
    'css-left-right-top-bottom' => [['left', ''], ['right', ''], ['top', ''], ['bottom', '']],
    'css-letter-spacing' => [['letter-spacing', '']],
    'css-line-height' => [['line-height', '']],
    'css-list-style-image' => [['list-style-image', '']],
    'css-list-style-position' => [['list-style-position', '']],
    'css-list-style-type' => [['list-style-type', '']],
    'css-intrinsic-size' => [
      %w[width fit-content], %w[height fit-content], %w[min-width fit-content],
      %w[min-height fit-content], %w[max-width fit-content], %w[max-height fit-content],
      %w[inline-size fit-content], %w[block-size fit-content],
      %w[width min-content], %w[height min-content], %w[min-width min-content],
      %w[min-height min-content], %w[max-width min-content], %w[max-height min-content],
      %w[inline-size min-content], %w[block-size min-content],
      %w[width max-content], %w[height max-content], %w[min-width max-content],
      %w[min-height max-content], %w[max-width max-content], %w[max-height max-content],
      %w[inline-size max-content], %w[block-size max-content]
    ],
    'css-list-style' => [['list-style', '']],
    'css-margin' => [['margin', '']],
    'css-margin-block-start-end' => [['margin-block-start', ''], ['margin-block-end', '']],
    'css-margin-inline' => [['margin-block', ''], ['margin-inline', '']],
    'css-margin-inline-block' => [['margin-block', ''], ['margin-inline', '']],
    'css-margin-inline-start-end' => [['margin-inline-start', ''], ['margin-inline-end', '']],
    'css-max-inline-size' => [['max-inline-size', '']],
    'css-max-height' => [['max-height', '']],
    'css-max-width' => [['max-width', '']],
    'css-max-block-size' => [['max-block-size', '']],
    'css-min-block-size' => [['min-block-size', '']],
    'css-min-height' => [['min-height', '']],
    'css-min-width' => [['min-width', '']],
    'css-min-inline-size' => [['min-inline-size', '']],
    'css-mix-blend-mode' => [['mix-blend-mode', '']],
    'css-object-fit' => [['object-fit', '']],
    'css-object-position' => [['object-position', '']],
    'css-orphans' => [['orphans', '']],
    'css-opacity' => [['opacity', '']],
    'css-outline' => [['outline', '']],
    'css-outline-offset' => [['outline-offset', '']],
    'css-overflow' => [['overflow', '']],
    'css-overflow-wrap' => [['overflow-wrap', '']],
    'css-padding' => [['padding', '']],
    'css-padding-block-start-end' => [['padding-block-start', ''], ['padding-block-end', '']],
    'css-padding-inline-block' => [['padding-block', ''], ['padding-inline', '']],
    'css-padding-inline-start-end' => [['padding-inline-start', ''], ['padding-inline-end', '']],
    'css-position' => [['position', '']],
    'css-resize' => [['resize', '']],
    'css-shape-margin' => [['shape-margin', '']],
    'css-shape-outside' => [['shape-outside', '']],
    'css-scroll-snap' => [['scroll-snap-type', '']],
    'css-tab-size' => [['tab-size', '']],
    'css-table-layout' => [['table-layout', '']],
    'css-text-align' => [['text-align', '']],
    'css-text-align-last' => [['text-align-last', '']],
    'css-text-decoration-color' => [['text-decoration-color', '']],
    'css-text-decoration-thickness' => [['text-decoration-thickness', '']],
    'css-text-decoration' => [['text-decoration', '']],
    'css-text-decoration-line' => [['text-decoration-line', '']],
    'css-text-decoration-skip-ink' => [['text-decoration-skip-ink', '']],
    'css-text-decoration-style' => [['text-decoration-style', '']],
    'css-text-emphasis-position' => [['text-emphasis-position', '']],
    'css-text-emphasis' => [['text-emphasis', '']],
    'css-text-indent' => [['text-indent', '']],
    'css-text-justify' => [['text-justify', '']],
    'css-text-orientation' => [['text-orientation', '']],
    'css-text-overflow' => [['text-overflow', '']],
    'css-text-shadow' => [['text-shadow', '']],
    'css-text-transform' => [['text-transform', '']],
    'css-text-underline-offset' => [['text-underline-offset', '']],
    'css-text-underline-position' => [['text-underline-position', '']],
    'css-text-wrap' => [['text-wrap', '']],
    'css-transition' => [['transition', '']],
    'css-transform' => [['transform', '']],
    'css-user-select' => [['user-select', '']],
    'css-vertical-align' => [['vertical-align', '']],
    'css-visibility' => [['visibility', '']],
    'css-white-space' => [['white-space', '']],
    'css-white-space-collapse' => [['white-space-collapse', '']],
    'css-word-break' => [['word-break', '']],
    'css-word-spacing' => [['word-spacing', '']],
    'css-widows' => [['word-widows', '']],
    'css-width' => [['width', '']],
    'css-writing-mode' => [['writing-mode', '']],
    'css-z-index' => [['z-index', '']]
  }.freeze

  CSS_DIMENTIONS_MAPS = {
    'css-unit-ch' => 'ch',
    'css-unit-cm' => 'cm',
    'css-unit-em' => 'em',
    'css-unit-ex' => 'ex',
    'css-unit-in' => 'in',
    'css-unit-initial' => 'initial',
    'css-unit-mm' => 'mm',
    'css-unit-pc' => 'pc',
    'css-unit-percent' => '%',
    'css-unit-pt' => 'pt',
    'css-unit-px' => 'px',
    'css-unit-rem' => 'rem',
    'css-unit-vh' => 'vh',
    'css-unit-vmax' => 'vmax',
    'css-unit-vmin' => 'vmin',
    'css-unit-vw' => 'vw'
  }.freeze

  CSS_FUNCTIONS_MAPS = {
    'css-modern-color' => %w[lch oklch lab oklab],
    'css-linear-gradient' => ['linear-gradient'],
    'css-radial-gradient' => ['radial-gradient'],
    'css-rgb' => ['rgb'],
    'css-rgba' => ['rgba'],
    'css-unit-calc' => ['calc'],
    'css-function-clamp' => ['clamp'],
    'css-function-light-dark' => ['light-dark'],
    'css-function-max' => ['max'],
    'css-function-min' => ['min']
  }.freeze

  CSS_PSEUDO_SELECTORS_MAPS = {
    'css-pseudo-class-active' => 'active',
    'css-pseudo-class-checked' => 'checked',
    'css-pseudo-class-first-child' => 'first-child',
    'css-pseudo-class-first-of-type' => 'first-of-type',
    'css-pseudo-class-focus' => 'focus',
    'css-pseudo-class-has' => 'has',
    'css-pseudo-class-hover' => 'hover',
    'css-pseudo-class-lang' => 'lang',
    'css-pseudo-class-last-child' => 'last-child',
    'css-pseudo-class-last-of-type' => 'last-of-type',
    'css-pseudo-class-link' => 'link',
    'css-pseudo-class-not' => 'not',
    'css-pseudo-class-nth-child' => 'nth-child',
    'css-pseudo-class-nth-last-child' => 'nth-last-child',
    'css-pseudo-class-nth-last-of-type' => 'nth-last-of-type',
    'css-pseudo-class-nth-of-type' => 'nth-of-type',
    'css-pseudo-class-only-child' => 'only-child',
    'css-pseudo-class-only-of-type' => 'only-of-type',
    'css-pseudo-class-target' => 'target',
    'css-pseudo-class-visited' => 'visited',
    'css-pseudo-element-after' => 'after',
    'css-pseudo-element-before' => 'before',
    'css-pseudo-element-first-letter' => 'first-letter',
    'css-pseudo-element-first-line' => 'first-line',
    'css-pseudo-element-placeholder' => 'placeholder',
    'css-pseudo-element-marker' => 'marker'
  }.freeze

  AT_RULE_CSS_STATEMENTS_MAPS = {
    'css-at-font-face' => [['@font-face', '']],
    'css-at-import' => [['@import', '']],
    'css-at-keyframes' => [['@keyframes', '']],
    'css-at-media-device-pixel-ratio' => [['@media', 'device-pixel-ratio']],
    'css-at-media-orientation' => [['@media', 'orientation']],
    'css-at-media-prefers-color-scheme' => [['@media', 'prefers-color-scheme']],
    'css-at-media-prefers-reduced-motion' => [['@media', 'prefers-reduced-motion']],
    'css-at-media' => [['@media', '']],
    'css-at-media-hover' => [['@media', 'hover'], ['@media', 'any-hover']],
    'css-at-supports' => [['@supports', '']]
  }.freeze

  IMG_FORMATS_MAPS = {
    'image-apng' => 'apng',
    'image-avif' => 'avif',
    'image-base64' => 'base64',
    'image-bmp' => 'bmp',
    'image-gif' => 'gif',
    'image-hdr' => 'hdr',
    'image-heif' => 'heif',
    'image-ico' => 'ico',
    'image-jpg' => 'jpg',
    'image-mp4' => 'mp4',
    'image-png' => 'png',
    'image-ppm' => 'ppm',
    'image-svg' => 'svg',
    'image-tiff' => 'tiff',
    'image-webp' => 'webp'
  }.freeze

  LINK_TYPES_MAP = {
    'html-anchor-links' => 'anchor',
    'html-mailto-links' => 'mailto'
  }.freeze

  attr_reader :data

  def initialize
    conn = Faraday.new do |f|
      f.request :json # encode req bodies as JSON
      f.request :retry, {
        max: 2,
        interval: 0.05,
        interval_randomness: 0.5,
        backoff_factor: 2
      } # retry transient failures
      f.response :follow_redirects # follow redirects
      f.response :json # decode response bodies as JSON
      f.response :raise_error # raise error on bad HTTP code
    end
    @data = conn.get('https://www.caniemail.com/api/data.json').body['data']
  end

  # Metrics/AbcSize
  def generate(file)
    rules = {
      html_tags: generate_html_tags,
      html_attributes: generate_html_attributes,
      css_properties: generate_css_properties,
      css_selector_types: generate_css_selector_types,
      css_dimentions: generate_css_dimentions,
      css_functions: generate_css_functions,
      css_pseudo_selectors: generate_css_pseudo_selectors,
      at_rule_css_statements: generate_at_rule_css_statements,
      img_formats: generate_img_formats,
      link_types: generate_link_types,
      css_variables: generate_for_single_key('css-variables'),
      css_important: generate_for_single_key('css-important'),
      html5_doctype: generate_for_single_key('html-doctype')
    }

    File.write(file, JSON.dump(rules))
    warn_about_now_covered_rules
  end

  private

  def warn_about_now_covered_rules # rubocop:disable Metrics/MethodLength, Metrics/AbcSize, Metrics/CyclomaticComplexity, Metrics/PerceivedComplexity
    skipped_rules = %w[
      bimi
      css-comments
      html-cellpadding
      html-cellspacing
      html-comments
    ]

    rules_without_apply = data.filter do |r|
      !SINGLE_KEY_MAP.include?(r['slug']) &&
        !HTML_TAGS_MAPS.key?(r['slug']) &&
        !HTML_ATTRIBUTES_MAPS.key?(r['slug']) &&
        !CSS_PROPERTIES_MAPS.key?(r['slug']) &&
        !CSS_SELECTOR_TYPES_MAPS.key?(r['slug']) &&
        !CSS_DIMENTIONS_MAPS.key?(r['slug']) &&
        !CSS_FUNCTIONS_MAPS.key?(r['slug']) &&
        !CSS_PSEUDO_SELECTORS_MAPS.key?(r['slug']) &&
        !AT_RULE_CSS_STATEMENTS_MAPS.key?(r['slug']) &&
        !IMG_FORMATS_MAPS.key?(r['slug']) &&
        !LINK_TYPES_MAP.key?(r['slug']) &&
        !skipped_rules.include?(r['slug']) # skip some rules
    end

    return if rules_without_apply.empty?

    $stdout.puts "WARN, This rules was skipped: #{rules_without_apply.map { |r| r['slug'] }.join(', ')}"
  end

  def generate_html_tags
    generate_multi_level_maps(HTML_TAGS_MAPS)
  end

  def generate_html_attributes
    generate_multi_level_maps(HTML_ATTRIBUTES_MAPS)
  end

  def generate_css_properties
    generate_multi_level_maps(CSS_PROPERTIES_MAPS)
  end

  def generate_css_selector_types
    generate_one_level_maps(CSS_SELECTOR_TYPES_MAPS)
  end

  def generate_css_dimentions
    generate_one_level_maps(CSS_DIMENTIONS_MAPS)
  end

  def generate_css_functions
    generate_one_level_maps(CSS_FUNCTIONS_MAPS)
  end

  def generate_css_pseudo_selectors
    generate_one_level_maps(CSS_PSEUDO_SELECTORS_MAPS)
  end

  def generate_at_rule_css_statements
    generate_multi_level_maps(AT_RULE_CSS_STATEMENTS_MAPS)
  end

  def generate_img_formats
    generate_one_level_maps(IMG_FORMATS_MAPS)
  end

  def generate_link_types
    generate_one_level_maps(LINK_TYPES_MAP)
  end

  def generate_for_single_key(key)
    rule = data.detect { |r| r['slug'] == key }
    if rule.present?
      {
        notes: rule['notes_by_num'],
        stats: normalize_support(rule['stats']),
        url: rule['url']
      }
    else
      {}
    end
  end

  def generate_one_level_maps(maps)
    maps.each_with_object({}) do |(k, v), agg|
      rule = data.detect { |r| r['slug'] == k }
      next unless rule.present?

      next if count_not_support(rule['stats']).zero? # supported in all clients

      if v.is_a?(Array)
        v.each do |vk|
          agg[vk] = {
            notes: rule['notes_by_num'] || [],
            stats: normalize_support(rule['stats']),
            url: rule['url'] || '',
            description: rule['description'] || ''
          }
        end
      else
        agg[v] = {
          notes: rule['notes_by_num'] || [],
          stats: normalize_support(rule['stats']),
          url: rule['url'] || '',
          description: rule['description'] || ''
        }
      end

    end
  end

  def generate_multi_level_maps(maps)
    maps.each_with_object({}) do |(k, v), agg|
      rule = data.detect { |r| r['slug'] == k }
      next unless rule.present?

      next if count_not_support(rule['stats']).zero? # supported in all clients

      v.each do |item|
        agg[item[0]] ||= {}
        agg[item[0]][item[1]] = {
          notes: rule['notes_by_num'] || [],
          stats: normalize_support(rule['stats']),
          url: rule['url'] || '',
          description: rule['description'] || ''
        }
      end
    end
  end

  def normalize_support(hash)
    hash.reduce({}) do |agg, (k, v)|
      case v
      when Hash
        agg.merge(
          k => normalize_support(v)
        )
      when String
        agg.merge(
          k => v.split.map { |rd| rd.delete('#') }.map(&:downcase)
        )
      end
    end
  end

  def count_not_support(hash)
    hash.reduce(0) do |agg, (_k, v)|
      case v
      when Hash
        agg + count_not_support(v)
      when String
        agg + (v.casecmp('y').zero? ? 0 : 1)
      end
    end
  end

end

# Generate JSON doc from caniemail data - https://www.caniemail.com/api/data.json
namespace :caniemail do
  desc 'Generate JSON doc from caniemail data'
  task :generate do |_t, _args|
    CaniuseGenerator.new.generate(File.expand_path('../../wasm_parser/parser/caniuse.json', __dir__))
    $stdout.puts 'Work done'
  end
end
