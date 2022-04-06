defmodule ChatWeb.RoomChannelAuthTest do
  use ChatWeb.ChannelCase

  setup do
    {:ok, _, socket} =
      ChatWeb.UserSocket
      |> socket("user_id", %{some: :assign})
      |> subscribe_and_join(ChatWeb.RoomChannel, "room:lobby")

    %{socket: socket}
  end

  test "ping replies with status ok", %{socket: socket} do
    ref = push(socket, "ping", %{"hello" => "there"})
    assert_reply ref, :ok, %{"hello" => "there"}
  end

  test "shout broadcasts to room:lobby", %{socket: socket} do
    push(socket, "shout", %{"hello" => "all"})
    assert_broadcast "shout", %{"hello" => "all"}
  end

  test "broadcasts are pushed to the client", %{socket: socket} do
    broadcast_from!(socket, "broadcast", %{"some" => "data"})
    assert_push "broadcast", %{"some" => "data"}
  end

  test "sending message to a room", %{socket: socket} do
    push(socket, "send_message", %{"hello" => "all"})
    assert_broadcast "receive_message", %{"hello" => "all"}
  end
end
