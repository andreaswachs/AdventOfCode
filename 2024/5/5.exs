defmodule Solution do
  def run do
    if length(System.argv()) > 0 do
      tests()
      exit(0)
    end

    input =
      IO.read(:stdio, :eof)
      |> String.split("\n\n")

    rules = input |> Enum.at(0)
    updates = input |> Enum.at(1)

    rules_map = build_rule_map(rules)

    update_sequences_ints_with_idx =
      updates
      |> String.split("\n")
      |> Enum.with_index()
      |> Enum.map(fn {line, i} ->
        # converting to integer list
        {line |> String.split(",") |> Enum.map(&(Integer.parse(&1) |> elem(0))), i}
      end)

    valid_updates =
      update_sequences_ints_with_idx
      # checking that all the updates are in the right order, keep the indices
      |> Enum.map(&{check_ok(rules_map, elem(&1, 0)) |> List.flatten(), elem(&1, 1)})
      # Filtering out only sequences of updates that are valid
      |> Enum.filter(&Enum.all?(elem(&1, 0)))
      # Getting the indices of the valid updates
      |> Enum.map(&elem(&1, 1))

    _part1 =
      update_sequences_ints_with_idx
      |> Enum.filter(&(elem(&1, 1) in valid_updates))
      |> Enum.map(&elem(&1, 0))
      |> sum_middle_values()
      |> IO.inspect(label: "Part 1")

    _all_wrong_updates =
      update_sequences_ints_with_idx
      |> Enum.reject(&(elem(&1, 1) in valid_updates))
      |> Enum.map(&elem(&1, 0))
      |> Enum.map(&reconstruct(rules_map, &1))
      |> sum_middle_values()
      |> IO.inspect(label: "Part 2")
  end

  @doc """
    Sums the values of the middle element in the given sequences
  """
  def sum_middle_values(sequences) do
    sequences
    |> Enum.map(&(Enum.drop(&1, div(length(&1), 2)) |> Enum.take(1)))
    |> List.flatten()
    |> Enum.sum()
  end

  @doc """
    Builds the rule map, which is a mappping of a given number
    that poinst to a list of numbers which must come after it
  """
  def build_rule_map(rules) do
    rules
    |> String.split("\n")
    |> Enum.reduce(%{}, fn rule, acc ->
      [key, value] = String.split(rule, "|")
      {k, _} = Integer.parse(key)
      {v, _} = Integer.parse(value)
      x = Map.get(acc, k, [])
      Map.put(acc, k, [v | x])
    end)
  end

  @doc """
    Given a list of numbers and the rules map, this function
    reconstructs the only valid sequence by brute force
  """
  defp reconstruct(rules_map, numbers) do
    # sequence
    reconstruct_do(rules_map, Enum.drop(numbers, 1), [Enum.take(numbers, 1)])
  end

  @doc """
    The auxillary function to reconstruct/1
  """
  defp reconstruct_do(rules_map, numbers, sequences) do
    case numbers do
      [] ->
        # we should expect there only to be one result in the end
        sequences |> hd

      [x | tail] ->
        new_seqs = insert(rules_map, x, sequences)
        reconstruct_do(rules_map, tail, new_seqs)
    end
  end

  @doc """
    Inserts the given number into the given sequences and returns only valid sequences

    OBS: might return multiple sequences as there might intermittently be
    multiple valid sequences in processing
  """
  defp insert(rules_map, number, sequences) do
    new_sequences =
      Enum.map(sequences, fn sequence ->
        for i <- 0..length(sequence) do
          List.insert_at(sequence, i, number)
        end
      end)

    new_sequences
    |> Enum.map(
      &Enum.map(&1, fn s -> {s, check_ok(rules_map, s) |> List.flatten() |> Enum.all?()} end)
    )
    |> List.flatten()
    |> Enum.filter(fn {_, ok} -> ok end)
    |> Enum.map(&elem(&1, 0))
  end

  @doc """
    Given a rules_map and a sequence, checks if the updates are in the right order
  """
  def check_ok(rules_map, sequence) do
    for i <- 1..(length(sequence) - 1) do
      a = Enum.at(sequence, i)

      for j <- 0..(i - 1) do
        b = Enum.at(sequence, j)

        if not check_following(rules_map, b, a) do
          false
        else
          true
        end
      end
    end
  end

  # Checking that it's ok that b comes after a in the number sequence
  def check_following(rules_map, a, b) do
    a not in Map.get(rules_map, b, [])
  end

  @doc """
    Sourced from the test input, ran through the build_rule_map/1 function
  """
  def test_rules_map do
    %{
      29 => [13],
      47 => [29, 61, 13, 53],
      53 => [13, 29],
      61 => [29, 53, 13],
      75 => [13, 61, 47, 53, 29],
      97 => [75, 53, 29, 47, 61, 13]
    }
  end

  @doc """
    Tests to validate correctness in single case(s)
  """
  def tests do
    [97, 75, 47, 29, 13] = reconstruct(test_rules_map(), [97, 13, 75, 29, 47])
  end
end

Solution.run()
