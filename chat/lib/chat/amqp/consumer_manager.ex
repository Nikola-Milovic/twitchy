defmodule Chat.AMQP.ConsumerManager do
  use AMQP
  use Supervisor
  require Logger
  # TODO try out https://github.com/meltwater/gen_rmq
  def start_link do
    GenServer.start_link(__MODULE__, :ok, name: __MODULE__)
  end

  def init(:ok) do
    children = [
      Chat.AMQP.AccountConsumer
    ]

    Chat.AMQP.ConnectionManager.request_conn(self)

    opts = [strategy: :one_for_one, name: Chat.AMQP.ConsumerManager]
    Supervisor.start_link(children, opts)
  end

  def connect(conn) do
    GenServer.cast(__MODULE__, {:connected, conn})
  end

  def handle_cast({:connected, conn}) do
    Logger.log("debug", "Consumer manager connected to AMQP")
    {:noreply, {conn, %{}}}
  end

  def request_channel(consumer) do
    GenServer.cast(__MODULE__, {:chan_request, consumer})
  end

  def handle_cast({:chan_request, consumer}, {conn, channel_mappings}) do
    new_mapping = store_channel_mapping(conn, consumer, channel_mappings)
    channel = Map.get(new_mapping, consumer)
    consumer.channel_available(channel)

    {:noreply, {conn, new_mapping}}
  end

  defp store_channel_mapping(conn, consumer, channel_mappings) do
    Map.put_new_lazy(channel_mappings, consumer, fn -> create_channel(conn) end)
  end

  defp create_channel(conn) do
    {:ok, chan} = Channel.open(conn)
    chan
  end
end
