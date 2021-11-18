# frozen_string_literal: true

Encoding.default_external = Encoding::UTF_8
Encoding.default_internal = Encoding::UTF_8

require 'lib/middleman_patches'

::Middleman::Extensions.register(:auto_blank_links) do
  require 'lib/auto_blank_links_extension'
  ::AutoBlankLinksExtension
end

###
# Blog settings
###

Time.zone = 'Kyiv'

# Static pages
proxy '/faq.html', '/static_pages/faq.html', ignore: true

###
# Helpers
###

require 'lib/defaults_site_helpers'
helpers DefaultSiteHelpers
require 'lib/assets_helpers'
helpers AssetsHelpers

set :images_dir, 'images'
set :markdown_engine, :kramdown
set :markdown, filter_html: false, fenced_code_blocks: true, smartypants: true
set :encoding, 'utf-8'

assets_dir = File.expand_path('.tmp/dist', __dir__)

activate :external_pipeline,
         name: :webpack,
         command: "yarn run assets:#{build? ? 'build' : 'watch'}",
         source: assets_dir,
         latency: 2,
         ignore_exit_code: true

activate :auto_blank_links,
         ignore_hostnames: ['vmail.leopard.in.ua', 'www.vmail.leopard.in.ua']

# Build-specific configuration
configure :build do
  config[:site_urls_base] = 'https://vmail.leopard.in.ua'
  # min html
  activate :minify_html
  # gzip
  activate :gzip, exts: %w[.css .htm .html .js .svg .xhtml]
end

# after_build do
#   system('yarn run optimize')
# end
