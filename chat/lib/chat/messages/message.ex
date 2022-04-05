defmodule Chat.Messages.Message do
  use Ecto.Schema
  import Ecto.Changeset

  schema "messages" do
    field(:channel, :string)
    field(:contents, :string)
    belongs_to(:user, Chat.Users.User)

    timestamps()
  end

  @spec changeset(
          {map, map}
          | %{
              :__struct__ => atom | %{:__changeset__ => map, optional(any) => any},
              optional(atom) => any
            },
          :invalid | %{optional(:__struct__) => none, optional(atom | binary) => any}
        ) :: Ecto.Changeset.t()
  @doc false
  def changeset(message, attrs) do
    message
    |> cast(attrs, [:channel, :contents, :user_id])
    |> assoc_constraint(:user)
    |> validate_required([:user_id, :channel, :contents])
  end
end
