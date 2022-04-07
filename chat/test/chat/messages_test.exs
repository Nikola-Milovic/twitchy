defmodule Chat.MessagesTest do
  use Chat.DataCase

  alias Chat.Messages

  describe "messages" do
    alias Chat.Messages.Message

    import Chat.MessagesFixtures

    @invalid_attrs %{channel: nil, contents: nil, user_id: nil}

    test "list_messages/0 returns all messages" do
      user = Chat.UsersFixtures.user_fixture()
      message1 = message_fixture_no_preload(user)
      message2 = message_fixture_no_preload(user)
      assert Messages.list_messages(user.id) == [message1, message2]
    end

    test "get_message!/1 returns the message with given id" do
      message = message_fixture_with_user()
      assert Messages.get_message!(message.id) == message
    end

    test "create_and_populate_message/1 with valid data creates a message" do
      user = Chat.UsersFixtures.user_fixture()

      valid_attrs = %{channel: "some channel", contents: "some contents", user_id: user.id}

      assert {:ok, %Message{} = message} = Messages.create_and_populate_message(valid_attrs)
      assert message.channel == "some channel"
      assert message.contents == "some contents"
      assert message.user_id == user.id
      assert message.user == user
    end

    test "create_and_populate_message/1 with invalid data returns error changeset" do
      assert {:error, %Ecto.Changeset{}} = Messages.create_and_populate_message(@invalid_attrs)
    end

    test "update_message/2 with valid data updates the message" do
      message = message_fixture_with_user()

      update_attrs = %{
        channel: "some updated channel",
        contents: "some updated contents"
      }

      assert {:ok, %Message{} = message} = Messages.update_message(message, update_attrs)
      assert message.channel == "some updated channel"
      assert message.contents == "some updated contents"
    end

    test "update_message/2 with invalid data returns error changeset" do
      message = message_fixture_with_user()
      assert {:error, %Ecto.Changeset{}} = Messages.update_message(message, @invalid_attrs)
      assert message == Messages.get_message!(message.id)
    end

    test "delete_message/1 deletes the message" do
      message = message_fixture_with_user()
      assert {:ok, %Message{}} = Messages.delete_message(message)
      assert_raise Ecto.NoResultsError, fn -> Messages.get_message!(message.id) end
    end

    test "change_message/1 returns a message changeset" do
      message = message_fixture_with_user()
      assert %Ecto.Changeset{} = Messages.change_message(message)
    end
  end
end
