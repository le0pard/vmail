# frozen_string_literal: true

module DefaultSiteHelpers

  def default_title_helper
    'VMail - check HTML & CSS compatibility with email clients'
  end

  def default_keywords_helper
    'vmail, html, css, validate, email clients, check, ruby, sass, wasm, golang, svelte'
  end

  def default_description_helper
    'VMail - check HTML & CSS compatibility with email clients'
  end

  def full_url(url)
    "#{config[:site_urls_base] || ''}#{url_for(url)}"
  end

end
