require "kemal"

macro render_view(filename)
  render "src/views/#{{{filename}}}.ecr", "src/views/layout.ecr"
end

get "/" do
  "Hello, world!"
  render_view "index"
end

Kemal.run(5000)