defmodule Chat.AMQP.AccountConsumerTest do
  use ExUnit.Case, async: true

  describe "handle_info :basic_deliver" do
    test "when there is :account_created event, will create account" do
      payload =
        "{\"type\":\"account_created\",\"payload\":\"{\\\"id\\\":1,\\\"email\\\":\\\"test@gmail.com\\\",\\\"username\\\":\\\"test username\\\"}\"}"

      Process.send(
        Chat.AMQP.AccountConsumer,
        {
          :basic_deliver,
          payload,
          %{delivery_tag: 1, redelivered: false}
        },
        []
      )
    end
  end
end
