defmodule Solution do
  def run do
    input = 
      IO.read(:stdio, :eof)
      |> String.trim()

    part1 =
      mul_scan(input)
      |> List.flatten()
      |> Enum.map(&mul_transform/1)
      |> Enum.sum()
      |> IO.inspect(label: "Part 1")


    _ = part1
  end

  defp mul_scan(x) do
    Regex.scan(~r/mul\(\d\d?\d?,\d\d?\d?\)/, x)
  end

  defp mul_scan_cond(x) do
  end

  defp do_mul_scan_cond(x, conditional_state, muls) do
    case x do
      "mul" <> rest -> # TODO
      "don't()" <> rest -> do_mul_scan_cond(List.dr)
    end
  end

  defp mul_transform(stmt) do
    Regex.named_captures(~r/(?<left>\d*),(?<right>\d*)/, stmt)
    |> Map.values()
    |> Enum.map(&to_int/1)
    |> Enum.product()
  end

  defp to_int(x), do: Integer.parse(x) |> elem(0)
end
Solution.run()
