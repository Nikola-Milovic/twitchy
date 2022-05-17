defmodule Chat.AMQP.AccountConsumer do
  use GenServer
  use AMQP
  require Logger

  @exchange "accounts_topic"
  @account_created_key "account.created"
  @queue "accounts_queue"
  @queue_error "#{@queue}_error"

  def start_link(_args) do
    GenServer.start_link(__MODULE__, [], name: __MODULE__)
  end

  @impl GenServer
  def init(_opts) do
    Chat.AMQP.Connection.request_channel(__MODULE__)
    Logger.log(:info, "Started AMQP account consumer")
    {:ok, nil}
  end

  def channel_available(chan) do
    GenServer.cast(__MODULE__, {:channel_available, chan})
  end

  @impl GenServer
  def handle_cast({:channel_available, chan}, _state) do
    Logger.log(:info, "Received channel, setting up consumer")
    account_created_queue(chan)
    bind_to_exchange(chan)
    # Limit unacknowledged messages to 10
    :ok = Basic.qos(chan, prefetch_count: 10)
    # Register the GenServer process as a consumer
    {:ok, _consumer_tag} = Basic.consume(chan, @queue)
    {:noreply, chan}
  end

  # Confirmation sent by the broker after registering this process as a consumer
  @impl GenServer
  def handle_info({:basic_consume_ok, %{consumer_tag: consumer_tag}}, chan) do
    Logger.debug("Registered as consumer")
    {:noreply, chan}
  end

  # Sent by the broker when the consumer is unexpectedly cancelled (such as after a queue deletion)
  @impl GenServer
  def handle_info({:basic_cancel, %{consumer_tag: consumer_tag}}, chan) do
    Logger.debug("Consumer cancelled")
    {:stop, :normal, chan}
  end

  # Confirmation sent by the broker to the consumer process after a Basic.cancel
  @impl GenServer
  def handle_info({:basic_cancel_ok, %{consumer_tag: consumer_tag}}, chan) do
    {:noreply, chan}
  end

  @impl GenServer
  def handle_info({:basic_deliver, payload, %{delivery_tag: tag, redelivered: redelivered}}, chan) do
    # You might want to run payload consumption in separate Tasks in production
    consume(chan, tag, redelivered, payload)
    {:noreply, chan}
  end

  defp account_created_queue(chan) do
    {:ok, _} = Queue.declare(chan, @queue, durable: true)

    # # Messages that cannot be delivered to any consumer in the main queue will be routed to the error queue
    # {:ok, _} =
    #   Queue.declare(chan, @queue,
    #     durable: true,
    #     arguments: [
    #       {"x-dead-letter-exchange", :longstr, ""},
    #       {"x-dead-letter-routing-key", :longstr, @queue_error}
    #     ]
    #   )
  end

  defp bind_to_exchange(chan) do
    :ok = Exchange.topic(chan, @exchange, durable: true)
    :ok = Queue.bind(chan, @queue, @exchange, routing_key: @account_created_key)
  end

  defp consume(channel, tag, redelivered, payload) do
    case Chat.Handler.AccountEvents.handle_event(payload) do
      :ok ->
        Basic.ack(channel, tag)
        :ok

      {:error, err} ->
        Basic.reject(channel, tag, requeue: not redelivered)
        {:error, err}
    end
  rescue
    # Requeue unless it's a redelivered message.
    # This means we will retry consuming a message once in case of exception
    # before we give up and have it moved to the error queue
    #
    # You might also want to catch :exit signal in production code.
    # Make sure you call ack, nack or reject otherwise consumer will stop
    # receiving messages.
    exception ->
      :ok = Basic.reject(channel, tag, requeue: not redelivered)
      Logger.error("Error consuming event #{payload}")
      IO.inspect(exception)
  end
end
