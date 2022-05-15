defmodule Chat.AMQP.AccountConsumer do
  use GenServer
  use AMQP
  require Logger

  def start_link do
    GenServer.start_link(__MODULE__, [], [])
  end

  @exchange "accounts_topic"
  @account_created_key "account.created"
  @queue "accounts_queue"
  @queue_error "#{@queue}_error"

  def init(_opts) do
    Chat.AMQP.ConsumerManager.request_channel(__MODULE__)
    {:ok, nil}
  end

  def channel_available(chan) do
    GenServer.cast(__MODULE__, {:channel_available, chan})
  end

  def handle_cast({:channel_available, chan}, _state) do
    Logger.log("debug", "Account consumer setup")
    setup_account_created_queue(chan)
    # Limit unacknowledged messages to 10
    :ok = Basic.qos(chan, prefetch_count: 10)
    # Register the GenServer process as a consumer
    {:ok, _consumer_tag} = Basic.consume(chan, @queue)
    {:noreply, chan}
  end

  # Confirmation sent by the broker after registering this process as a consumer
  def handle_info({:basic_consume_ok, %{consumer_tag: consumer_tag}}, chan) do
    {:noreply, chan}
  end

  # Sent by the broker when the consumer is unexpectedly cancelled (such as after a queue deletion)
  def handle_info({:basic_cancel, %{consumer_tag: consumer_tag}}, chan) do
    {:stop, :normal, chan}
  end

  # Confirmation sent by the broker to the consumer process after a Basic.cancel
  def handle_info({:basic_cancel_ok, %{consumer_tag: consumer_tag}}, chan) do
    {:noreply, chan}
  end

  def handle_info({:basic_deliver, payload, %{delivery_tag: tag, redelivered: redelivered}}, chan) do
    # You might want to run payload consumption in separate Tasks in production
    consume(chan, tag, redelivered, payload)
    {:noreply, chan}
  end

  defp setup_account_created_queue(chan) do
    {:ok, _} = Queue.declare(chan, @queue_error, durable: true)

    # Messages that cannot be delivered to any consumer in the main queue will be routed to the error queue
    {:ok, _} =
      Queue.declare(chan, @queue,
        durable: true,
        arguments: [
          {"x-dead-letter-exchange", :longstr, ""},
          {"x-dead-letter-routing-key", :longstr, @queue_error}
        ]
      )

    :ok = Exchange.topic(chan, @exchange, durable: true)
    :ok = Queue.bind(chan, @queue, @exchange, routing_key: @account_created_key)
  end

  defp consume(channel, tag, redelivered, payload) do
    IO.puts("Received event #{payload}")
    :ok = Basic.ack(channel, tag)

    # if number <= 10 do
    #   :ok = Basic.ack(channel, tag)
    #   IO.puts("Consumed a #{number}.")
    # else
    #   :ok = Basic.reject(channel, tag, requeue: false)
    #   IO.puts("#{number} is too big and was rejected.")
    # end
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
      IO.puts("Error consuming event #{payload}, on channel #{channel}")
  end
end
