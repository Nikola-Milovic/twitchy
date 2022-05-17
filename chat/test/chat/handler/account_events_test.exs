defmodule Chat.Handler.AccountEventsTest do
  use ExUnit.Case, async: true

  describe "handle_event/3" do
    test "when there is correct payload, will return :ok" do
      payload =
        "{\"type\":\"account_created\",\"payload\":\"{\\\"id\\\":1,\\\"email\\\":\\\"test@gmail.com\\\",\\\"username\\\":\\\"test username\\\"}\"}"

      assert :ok = Chat.Handler.AccountEvents.handle_event(payload)
    end
  end
end
