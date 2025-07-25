defmodule Solution do
  def run do
    prep = 
      IO.read(:stdio, :eof)
      |> String.split("\n")
      |> Enum.map(&(String.split(&1, " ")) |> Enum.map(fn e -> Integer.parse(e) |> elem(0) end))

    prep
    |> part1()
    |> IO.inspect(label: "Part 1")

    prep
    |> part2()
    |> IO.inspect(label: "Part 2")
  end

  def part1(prep) do
    prep 
    |> Enum.filter(&part1_verify/1)
    |> Enum.count()
  end

  def part2(prep) do
    prep
    |> Enum.map(&(part2_verify(&1)))
    |> part2_filter()
    |> Enum.count()
  end

  def part1_verify([a | [b | [c | _]]] = lst) do
    cond do
      a < b and b < c and diff_ok?(a, b) and diff_ok?(b, c) -> Enum.drop(lst, 1) |> part1_verify()
      a > b and b > c and diff_ok?(a, b) and diff_ok?(b, c) -> Enum.drop(lst, 1) |> part1_verify()
      true -> false
    end
  end

  def part1_verify([a | [b | _]] = lst) do
    cond do
      a < b and diff_ok?(a, b) -> Enum.drop(lst, 1) |> part1_verify()
      a > b and diff_ok?(a, b) -> Enum.drop(lst, 1) |> part1_verify()
      true -> false
    end
  end

  def part1_verify([_]) do
    true
  end

  def diff_ok?(a, b) do
    abs(a - b) in [1, 2, 3]
  end

  def part2_verify(lst) do
    scanned =
      Enum.scan(lst, {:begin, Enum.at(lst, 0), :ok}, fn x, prev -> 
        case prev do
          {:begin, y, _} when x == y -> {:begin, x, :ok}
          {_, y, _} when x > y -> 
            if diff_ok?(y, x) do 
              {:up, x, :ok}
            else 
              {:up, x, :bad}
            end
          {_, y, _} when x < y -> 
            if diff_ok?(y, x) do
              {:down, x, :ok} 
            else
              {:down, x, :bad}
            end
          {_, y, _} when x == y -> {:eq, x, :bad}
        end
      end)
    scanned
  end

  def faults(lst) do
    case lst do
      [] -> 0
      [{:up, _, _} , {:down, _, _} | _] -> 1 + faults(Enum.drop(lst, 2))
      [{:down, _, _} , {:up, _, _} | _] -> 1 + faults(Enum.drop(lst, 2))
      [{_, _, _} , {:eq, _, _} | _] -> 1 + faults(Enum.drop(lst, 2))
      _ -> 0 + faults(Enum.drop(lst, 1))
    end
  end

  def part2_filter(lsts) do 
    lsts
    |> Enum.filter(&faults(&1) < 2)
    |> Enum.filter(fn lst -> 
      Enum.all?(lst, fn e -> 
        case e do
          {:eq, _, _} -> true # edge case
          {_, _, :bad} -> false
          _ -> true
        end
      end)
    end)
  end
end

Solution.run()
