defmodule Chat.UsersFixtures do
  @moduledoc """
  This module defines test helpers for creating
  entities via the `Chat.Users` context.
  """

  @doc """
  Generate a user.
  """
  def user_fixture(attrs \\ %{}) do
    id = :rand.uniform(10000)
    
    {:ok, user} =
      attrs
      |> Enum.into(%{
        id: id,
        username: "some username"
      })
      |> Chat.Users.create_user()

    user
  end
end
