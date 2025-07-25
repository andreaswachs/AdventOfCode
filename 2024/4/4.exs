defmodule Solution do

  @type grid :: [[String.t()]]
  @type location :: {non_neg_integer(), non_neg_integer()}
  @type location_with_direction :: {non_neg_integer(), non_neg_integer(), atom()}
  @type direction :: atom()
  @type token :: String.t()
  @type word :: String.t()

  defp debug(what, opts \\ []) do
    if System.argv() |> Enum.find(nil, fn a -> a == "debug" end) != nil do
      IO.inspect(what, opts)
    else
      what
    end
  end

  def run do
    grid =
      read_all()
      |> String.split("\n")
      |> Enum.map(&String.split(&1, ""))
      |> Enum.map(&Enum.drop(&1, 1))
      |> Enum.map(&Enum.drop(&1, -1))

    _part1 = grid |> part1("XMAS") |> IO.inspect(label: "Part 1")

    _part2 = grid |> part2() |> IO.inspect(label: "Part 2")
  end

  @spec locations_for(grid :: grid, letter :: token) :: [location()]
  defp locations_for(grid, letter) do
    grid
      |> Enum.with_index()
      |> Enum.flat_map(fn {row, y} ->
        row
        |> Enum.with_index()
        |> Enum.filter(fn {cell, _} -> cell == letter end)
        |> Enum.map(fn {_, x} -> {y, x} end)
      end)
  end

  @spec part2(grid :: grid) :: non_neg_integer()
  def part2(grid) do
      possible_crosses = locations_for(grid, "A")
      Enum.map(possible_crosses, fn location ->
        check_cross(grid, location)
      end)
      |> Enum.filter(&(&1))
      |> Enum.count()
  end

  @spec check_cross(grid :: grid(), location :: location()) :: boolean()
  defp check_cross(grid, location) do
    # First we filter out locations that clip outside the bounds of the grid
    yb = length(grid)
    xb = length(Enum.at(grid, 0))

    if will_be_inside_bounds?(yb, xb, location) do
      case get_cross(grid, location) do
        "MASMAS" -> true
        "MASSAM" -> true
        "SAMSAM" -> true
        "SAMMAS" -> true
        _ -> false
      end
    else
      false
    end
  end

  @spec get_cross(grid :: grid(), location :: location()) :: String.t()
  defp get_cross(grid, {y, x}) do
    [
      [at(grid, {y-1, x-1}), at(grid, {y, x}), at(grid, {y+1, x+1})],
      [at(grid, {y+1, x-1}), at(grid, {y, x}), at(grid, {y-1, x+1})]
    ]
    |> Enum.map(fn row ->
      row
      |> Enum.filter(&(&1 != nil))
      |> Enum.join()
    end)
    |> Enum.join()
    |> debug(label: "words")
  end

  @spec will_be_inside_bounds?(yb :: non_neg_integer(), xb :: non_neg_integer(), location :: location()) :: boolean()
  defp will_be_inside_bounds?(yb, xb, {y, x}) do
    x >= 1 and y >= 1 and y < (yb-1) and x < (xb-1)
  end

  @spec part1(grid :: grid(), word :: word) :: non_neg_integer()
  defp part1(grid, word) do
    first_token = word |> String.graphemes() |> hd()
    grid |> debug()
    starting_locations = locations_for(grid, first_token) |> debug(label: "starting locations")
    following_tokens = String.graphemes(word) |> tl() |> debug(label: "graphemes")
    all_directions = possible_directions() |> Enum.map(&elem(&1, 2)) |> debug()

    starting_locations
    |> Enum.flat_map(fn location ->
      Enum.map(all_directions, fn direction -> try_traverse(grid, following_tokens, location, direction) end)
    end)
    |> Enum.filter(&(&1))
    |> Enum.count()
  end

  @spec try_traverse(grid :: grid(), tokens :: [String.t()], location :: location(), direction :: atom()) :: boolean()
  defp try_traverse(grid, tokens, location, direction) do
    if length(tokens) == 0 do
      true
    else
      next_token = hd(tokens)
      next_location = generate_next_location(grid, location, direction) |> debug(label: "next locations")
      n = at(grid, next_location)
      if n == next_token do
        try_traverse(grid, tl(tokens), next_location, direction)
      else
        false
      end
    end
  end

  @spec at(grid :: grid(), location :: location() | nil) :: String.t() | nil
  defp at(_grid, nil) do
    nil
  end

  defp at(grid, {y, x, _}) do
    Enum.at(Enum.at(grid, y), x)
  end

  defp at(grid, {y, x}) do
    Enum.at(Enum.at(grid, y), x)
  end

  @spec possible_directions() :: [{integer(), integer(), atom()}]
  defp possible_directions() do
    [
      {0, 1, :forward},
      {0, -1, :backward},
      {1, 0, :down},
      {-1, 0, :up},
      {1, 1, :down_right},
      {-1, 1, :down_left},
      {1, -1, :up_right},
      {-1, -1, :up_left}
    ]
  end

  @spec generate_next_location(grid :: grid(), location :: location(), direction :: direction()) :: location_with_direction() | nil
  defp generate_next_location(grid, location, direction) do

    debug(length(grid), label: "grid length")
    # There might be a third element in location, thats why we dont
    # destruct it with pattern matching
    y = elem(location, 0)
    x = elem(location, 1)

    {dy, dx, _} = possible_directions()
      |> Enum.find(fn {_, _, dirr} -> dirr == direction end)

    r = {ny, nx, _} = {y + dy, x + dx, direction}

    cond do
      ny < 0 or nx < 0 -> nil
      ny >= length(grid) or nx >= length(Enum.at(grid, 0)) -> nil
      true -> r
    end
  end

  defp read_all() do
    IO.read(:stdio, :all)
  end
end

Solution.run()
