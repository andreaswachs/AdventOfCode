defmodule Solution do

  def run do
    :ets.new(:data, [:set, :named_table])
    :ets.new(:visited, [:set, :named_table])

    grid = 
      IO.read(:stdio, :eof) 
      |> String.split("\n") 

    grid |> load_to_ets()


    grid |> part2() |> Enum.uniq() |> IO.inspect()
    # part1()
    # print_grid(length(grid)-1 , length(Enum.at(grid, 0) |> String.graphemes())-1)
  end

  def part1() do
    [position()]
    |> do_part1()
    |> Enum.uniq()
    |> Enum.count()
    |> IO.inspect(label: "part1()")
  end

  def all_dots_pos(grid) do
    grid
    |> Enum.with_index()
    |> Enum.flat_map(fn {row, y} ->
      row
      |> String.graphemes()
      |> Enum.with_index()
      |> Enum.flat_map(fn {c, x} -> 
        case c do
          "." -> [{y, x}]
          _ -> []
        end
      end)
    end)
  end

  defp part2(grid) do
    grid 
    |> all_dots_pos() 
    |> Enum.with_index() 
    |> Enum.map(fn {pos, i} ->
      :ets.insert(:data, {pos, '#'})
      res = do_part2(i, 0)
      :ets.insert(:data, {pos, '.'})
      res
    end)
  end

  defp do_part2(iteration, recurse_it) do
    dir = direction()

    case recurse_it do
      0 -> 
        reset_position()
        reset_direction()
      _ -> :noop
    end

    case position() |> next_pos(dir) |> get_i(dir, iteration) do
      :goal -> :escaped |> IO.inspect(label: "escaped at #{recurse_it}")
      :loop -> :loop
      pos -> 
        case pos do
          {_, "#"} -> 
            direction() |> rotate() |> set_direction()
          _ -> :ignore
        end
        dir = direction()
        position() 
        |> next_pos(dir) 
        |> set_position() 
        |> mark_visited(dir, iteration)
        do_part2(iteration, recurse_it + 1)
    end
  end

  defp do_part1(moves) do
    case position() |> next_pos(direction()) |> get() do
      :goal -> moves
      pos -> 
        case pos do
          {_, "#"} -> 
            direction() |> rotate() |> set_direction()
          _ -> :ignore
        end
        [position() |> next_pos(direction()) |> set_position() |> mark_visited(direction()) | moves]
        |> do_part1()

    end
  end

  def get(pos) do
    case :ets.lookup(:data, pos) do
      [p] -> p
      _ -> :goal
    end
  end

  def get_i({y, x} = pos, dir, iteration \\ 0) do
    case :ets.lookup(:visited, {y, x, dir, iteration}) do
      [{^y, ^x, ^dir, ^iteration, _}] -> :loop |> IO.inspect(label: "LOOOP!")
       _ -> get(pos)
    end
  end

  def position() do
    case :ets.lookup(:data, :position) do
      [{:position, pos}] -> pos
      _ -> nil # we should never reach this case
    end
  end

  def set_position({_y, _x} = pos) do
    :ets.insert(:data, {:position, pos})
    pos
  end

  def reset_position() do
    case :ets.lookup(:data, :start_position) do
      [{:start_position, pos}] -> set_position(pos)
        _ -> IO.puts("reset position failed")
    end
  end

  def direction() do
    case :ets.lookup(:data, :direction) do
      [{:direction, dir}] -> dir
      _ -> nil # we should never reach this case
    end
  end

  def set_direction(dir) do
    :ets.insert(:data, {:direction, dir})
  end

  def reset_direction() do
    case :ets.lookup(:data, :start_direction) do
      [{:start_direction, dir}] -> :ets.insert(:data, {:direction, dir})
      _ -> IO.puts("reset direction failed")
    end
  end

  def next_pos({y, x}, dir) do
    case dir do
      :up -> {y - 1, x}
      :right -> {y, x + 1}
      :left -> {y, x - 1}
      :down -> {y + 1, x}
    end
  end

  def rotate(dir) do
    case dir do
      :up -> :right
      :right -> :down
      :down -> :left
      :left -> :up
    end
  end


  def load_to_ets(grid) do
    grid
    |> Enum.with_index()
    |> Enum.each(fn {row, y} ->
      row
      |> String.graphemes()
      |> Enum.with_index()
      |> Enum.each(fn {c, x} -> 
        :ets.insert(:data, {{y, x}, c})
        case c do
          "^" -> 
            :ets.insert(:data, {:position, {y, x}})
            :ets.insert(:data, {:start_position, {y, x}})
            :ets.insert(:data, {:direction, :up})
            :ets.insert(:data, {:start_direction, :up})
          "v" ->
            :ets.insert(:data, {:position, {y, x}})
            :ets.insert(:data, {:start_position, {y, x}})
            :ets.insert(:data, {:direction, :down})
            :ets.insert(:data, {:start_direction, :down})
          ">" ->
            :ets.insert(:data, {:position, {y, x}})
            :ets.insert(:data, {:start_position, {y, x}})
            :ets.insert(:data, {:start_direction, :right})
          "<" ->
            :ets.insert(:data, {:position, {y, x}})
            :ets.insert(:data, {:start_position, {y, x}})
            :ets.insert(:data, {:start_direction, :left})
          _ -> true
        end
      end)
    end)
  end

  defp mark_visited({y, x} = pos, dir, iteration \\ 0) do
    s = {y, x, dir, iteration}
    IO.inspect("mark_visited(#{y}, #{x}, #{dir}, #{iteration}")

    :ets.insert(:data, {pos, "X"})
    :ets.insert(:visited, {s, :visited})
    pos
  end

  def print_grid(my, mx) do
    for y <- 0..my do
      for x <- 0..mx do
        case :ets.lookup(:data, {y, x}) do
          [{_, c}] -> c
          _ -> "@"
        end
      end
      |> Enum.join("")
    end
    |> Enum.join("\n")
    |> IO.puts()
  end
end

Solution.run()
