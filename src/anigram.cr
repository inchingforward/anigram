require "kemal"

macro render_view(filename)
  render "src/views/#{{{filename}}}.ecr", "src/views/layout.ecr"
end

get "/" do
  render_view "index"
end

get "/about" do
  render_view "not_yet_implemented"
end

get "/animations/new" do
  render_view "not_yet_implemented"
end

get "/animations/:uuid/edit" do
  render_view "edit_animation"
end

get "/api/animations/%s" do
  render_view "not_yet_implemented"
end

post "/api/animations" do
  render_view "not_yet_implemented"
end

Kemal.run(5000)