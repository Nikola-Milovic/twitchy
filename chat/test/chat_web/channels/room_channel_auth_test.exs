defmodule ChatWeb.RoomChannelAuthTest do
  use ChatWeb.ChannelCase

  setup do
    user = Chat.UsersFixtures.user_fixture()

    {:ok, _, socket} =
      ChatWeb.UserSocket
      |> socket("user_id", %{user_id: user.id})
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
    push(socket, "send_message", %{"contents" => "hello world", "channel" => "testchannel"})

    assert_broadcast "receive_message", expected_message

    assert expected_message.user_id == socket.assigns[:user_id]
    assert expected_message.contents == "hello world"
    assert expected_message.channel == "testchannel"

    saved_message = Chat.Messages.get_message!(expected_message.id)

    assert saved_message.user_id == expected_message.user_id
    assert saved_message.contents == expected_message.contents
    assert saved_message.channel == expected_message.channel
  end

  test "sending invalid message to a room will return error", %{socket: socket} do
    ref = push(socket, "send_message", %{"contents" => "hello world"})

    assert_reply ref, :error, %{reason: "invalid message format"}
  end
end
