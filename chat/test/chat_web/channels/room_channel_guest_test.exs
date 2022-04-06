defmodule ChatWeb.RoomChannelGuestTest do
  use ChatWeb.ChannelCase

  setup do
    {:ok, _, socket} =
      ChatWeb.UserSocket
      |> socket("guest", %{guest: true})
      |> subscribe_and_join(ChatWeb.RoomChannel, "room:lobby")

    %{socket: socket}
  end

  test "guests cannot send messages to channels", %{socket: socket} do
    ref = push(socket, "send_message", %{"hello" => "all"})
    assert_reply ref, :error, %{reason: "unauthorized"}
  end

end
