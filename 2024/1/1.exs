defmodule Solution do
  def run do
    prep = 
      IO.read(:stdio, :eof)
      |> String.trim_trailing()
      |> String.split("\n")
      |> Enum.map(&(String.split(&1, "   ")))
      |> unzip()
      |> sort()

    prep
    |> part1()
    |> IO.inspect(label: "Part 1")

    prep
    |> part2()
    |> IO.inspect(label: "Part 2")
  end

  def part1(lst) do
    Enum.zip(elem(lst, 0), elem(lst, 1))
    |> Enum.reduce(0, fn x,acc ->
      {a, b} = x
      abs(a - b) + acc
    end)
  end

  def part2({l1, l2}) do
    f = Enum.frequencies(l2)
    Enum.reduce(l1, 0, fn x, acc ->
      Map.get(f, x, 0) * x + acc
    end)
  end

  def unzip(lst) do
    Enum.reduce(lst, {[], []}, fn e, acc ->
      {l1, l2} = acc
      [a | [b | _]] = e
      {aa, _} = Integer.parse(a)
      {bb, _} = Integer.parse(b)
      {[aa | l1], [bb | l2]}
    end)
  end

  def sort({l1, l2}) do
    {Enum.sort(l1), Enum.sort(l2)}
  end
end

Solution.run()
