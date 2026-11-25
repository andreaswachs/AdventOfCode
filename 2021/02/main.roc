app [main!] {
	cli: platform "https://github.com/roc-lang/basic-cli/releases/download/0.28.1/7brgwR8YUqfeGrkdfF1V4GPPN7DnP7ftqub9KcrTEZtr.tar.zst",
	roc: "nightly-2026-11-14-a87808d",
}

import cli.IOErr
import cli.Path
import cli.Stdout

main! = |args| {
	path = Path.from_os_str(args.first()?)
	moves = input!(path)?
	solve_part1!(moves)?
	solve_part2!(moves)?

	Ok({})
}

solve_part1! : List((Str, I64)) => Try({}, [StdoutErr(IOErr), Interrupted])
solve_part1! = |moves| {
	result_pair = moves.fold(
		(0.I64, 0.I64),
		|(x, y), (direction, dist)| {
			match direction {
				"forward" => (x + dist, y)
				"down" => (x, y + dist)
				"up" => (x, y - dist)
				bad => crash "INVALID DIRECTION: ${Str.inspect(bad)}"
			}
		},
	)

	Stdout.line!("[DEBUG]: Answer part 1: ${Str.inspect(result_pair.0 * result_pair.1)}")?
	Ok({})
}

solve_part2! : List((Str, I64)) => Try({}, [StdoutErr(IOErr), Interrupted])
solve_part2! = |moves| {
	result_pair = moves.fold(
		(0.I64, 0.I64, 0.I64),
		|(x, y, aim), (direction, dist)| {
			match direction {
				"forward" => (x + dist, y + (aim * dist), aim)
				"down" => (x, y, aim + dist)
				"up" => (x, y, aim - dist)
				bad => crash "INVALID DIRECTION: ${Str.inspect(bad)}"
			}
		},
	)

	Stdout.line!("[DEBUG]: Answer part 2: ${Str.inspect(result_pair.0 * result_pair.1)}")?
	Ok({})
}

input! : Path => Try(List((Str, I64)), [MalformedMove(Str), PathErr(IOErr, Path.Path)])
input! = |path| {
	text = Path.read_utf8!(path)?
	text.trim().split_on("\n").map_try(parse_move)
}

parse_move : Str -> Try((Str, I64), [MalformedMove(Str)])
parse_move = |line| {
	parts = line.split_on(" ")
	direction = parts.get(0).map_err(|_| MalformedMove(line))
	distance = parts.get(1).map_err(|_| MalformedMove(line))

	match (direction, distance) {
		(Ok(dir), Ok(dist)) =>
			I64.from_str(dist)
				.map_err(|_| MalformedMove(line))
				.map_ok(|amount| (dir, amount))
		_ => Err(MalformedMove(line))
	}
}
