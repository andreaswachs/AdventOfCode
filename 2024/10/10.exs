defmodule Solution do
  def run do
    input = IO.read(:stdio, :eof) |> String.split("\n")
    init_map()

    trailhead_pos =
      input
      |> Enum.with_index()
      |> Enum.map(fn {row, y} ->
        row
        |> String.graphemes()
        |> Enum.with_index()
        |> Enum.map(fn {symbol, x} ->
          tile = if(symbol == ".", do: "99", else: symbol) |> Integer.parse() |> elem(0)
          update_map({x, y}, tile)

          case tile do
            0 -> {x, y}
            _ -> :skip
          end
        end)
      end)
      |> List.flatten()
      |> Enum.reject(&(&1 == :skip))

    _part1 =
      part1(trailhead_pos)
      |> IO.inspect(label: "Part 1")

    _part2 =
      part2(trailhead_pos)
      |> IO.inspect(label: "Part 2")
  end

  def part1(trailheads) do
    trailheads
    |> Enum.map(fn p ->
      do_part1(p, p) 
      |> List.flatten() 
      |> Enum.uniq() 
      |> Enum.count()
    end)
    |> Enum.sum()
  end

  defp do_part1(start_pos, curr_pos) do
    value = get_map(curr_pos)

    case value do
      9 ->
        {start_pos, curr_pos}

      _ ->
        curr_pos
        |> get_neighbouring_positions()
        |> Enum.filter(fn {_, _, v} -> value + 1 == v end)
        |> Enum.map(fn {x, y, _} -> {x, y} end)
        |> Enum.map(&do_part1(start_pos, &1))
    end
  end
  
  def part2(trailheads) do
    trailheads
    |> Enum.map(&do_part2/1)
    |> Enum.sum()
  end

  defp do_part2(curr_pos) do
    value = get_map(curr_pos)

    case value do
      9 ->
        1

      _ ->
        curr_pos
        |> get_neighbouring_positions()
        |> Enum.filter(fn {_, _, v} -> value + 1 == v end)
        |> Enum.map(fn {x, y, _} -> {x, y} end)
        |> Enum.map(&do_part2/1)
        |> Enum.sum()
    end

  end

  defp get_neighbouring_positions({x, y}) do
    # y axis is inverted
    deltas =
      [
        # left
        {-1, 0},
        # up
        {0, -1},
        # right
        {1, 0},
        # down
        {0, 1}
      ]

    deltas
    |> Enum.map(fn {dx, dy} -> {dx + x, dy + y} end)
    |> Enum.map(fn {x, y} -> {x, y, get_map({x, y})} end)
    |> Enum.reject(&(elem(&1, 2) == :not_found))
  end

  def init_map do
    :ets.new(:map, [:set, :named_table])
    :ok
  end

  def init_queue do
    :ets.new(:queue, [:set, :named_table])
    :ok
  end

  def update_map({_x, _y} = pos, tile) do
    :ets.insert(:map, {pos, tile})
    :ok
  end

  def get_map({x, y} = pos) do
    case :ets.lookup(:map, pos) do
      [{{^x, ^y}, tile}] -> tile
      _ -> :not_found
    end
  end

  def enqueue({_x, _y} = position) do
    new_queue = Enum.concat([position], get_queue())
    :ets.insert(:queue, {:queue, new_queue})
    :ok
  end

  def dequeue do
    case get_queue() do
      [x | tail] ->
        :ets.insert(:queue, {:queue, tail})
        x

      _ ->
        :empty
    end
  end

  def get_queue do
    case :ets.lookup(:queue, :queue) do
      [{:queue, queue}] -> queue
      _ -> []
    end
  end
end

Solution.run()
