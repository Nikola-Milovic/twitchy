defmodule Chat.Handler.AccountEvents do
  alias AMQP.Basic
  require Logger

  def handle_event(raw_event) do
    with {:ok, parsed_event} = decode_event(raw_event),
         {:ok, user} = do_handle_event(parsed_event) do
      :ok
    else
      {:error, error} ->
        {:error, error}
    end
  end

  defp do_handle_event(%{
         type: :account_created,
         payload: %{id: id, email: email, username: username}
       }) do
    Chat.Users.create_user(%{
      id: id,
      username: username
    })
  end

  defp do_handle_event(ev) do
    IO.inspect(ev)
    {:error, "Unknown event type"}
  end

  # Key to atoms is dangerous here, but our queues should be protected/ safe
  defp decode_event(raw_event) when is_binary(raw_event) do
    with {:ok, %{type: type, payload: payload_string}} =
           Poison.decode(raw_event, %{keys: :atoms}),
         {:ok, payload} = Poison.decode(payload_string, %{keys: :atoms}) do
      data = %{type: event_type_to_atom(type), payload: payload}
      Logger.log(:debug, "Handled event #{inspect(data)}")
      {:ok, data}
    else
      {:error, err} ->
        {:error, err}
    end
  end

  defp decode_event(unknown_event) do
    IO.inspect(unknown_event)
    {:error, "Unknown event data type"}
  end

  defp event_type_to_atom(str) when str in ~w(account_created),
    do: String.to_atom(str)
end
