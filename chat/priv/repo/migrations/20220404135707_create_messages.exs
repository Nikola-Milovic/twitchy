defmodule Chat.Repo.Migrations.CreateMessages do
  use Ecto.Migration

  def change do
    create table(:messages) do
      add(:user_id, references(:users, on_delete: :delete_all))
      add(:contents, :string)
      add(:channel, :string)

      timestamps()
    end

    create index(:messages, [:user_id])
  end
end
