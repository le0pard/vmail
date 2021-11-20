# frozen_string_literal: true

require 'erb'
require 'active_support'
require 'active_support/core_ext'
require 'rails-html-sanitizer'

module AssetsHelpers

  include ActionView::Helpers::SanitizeHelper

  def assets_manifest
    public_manifest_path = File.expand_path File.join(
      File.dirname(__FILE__),
      '../.tmp/dist/assets-manifest.json'
    )
    if File.exist?(public_manifest_path)
      JSON.parse(File.read(public_manifest_path))
    else
      {}
    end
  end

  def javascript_pack_tag(name)
    file_name = "#{name}.js"
    %(
      <script
        src="#{asset_pack_path(file_name)}"
        defer="defer"
        async="async"
        data-turbo-track="reload"></script>
    )
  end

  def stylesheet_pack_tag(name)
    file_name = "#{name}.css"
    %(
      <link
        href="#{asset_pack_path(file_name)}"
        rel="stylesheet"
        media="all" />
    )
  end

  def asset_pack_path(name)
    assets_manifest.dig(name.to_s, 'src') || raise("asset #{name} not found in #{assets_manifest.inspect}")
  end

  def asset_pack_integrity(name)
    assets_manifest.dig(name.to_s,
                        'integrity') || raise("integrity for asset #{name} not found in #{assets_manifest.inspect}")
  end

  def current_link_class(path = '/')
    current_page.path == path ? 'active' : ''
  end

  def svg_sprite_icons
    svg_html_safe File.new(File.expand_path('../source/svg/sprite.svg', __dir__)).read
  end

  def svg_sprite(name, options = {})
    options[:class] = [
      "svg-icon svg-icon--#{name}",
      options[:size] ? "svg-icon--#{options[:size]}" : nil,
      options[:class]
    ].compact.join(' ')

    icon = "<svg class='svg-icon__cnt'><use href='##{name}-svg-icon'/></svg>"

    svg_html_safe "
      <div class='#{options[:class]}'>
        #{wrap_svg_spinner icon, options[:class]}
      </div>
    "
  end

  def wrap_svg_spinner(html, klass)
    if klass.include?('spinner')
      svg_html_safe "<div class='icon__spinner'>#{html}</div>"
    else
      html
    end
  end

  def svg_html_safe(html)
    html.respond_to?(:html_safe) ? html.html_safe : html
  end

  def sanitize_tags(html)
    Rails::Html::FullSanitizer.new.sanitize(html)
  end

end
