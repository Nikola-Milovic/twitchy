defmodule Chat.AMQP.ConnectionManager do
  use Supervisor
  require Logger

  def start_link(opts \\ []) do
    GenServer.start_link(__MODULE__, :ok, name: __MODULE__)
  end

  def init(:ok) do
    children = [
      Chat.AMQP.ConsumerManager
    ]

    opts = [strategy: :one_for_one, name: Chat.AMQP.ConnectionManager]
    Supervisor.start_link(children, opts)
    establish_new_connection()
  end

  # TODO add try limit
  defp establish_new_connection() do
    case AMQP.Connection.open(connection_options()) do
      {:ok, conn} ->
        Logger.log("debug", "Connected to AMQP")
        Process.link(conn.pid)
        {:ok, {conn, %{}}}

      {:error, reason} ->
        IO.puts("failed for #{inspect(reason)}")
        :timer.sleep(5000)
        establish_new_connection()
    end
  end

  def connection_options() do
    user = System.get_env("RABBITMQ_USER")
    password = System.get_env("RABBITMQ_PASSWORD")
    host = System.get_env("RABBITMQ_HOST")
    port = System.get_env("RABBITMQ_PORT")

    "amqp://#{user}:#{password}@#{host}:#{port}/"
  end

  def request_conn(consumer) do
    GenServer.cast(__MODULE__, {:conn_request, consumer})
  end

  def handle_cast({:conn_request, consumer}, {conn}) do
    consumer.connect(conn)

    {:noreply, {conn}}
  end
end
