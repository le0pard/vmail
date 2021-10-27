# frozen_string_literal: true

require 'shellwords'

namespace :wasm do
  desc 'Generate wasm file'
  task :compile do |_t, _args|
    wasm_dir = File.expand_path('../../wasm', __dir__)
    out_file = File.expand_path('../../source/parser.wasm', __dir__)
    command_with_args = Shellwords.split('build -ldflags="-s -w" -o')

    system(
      {
        'GOOS' => 'js',
        'GOARCH' => 'wasm'
      },
      'go',
      *command_with_args,
      out_file,
      unsetenv_others: false,
      exception: true,
      chdir: wasm_dir
    )
    puts 'Finished'
  end
end
