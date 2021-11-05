# frozen_string_literal: true

require 'faraday'
require 'faraday_middleware'
require 'active_support/all'

class CaniuseGenerator

  HTML_TAGS_MAPS = {
    'html-address' => [['address', '']],
    'html-audio' => [['audio', '']],
    'html-bdi' => [['bdi', '']],
    'html-blockquote' => [['blockquote', '']],
    'html-button-reset' => [['button', 'type||reset']],
    'html-button-submit' => [['button', 'type||submit']],
    'html-del' => [['del', '']],
    'html-dialog' => [['dialog', '']],
    'html-dir' => [['dir', '']],
    'html-div' => [['div', '']],
    'html-form' => [['form', '']],
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
    'html-meter' => [['meter', '']],
    'html-object' => [['object', '']],
    'html-p' => [['p', '']],
    'html-picture' => [['picture', '']],
    'html-progress' => [['progress', '']],
    'html-rp' => [['rp', '']],
    'html-rt' => [['rt', '']],
    'html-ruby' => [['ruby', '']],
    'html-select' => [['select', '']],
    'html-span' => [['span', '']],
    'html-strike' => [['strike', '']],
    'html-strong' => [['strong', '']],
    'html-style' => [['style', '']],
    'html-svg' => [['svg', '']],
    'html-table' => [['table', '']],
    'html-textarea' => [['textarea', '']],
    'html-video' => [['video', '']],
    'html-wbr' => [['wbr', '']]
  }.freeze

  HTML_ATTRIBUTES_MAPS = {
    'html-aria-describedby' => [['aria-describedby', '']],
    'html-aria-hidden' => [['aria-hidden', '']],
    'html-aria-label' => [['aria-label', '']],
    'html-aria-labelledby' => [['aria-labelledby', '']],
    'html-aria-live' => [['aria-live', '']],
    'html-background' => [['background', '']],
    'html-height' => [['height', '']],
    'html-lang' => [['lang', '']],
    'html-loading-attribute' => [['loading', '']],
    'html-required' => [['required', '']],
    'html-role' => [['role', '']],
    'html-srcset' => [['srcset', ''], ['sizes', '']],
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
    'css-border-image' => [['border-image', '']],
    'css-border-radius' => [['border-radius', '']],
    'css-border' => [['border', '']],
    'css-box-shadow' => [['box-shadow', '']],
    'css-box-sizing' => [['box-sizing', '']],
    'css-caption-side' => [['caption-side', '']],
    'css-clip-path' => [['clip-path', '']],
    'css-column-count' => [['column-count', '']],
    'css-direction' => [['direction', '']],
    'css-display-flex' => [%w[display flex]],
    'css-display-grid' => [%w[display grid]],
    'css-display-none' => [%w[display none]],
    'css-filter' => [['filter', '']],
    'css-flex-direction' => [['flex-direction', '']],
    'css-flex-wrap' => [['flex-wrap', '']],
    'css-float' => [['float', '']],
    'css-font-weight' => [['font-weight', '']],
    'css-font' => [['font', '']],
    'css-height' => [['height', '']],
    'css-inline-size' => [['inline-size', '']],
    'css-justify-content' => [['justify-content', '']],
    'css-left-right-top-bottom' => [['left', ''], ['right', ''], ['top', ''], ['bottom', '']],
    'css-letter-spacing' => [['letter-spacing', '']],
    'css-line-height' => [['line-height', '']],
    'css-list-style-image' => [['list-style-image', '']],
    'css-list-style-position' => [['list-style-position', '']],
    'css-list-style-type' => [['list-style-type', '']],
    'css-list-style' => [['list-style', '']],
    'css-margin' => [['margin', ''], ['margin-top', ''], ['margin-bottom', ''], ['margin-left', ''],
                     ['margin-right', '']],
    'css-max-width' => [['max-width', '']],
    'css-mix-blend-mode' => [['mix-blend-mode', '']],
    'css-object-fit' => [['object-fit', '']],
    'css-object-position' => [['object-position', '']],
    'css-opacity' => [['opacity', '']],
    'css-overflow' => [['overflow', '']],
    'css-padding' => [['padding', ''], ['padding-top', ''], ['padding-bottom', ''], ['padding-left', ''],
                      ['padding-right', '']],
    'css-position' => [['position', '']],
    'css-text-align' => [['text-align', '']],
    'css-text-decoration-color' => [['text-decoration-color', '']],
    'css-text-decoration-thickness' => [['text-decoration-thickness', '']],
    'css-text-decoration' => [['text-decoration', '']],
    'css-text-indent' => [['text-indent', '']],
    'css-text-overflow' => [['text-overflow', '']],
    'css-text-shadow' => [['text-shadow', '']],
    'css-text-transform' => [['text-transform', '']],
    'css-text-underline-offset' => [['text-underline-offset', '']],
    'css-transform' => [['transform', '']],
    'css-vertical-align' => [['vertical-align', '']],
    'css-visibility' => [['visibility', '']],
    'css-white-space' => [['white-space', '']],
    'css-width' => [['width', '']],
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
    'css-rgb' => 'rgb',
    'css-rgba' => 'rgba',
    'css-unit-calc' => 'calc',
    'css-function-clamp' => 'clamp',
    'css-function-max' => 'max',
    'css-function-min' => 'min'
  }.freeze

  CSS_PSEUDO_SELECTORS_MAPS = {
    'css-pseudo-class-active' => 'active',
    'css-pseudo-class-checked' => 'checked',
    'css-pseudo-class-first-child' => 'first-child',
    'css-pseudo-class-first-of-type' => 'first-of-type',
    'css-pseudo-class-focus' => 'focus',
    'css-pseudo-class-hover' => 'hover',
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
    'css-pseudo-element-placeholder' => 'placeholder'
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
    'css-at-supports' => [['@fsupports', '']]
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
      css_variables: generate_css_variables
    }

    File.open(file, 'w') { |f| f.write JSON.dump(rules) }
  end

  private

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

  def generate_css_variables
    rule = data.detect { |r| r['slug'] == 'css-variables' }
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

      agg[v] = {
        notes: rule['notes_by_num'],
        stats: normalize_support(rule['stats']),
        url: rule['url']
      }
    end
  end

  def generate_multi_level_maps(maps)
    maps.each_with_object({}) do |(k, v), agg|
      rule = data.detect { |r| r['slug'] == k }
      next unless rule.present?

      v.each do |item|
        agg[item[0]] ||= {}
        agg[item[0]][item[1]] = {
          notes: rule['notes_by_num'] || [],
          stats: normalize_support(rule['stats']) ,
          url: rule['url'] || ''
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
          k => v.split.map { |rd| rd.delete('#') }
        )
      end
    end
  end

end

# Generate JSON doc from caniemail data - https://www.caniemail.com/api/data.json
namespace :caniemail do
  desc 'Generate JSON doc from caniemail data'
  task :generate do |_t, _args|
    CaniuseGenerator.new.generate(File.expand_path('../../wasm/parser/caniuse.json', __dir__))
    puts 'Work done'
  end
end
