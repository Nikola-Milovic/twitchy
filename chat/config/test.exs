import Config

# # Configure your database
# #
# # The MIX_TEST_PARTITION environment variable can be used
# # to provide built-in test partitioning in CI environment.
# # Run `mix help test` for more information.
config :chat, Chat.Repo,
  # username: "postgres",
  # password: "postgres",
  # hostname: "localhost",
  # database: "chat_test#{System.get_env("MIX_TEST_PARTITION")}",
  pool: Ecto.Adapters.SQL.Sandbox

# pool_size: 10
config :logger, :console,
  format: "[$level]$metadata $message\n",
  metadata: [:file, :line]

# We don't run a server during test. If one is required,
# you can enable the server option below.
config :chat, ChatWeb.Endpoint,
  http: [ip: {127, 0, 0, 1}, port: 4002],
  secret_key_base: "ci+F6m++pZcZaKeSMRGGVwVGgl6eIPtS25wHNRnhf6/FW1nJPZs3QQ8zFcC+AW6/",
  server: false

# Print only warnings and errors during test
config :logger, level: :warn

# Initialize plugs at runtime for faster test compilation
config :phoenix, :plug_init_mode, :runtime
