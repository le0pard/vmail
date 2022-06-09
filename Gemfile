# frozen_string_literal: true

source 'https://rubygems.org'

git_source(:github) { |repo_name| "https://github.com/#{repo_name}" }

gem 'middleman', github: 'middleman/middleman'
# EXTENSIONS
gem 'middleman-minify-html', github: 'middleman/middleman-minify-html' # min html
# UTILS
gem 'actionpack', '>= 5.2.4.2', require: false
gem 'actionview', '>= 5.2.4.2', require: false
gem 'activesupport', '>= 5.2.4.2', require: false
gem 'addressable', '>= 2.7.0'
gem 'builder', '>= 3.2.2' # XML builder
gem 'erubis', '>= 2.7'
gem 'kramdown', '>= 2.3.0' # faster markdown
gem 'multi_json', '>= 1.15.0'
gem 'nokogiri', '>= 1.11.0.rc4'
gem 'oj', '>= 2.10.4' # faster JSON
gem 'rack'
gem 'rails-html-sanitizer', '>= 1.0.1', require: false
gem 'rake'

# Dev
group :development do
  gem 'faraday', '~> 1'
  gem 'faraday_middleware'
  gem 'lefthook', require: false
  gem 'rubocop', '>= 1.7.0', require: false
  gem 'rubocop-performance', '>= 1.9.1', require: false
  gem 'rubocop-rake', '>= 0.5.1', require: false
end
