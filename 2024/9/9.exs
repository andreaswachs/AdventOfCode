defmodule Solution do
  @type file_size :: non_neg_integer()
  @type id :: non_neg_integer()
  @type file :: {:file, id(), file_size()}

  def run() do
    input =
      IO.read(:stdio, :eof)
      |> String.graphemes()
      |> Enum.map(fn n -> Integer.parse(n) |> elem(0) end)

    _part1 =
      input
      |> gen_disk_map()
      |> format_part1()
      |> part1()
      |> IO.inspect(label: "Part 1")
  end

  def part1(mapping) do
    mapping
    |> Enum.with_index()
    |> Enum.map(fn {e, idx} ->
      case e do
        {:file, id, _} -> id * idx
        _ -> 0
      end
    end)
    |> Enum.sum()
  end

  defp format_part2(mapping) do
  end

  defp scan_blocks(mapping) do
    mapping
    |> Enum.reduce({[], nil}, fn
      {key, num}, {acc, nil} ->
        {acc ++ [{key, num}], num}

      {key, num}, {acc, prev_num} when num <= prev_num ->
        {acc ++ [{key, num}], num}

      {key, num}, {acc, prev_num} ->
        if prev_num + 1 != num do
          {acc ++ [{key, num}], num}
        else
          {acc, num}
        end
    end)
    |> elem(0)
  end

  def scan_free_spaces(mapping) do
    filtered_mapping =
      mapping
      |> Enum.with_index()
      |> Enum.filter(&(elem(&1, 0) == :free_space))

    filtered_mapping
    |> Enum.reduce({[], nil}, fn
      {key, num}, {acc, nil} ->
        {acc ++ [{key, num}], num}

      {key, num}, {acc, prev_num} when num <= prev_num ->
        {acc ++ [{key, num}], num}

      {key, num}, {acc, prev_num} ->
        if prev_num + 1 != num do
          {acc ++ [{key, num}], num}
        else
          {acc, num}
        end
    end)
    |> elem(0)
    |> Enum.map(&elem(&1, 1))
    |> Enum.map(fn idx ->
      {idx,
       mapping
       |> Enum.drop(idx)
       |> Enum.take_while(&(&1 == :free_space))
       |> Kernel.length()}
    end)
  end

  def format_part1(mapping) do
    do_format_part1(mapping, 0, length(mapping) - 1)
  end

  defp do_format_part1(mapping, left, right) do
    cond do
      left >= right ->
        mapping

      true ->
        case {Enum.at(mapping, left), Enum.at(mapping, right)} do
          {:free_space, {:file, _, _} = file} ->
            List.replace_at(mapping, left, file)
            |> List.replace_at(right, :free_space)
            |> do_format_part1(left + 1, right - 1)

          {{:file, _, _}, _} ->
            do_format_part1(mapping, left + 1, right)

          {_, :free_space} ->
            do_format_part1(mapping, left, right - 1)
        end
    end
  end

  def gen_disk_map(input) do
    do_gen_disk_map(input, [], 0, :file) |> Enum.reverse()
  end

  defp do_gen_disk_map(input, acc, file_count, mode) do
    case input do
      [] ->
        acc

      [x | tail] ->
        case mode do
          :file ->
            do_gen_disk_map(
              tail,
              List.duplicate({:file, file_count, x}, x) ++ acc,
              file_count + 1,
              :free_space
            )

          :free_space ->
            do_gen_disk_map(tail, List.duplicate(:free_space, x) ++ acc, file_count, :file)
        end
    end
  end

  defp visualize(disk_map) do
    Enum.map(disk_map, fn e ->
      case e do
        {:file, id, _} -> "#{id}"
        _ -> "."
      end
    end)
    |> Enum.join()
  end

  def test do
    disk_map = gen_disk_map([1, 2, 3, 4, 5])
    "0..111....22222" = visualize(disk_map)

    "00...111...2...333.44.5555.6666.777.888899" =
      gen_disk_map([2, 3, 3, 3, 1, 3, 3, 1, 2, 1, 4, 1, 4, 1, 3, 1, 4, 0, 2])
      |> visualize()

    "022111222......" =
      disk_map
      |> format_part1()
      |> visualize()

    disk_map |> scan_free_spaces() |> IO.inspect(label: "Scan free spaces")
  end
end

if length(System.argv()) > 0 do
  Solution.test()
else
  Solution.run()
end
