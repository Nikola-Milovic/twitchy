defmodule ChatWeb.RoomChannel do
  use ChatWeb, :channel

  @impl true
  def join("room:lobby", _params, socket) do
    {:ok, socket}
  end

  def join("room:" <> _private_room_id, _params, _socket) do
    {:error, %{reason: "unauthorized"}}
  end

  # Channels can be used in a request/response fashion
  # by sending replies to requests from the client
  @impl true
  def handle_in("ping", payload, socket) do
    {:reply, {:ok, payload}, socket}
  end

  # It is also common to receive messages from the client and
  # broadcast to everyone in the current topic (messaging:lobby).
  @impl true
  def handle_in("shout", payload, socket) do
    broadcast(socket, "shout", payload)
    {:noreply, socket}
  end

  @impl true
  def handle_in("send_message", %{"channel" => chan, "contents" => conts}, socket) do
    # Check if user is guest
    case Map.get(socket.assigns, :guest, false) do
      true ->
        {:reply, {:error, %{reason: "unauthorized"}}, socket}

      false ->
        message = %{
          "channel" => chan,
          "contents" => conts,
          "user_id" => socket.assigns[:user_id]
        }

        case Chat.Messages.create_and_populate_message(message) do
          {:ok, message} ->
            broadcast(socket, "receive_message", message)

          {:error, reason} ->
            {:reply, {:error, %{reason: reason}}}
        end

        {:noreply, socket}
    end
  end

  # Invalid message format
  @impl true
  def handle_in("send_message", _, socket) do
    {:reply, {:error, %{reason: "invalid message format"}}, socket}
  end

  # # Add authorization logic here as required.
  # defp authorized?(_payload) do
  #   true
  # end
end
