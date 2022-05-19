defmodule ChatWeb.Router do
  use ChatWeb, :router
  import Phoenix.LiveDashboard.Router

  pipeline :browser do
    plug :accepts, ["html"]
  end

  pipeline :api do
    plug :accepts, ["json"]
  end

  scope "/api", ChatWeb do
    pipe_through :api
  end

  if Mix.env() == :dev do
    scope "/" do
      pipe_through :browser
      live_dashboard "/dashboard", ecto_repos: [Chat.Repo], metrics: ChatWeb.Telemetry
    end
  end
end
