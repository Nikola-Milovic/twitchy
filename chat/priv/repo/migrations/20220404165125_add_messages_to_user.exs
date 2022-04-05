defmodule Chat.Repo.Migrations.AddMessagesToUser do
  use Ecto.Migration

  def change do
    alter table(:users) do
      add(:messages, references(:messages))
    end
  end
end
