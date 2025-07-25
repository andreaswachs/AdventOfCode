defmodule Solution do
  def run do
    if length(System.argv()) > 0 do
      test()
      exit(0)
    end

    input =
      IO.read(:stdio, :eof)
      |> String.split("\n")

    _part1 =
      input
      |> Enum.map(&String.split(&1, ": "))
      |> Enum.map(fn l ->
        result = Enum.at(l, 0) |> Integer.parse() |> elem(0)

        numbers =
          Enum.at(l, 1) |> String.split(" ") |> Enum.map(fn n -> Integer.parse(n) |> elem(0) end)

        results = interlace_part1(numbers)

        {result, result in results}
      end)
      |> Enum.filter(fn {_, ok} -> ok end)
      |> Enum.map(fn {result, _} -> result end)
      |> Enum.sum()
      |> IO.inspect(label: "Part 1")

    # This contains a list of tuples {result, [n_1, n_2, n_3, ... n_n]}
    _part2 =
      input
      |> Enum.map(&String.split(&1, ": "))
      |> Enum.map(fn l ->
        result = Enum.at(l, 0) |> Integer.parse() |> elem(0)

        numbers =
          Enum.at(l, 1) |> String.split(" ") |> Enum.map(fn n -> Integer.parse(n) |> elem(0) end)

        results = interlace(numbers)

        {result, result in results}
      end)
      |> Enum.filter(fn {_, ok} -> ok end)
      |> Enum.map(fn {result, _} -> result end)
      |> Enum.sum()
      |> IO.inspect(label: "Part 2")
  end

  @doc """
    Will interlace a list of numbers with operators
  """
  defp interlace(numbers) do
    interlace_do(numbers, 0)
  end

  defp interlace_do(numbers, acc) do
    case numbers do
      [] ->
        [acc]

      [x | tail] ->
        interlace_do(tail, concat(acc, x)) ++
          interlace_do(tail, x * if(acc == 0, do: 1, else: acc)) ++
          interlace_do(tail, x + acc)
    end
  end

  def interlace_part1(numbers, acc \\ 0) do
    case numbers do
      [] ->
        [acc]

      [x | tail] ->
        interlace_part1(tail, x * if(acc == 0, do: 1, else: acc)) ++
          interlace_part1(tail, x + acc)
    end
  end

  defp concat(left, right) do
    case {left, right} do
      {0, 0} -> 0
      {0, _} -> right
      {_, 0} -> left
      {_, _} -> "#{left}#{right}" |> Integer.parse() |> elem(0)
    end
  end

  defp test do
    interlace([6, 8, 6, 15])
    |> IO.inspect()
  end
end

Solution.run()
