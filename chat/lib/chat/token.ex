defmodule Chat.Token do
  use Joken.Config, default_signer: :hs256

  @impl Joken.Config
  def token_config do
    default_claims(
      iss: "twitchy",
      default_exp: Joken.current_time() + 3000,
      skip: [:aud]
    )
    |> add_claim("sub", nil, &is_valid_id/1)
  end

  def is_valid_id(_id) do
    true
  end
end
