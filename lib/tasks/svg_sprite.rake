# frozen_string_literal: true

require 'nokogiri'
require 'erb'

class SvgGenerator

  def initialize(svg_path, templates_dir, templates_out_dir)
    @svg_path = svg_path
    @templates_dir = templates_dir
    @templates_out_dir = templates_out_dir
  end

  def files
    @files ||= Dir.entries(@svg_path).select { |f| File.extname(f) == '.svg' }
  end

  def read_svg(filename)
    file = File.join(@svg_path, filename)
    File.read(file)
  end

  def icons
    files.map do |name|
      file        = read_svg(name)
      doc         = Nokogiri::HTML::DocumentFragment.parse(file)

      doc.css('*').remove_attr('fill')

      svg         = doc.at_css('svg')
      viewbox     = svg['viewbox']
      g           = svg.search('g')
      container   = g.empty? ? svg : g

      shape       = container.children.map(&:to_s).join
      name        = File.basename(name, '.svg')

      { name: name, viewbox: viewbox, shape: shape }
    end
  end

  def optimize(code)
    code.gsub(/$\s+/, '')
  end

  def sprite(template)
    view    = File.read File.join(@templates_dir, "#{template}.erb")
    result  = ERB.new(view).result(binding)
    optimize(result)
  end

  def generate(template)
    path = File.join(@templates_out_dir, template)
    file = File.new(path, 'w')
    file.write sprite(template)
    file.close
  end

end

# generate svg sprite
namespace :svg_sprite do
  desc 'Create svg sprite'
  task :generate do |_t, _args|
    SvgGenerator.new(
      File.expand_path('../svg/svg_icons', __dir__),
      File.expand_path('../svg', __dir__),
      File.expand_path('../../source/svg', __dir__)
    ).generate('sprite.svg')
    puts 'Work done'
  end
end
