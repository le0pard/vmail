# frozen_string_literal: true

module DefaultSiteHelpers

  def default_title_helper
    'RWpod - подкаст про мир Ruby и Web технологии'
  end

  def default_keywords_helper
    'RWpod, Ruby, Web, подкаст, русский подкаст, скринкасты, программирование'
  end

  def default_description_helper
    'RWpod - подкаст про мир Ruby и Web технологии (для тех, кому нравится мыслить в Ruby стиле).'
  end

  def full_url(url)
    "#{config[:site_urls_base] || ''}#{url_for(url)}"
  end

end
