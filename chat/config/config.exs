# This file is responsible for configuring your application
# and its dependencies with the aid of the Config module.
#
# This configuration file is loaded before any dependency and
# is restricted to this project.

# General application configuration
import Config

config :chat,
  ecto_repos: [Chat.Repo]

# Configures the endpoint
config :chat, ChatWeb.Endpoint,
  url: [host: "localhost"],
  port: System.get_env("PORT"),
  pubsub_server: Chat.PubSub,
  live_view: [signing_salt: "VqpB9cdA"]

# Configure your database
config :chat, Chat.Repo,
  username: System.get_env("POSTGRES_USER"),
  password: System.get_env("POSTGRES_PASSWORD"),
  hostname: System.get_env("POSTGRES_HOST"),
  database: "#{System.get_env("POSTGRES_DB")}-#{Mix.env()}",
  port: 5432,
  show_sensitive_data_on_connection_error: true,
  pool_size: 10

config :joken,
  hs256: [
    signer_alg: "HS256",
    key_octet: System.get_env("JWT_SECRET")
  ]

# Configures Elixir's Logger
config :logger, :console,
  format: "$time $metadata[$level] $message\n",
  metadata: [:pid]

# Use Jason for JSON parsing in Phoenix
config :phoenix, :json_library, Jason

# Import environment specific config. This must remain at the bottom
# of this file so it overrides the configuration defined above.
import_config "#{config_env()}.exs"
