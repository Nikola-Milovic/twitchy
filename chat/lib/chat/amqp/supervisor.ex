defmodule Chat.AMQP.Supervisor do
  use Supervisor
  require Logger

  # TODO try out https://github.com/meltwater/gen_rmq
  def start_link(_args) do
    Supervisor.start_link(__MODULE__, :ok, name: __MODULE__)
  end

  @impl true
  def init(:ok) do
    children = [
      Chat.AMQP.Connection,
      Chat.AMQP.AccountConsumer
    ]

    opts = [strategy: :one_for_one, name: __MODULE__]
    Logger.log(:info, "Started AMQP supervisor")
    Supervisor.init(children, opts)
  end
end
