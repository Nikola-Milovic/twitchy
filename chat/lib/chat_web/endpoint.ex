defmodule ChatWeb.Endpoint do
  use Phoenix.Endpoint, otp_app: :chat

  socket "/socket", ChatWeb.UserSocket,
    websocket: true,
    longpoll: true

  # socket "/live", Phoenix.LiveView.Socket, websocket: [connect_info: [session: @session_options]]

  # Code reloading can be explicitly enabled under the
  # :code_reloader configuration of your endpoint.
  if code_reloading? do
    plug Phoenix.CodeReloader
    plug Phoenix.Ecto.CheckRepoStatus, otp_app: :chat
  end

  plug Plug.RequestId
  plug Plug.Telemetry, event_prefix: [:phoenix, :endpoint]

  plug Plug.Parsers,
    parsers: [:urlencoded, :multipart, :json],
    pass: ["*/*"],
    json_decoder: Phoenix.json_library()

  plug Plug.MethodOverride
  plug Plug.Head
  plug ChatWeb.Router
end
