defmodule Chat.UsersFixtures do
  @moduledoc """
  This module defines test helpers for creating
  entities via the `Chat.Users` context.
  """

  @doc """
  Generate a user.
  """
  def user_fixture(attrs \\ %{}) do
    id = :rand.uniform(10_000)

    {:ok, user} =
      attrs
      |> Enum.into(%{
        id: id,
        username: "some username"
      })
      |> Chat.Users.create_user()

    user
  end

  def generate_jwt(user) do
    {:ok, jwt, _claims} = Chat.Token.generate_and_sign(%{"sub" => user.id})

    jwt
  end
end
