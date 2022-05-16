defmodule Chat.AMQP.Connection do
  use GenServer
  use AMQP
  require Logger

  def start_link(_args) do
    GenServer.start_link(__MODULE__, :ok, name: __MODULE__)
  end

  @impl true
  def init(:ok) do
    establish_new_connection()
  end

  def request_channel(consumer) do
    GenServer.cast(__MODULE__, {:chan_request, consumer})
  end

  @impl true
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

  # TODO add try limit
  defp establish_new_connection() do
    case AMQP.Connection.open(connection_options()) do
      {:ok, conn} ->
        Logger.log(:info, "Connected to AMQP")
        Process.link(conn.pid)
        {:ok, {conn, %{}}}

      {:error, reason} ->
        Logger.info("failed for #{inspect(reason)}")
        :timer.sleep(5000)
        establish_new_connection()
    end
  end

  defp connection_options() do
    user = System.get_env("RABBITMQ_USER")
    password = System.get_env("RABBITMQ_PASSWORD")
    host = System.get_env("RABBITMQ_HOST")
    port = System.get_env("RABBITMQ_PORT")

    "amqp://#{user}:#{password}@#{host}:#{port}"
  end
end
