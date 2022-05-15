defmodule Chat.Users.User do
  @moduledoc """
  This is a Users module, containing all of the functionality needed to work with User entities (database models)
  """
  use Ecto.Schema
  import Ecto.Changeset

  @primary_key {:id, :id, autogenerate: false}
  schema "users" do
    field(:username, :string)
    has_many(:messages, Chat.Messages.Message)

    timestamps()
  end

  @doc false
  def changeset(user, attrs) do
    user
    |> cast(attrs, [:id, :username])
    |> unique_constraint(:username, message: "Username is already taken")
    |> unique_constraint(:id, message: "User with this ID already exists")
    |> validate_required([:username, :id])
  end
end
