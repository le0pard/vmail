# frozen_string_literal: true

require 'shellwords'

def compile_wasm(wasm_dir, out_file)
  command_args = Shellwords.split('-ldflags="-s -w"')
  system(
    {
      'GOOS' => 'js',
      'GOARCH' => 'wasm'
    },
    'go', 'build', *command_args,
    '-o', out_file,
    chdir: wasm_dir,
    unsetenv_others: false,
    exception: true
  )
end

namespace :wasm do
  desc 'Generate wasm parser file'
  task :parser do |_t, _args|
    wasm_dir = File.expand_path('../../wasm_parser', __dir__)
    out_file = File.expand_path('../../public/parser.wasm', __dir__)
    compile_wasm(wasm_dir, out_file)
    $stdout.puts 'Finished'
  end

  desc 'Generate wasm inliner file'
  task :inliner do |_t, _args|
    wasm_dir = File.expand_path('../../wasm_inliner', __dir__)
    out_file = File.expand_path('../../public/inliner.wasm', __dir__)
    compile_wasm(wasm_dir, out_file)
    $stdout.puts 'Finished'
  end
end
