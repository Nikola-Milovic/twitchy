defmodule Chat.Messages.Message do
  use Ecto.Schema
  import Ecto.Changeset

  schema "messages" do
    field(:channel, :string)
    field(:contents, :string)
    belongs_to(:user, Chat.Users.User)

    timestamps()
  end

  @doc false
  def changeset(message, attrs) do
    message
    |> cast(attrs, [:channel, :contents, :user_id])
    |> assoc_constraint(:user)
    |> validate_required([:user_id, :channel, :contents])
  end
end
