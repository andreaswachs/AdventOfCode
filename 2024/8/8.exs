defmodule Solution do
  def run do
    input = IO.read(:stdio, :eof)

    antennas =
      String.split(input, "\n")
      |> Enum.map(&String.trim/1)
      |> Enum.with_index()
      |> Enum.map(fn {line, y} ->
        line
        |> String.graphemes()
        |> Enum.with_index()
        |> Enum.filter(fn {char, _} -> char != "." end)
        |> Enum.map(fn {c, x} -> {c, {y, x}} end) # y is really x and the other way around - im confused
      end)
      |> List.flatten()
      |> Enum.reduce(%{}, fn {c, {x, y}}, acc ->
        Map.update(acc, c, [{x, y}], fn list -> [{x, y} | list] end)
      end)

    antennas
    |> Map.keys()
    |> Enum.map(&gen_all_pairs(&1, Map.get(antennas, &1)))
    |> IO.inspect()
  end

  @doc """
    Generates all combinations of antenna coordinates for antennas with the same character
  """
  def gen_all_pairs(antenna, antennas) do
    for [coord1, coord2] <- combinations(antennas, 2) do
      {coord1, coord2}
    end
    |> (fn pairs -> {antenna, pairs} end).()
  end

  defp combinations(_, 0), do: [[]]
  defp combinations([], _), do: []
  defp combinations([h|t], m) do
    (for l <- combinations(t, m-1), do: [h|l]) ++ combinations(t, m)
  end
end

Solution.run()
