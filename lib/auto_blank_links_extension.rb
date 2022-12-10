# frozen_string_literal: true

require 'nokogiri'
require 'addressable/uri'

class AutoBlankLinksExtension < Middleman::Extension

  option :ignore_hostnames, [], 'Website internal hostnames'
  option :ignore_pages, [], 'Patterns to avoid target blanks for pages'
  option :content_types, %w[text/html], 'Content types of resources that contain HTML'

  def initialize(app, _options_hash = ::Middleman::EMPTY_HASH, &block)
    super

    @ignore_pages = Array(options[:ignore_pages])
    @ignore_hostnames = Array(options[:ignore_hostnames])
  end

  def manipulate_resource_list_container!(resource_list)
    resource_list.by_binary(false).each do |r|
      type = r.content_type.try(:slice, /^[^;]*/)
      r.add_filter method(:process_links) if valid_content_type?(type) && !ignore?(r.destination_path)
    end
  end

  def valid_content_type?(content_type)
    options[:content_types].include?(content_type)
  end
  memoize :valid_content_type?

  def ignore?(path)
    @ignore_pages.any? { |ignore| ::Middleman::Util.path_match(ignore, path) }
  end
  memoize :ignore?

  def process_links(html)
    content = Nokogiri::HTML.parse(html)
    anchors = content.css('a[href]')
    anchors.each do |item|
      next unless processable_link?(item)

      add_target_blank_attribute(item)
      add_rel_attributes(item)
    end
    content.to_html
  end
  memoize :process_links

  private

  def external?(link)
    return false if link.nil?

    parsed_link = Addressable::URI.parse(link)
    @ignore_hostnames.all? { |host| parsed_link.host != host } if parsed_link&.absolute?
  end

  def mailto_link?(link)
    link.start_with?('mailto:')
  end

  def processable_link?(link)
    !mailto_link?(link['href']) && external?(link['href'])
  end

  def add_target_blank_attribute(link)
    link['target'] = '_blank'
  end

  def add_rel_attributes(link)
    link['rel'] = 'noopener noreferrer'
  end

end
