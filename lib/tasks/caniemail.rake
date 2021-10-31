# frozen_string_literal: true

require 'faraday'
require 'faraday_middleware'

retry_options = {
  max: 2,
  interval: 0.05,
  interval_randomness: 0.5,
  backoff_factor: 2
}

conn = Faraday.new do |f|
  f.request :json # encode req bodies as JSON
  f.request :retry, retry_options # retry transient failures
  f.response :follow_redirects # follow redirects
  f.response :json # decode response bodies as JSON
  f.response :raise_error # raise error on bad HTTP code
end

# type.key.value.family.platform
# type.key.value.lines

# Generate JSON doc from caniemail data - https://www.caniemail.com/api/data.json
namespace :caniemail do
  desc 'Generate JSON doc from caniemail data'
  task :generate do |_t, _args|
    response = conn.get('https://www.caniemail.com/api/data.json')
    data = response.body['data']
    nicenames = response.body['nicenames']
    puts 'Work done'
  end
end
