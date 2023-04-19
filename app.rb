# frozen_string_literal: true

$LOAD_PATH << File.join(__dir__, '..', 'lib')

require 'rubygems'
require 'bundler/setup'
require 'sinatra/base'
require 'sinatra/reloader' if ENV['RACK_ENV'] == 'development'

class Mbdo < Sinatra::Base
  configure :development do
    register Sinatra::Reloader
  end

  configure :production, :development do
    enable :logging
  end

  get '/' do
    erb :'index.html'
  end
end
