defmodule Chat.Handler.AccountEvents do
  alias AMQP.Basic
  require Logger

  def handle_event(event) do
    {:ok, data} = decode_event(event)
    IO.inspect(data)
    :ok
  end

  def handle_event(channel, tag, %{
        type: :account_created,
        payload: %{id: id, email: email, username: username}
      }) do
    :ok = Basic.ack(channel, tag)

    # if number <= 10 do
    #   :ok = Basic.ack(channel, tag)
    #   IO.puts("Consumed a #{number}.")
    # else
    #   :ok = Basic.reject(channel, tag, requeue: false)
    #   IO.puts("#{number} is too big and was rejected.")
    # end
  end

  defp decode_event(event) do
    with {:ok, %{"type" => type, "payload" => payload_string}} = Poison.decode(event),
         {:ok, payload} = Poison.decode(payload_string) do
      data = %{type: type, payload: payload}
      Logger.log(:info, "Handled event #{data}")
      {:ok, data}
    else
      {:error, err} ->
        Logger.log(:error, "Error decoding event #{event}, #{err}")
        {:error, err}
    end
  end
end
