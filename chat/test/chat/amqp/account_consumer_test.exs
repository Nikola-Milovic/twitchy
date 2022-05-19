defmodule Chat.AMQP.AccountConsumerTest do
  use Chat.DataCase

  defp generate_account_created_event() do
    id = :rand.uniform(10_000)
    email = "test#{id}@gmail.com"
    username = "test username#{id}"

    {id, email, username,
     "{\"type\":\"account_created\",\"payload\":\"{\\\"id\\\":#{id},\\\"email\\\":\\\"#{email}\\\",\\\"username\\\":\\\"#{username}\\\"}\"}"}
  end

  describe "handle_info :basic_deliver" do
    test "consuming :account_created event, will create user entry" do
      # GIVEN
      {id, _email, username, payload} = generate_account_created_event()

      # WHEN
      Process.send(
        Chat.AMQP.AccountConsumer,
        {
          :basic_deliver,
          payload,
          %{delivery_tag: 1, redelivered: false}
        },
        []
      )

      # Wait for GenServer cast to finish
      _ = :sys.get_state(Chat.AMQP.AccountConsumer)

      # SHOULD
      user = Chat.Users.get_user!(id)

      assert user.id == id
      assert user.username == username
    end
  end
end
