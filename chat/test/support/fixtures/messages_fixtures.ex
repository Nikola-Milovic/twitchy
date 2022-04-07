defmodule Chat.MessagesFixtures do
  @moduledoc """
  This module defines test helpers for creating
  entities via the `Chat.Messages` context.
  """

  @doc """
  Generate a message.
  """
  def message_fixture(user, attrs \\ %{}) do
    {:ok, message} =
      attrs
      |> Enum.into(%{
        channel: "some channel",
        contents: "some contents",
        user_id: user.id
      })
      |> Chat.Messages.create_and_populate_message()

    message
  end
  def message_fixture_no_preload(user, attrs \\ %{}) do
    {:ok, message} =
      attrs
      |> Enum.into(%{
        channel: "some channel",
        contents: "some contents",
        user_id: user.id
      })
      |> Chat.Messages.create_message()

    message
  end


  @spec message_fixture_with_user(any) ::
          nil | [%{optional(atom) => any}] | %{optional(atom) => any}
  def message_fixture_with_user(attrs \\ %{}) do
    user = Chat.UsersFixtures.user_fixture()

    {:ok, message} =
      attrs
      |> Enum.into(%{
        channel: "some channel",
        contents: "some contents",
        user_id: user.id
      })
      |> Chat.Messages.create_and_populate_message()

    message
  end
end
