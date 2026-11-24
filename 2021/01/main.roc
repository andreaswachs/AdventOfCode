app [main!] {
	cli: platform "https://github.com/roc-lang/basic-cli/releases/download/0.28.1/7brgwR8YUqfeGrkdfF1V4GPPN7DnP7ftqub9KcrTEZtr.tar.zst",
	roc: "nightly-2026-11-14-a87808d",
}

import cli.Cmd
import cli.Stdout

main! = |args| {
	input = read!(args.first()?)?
		.trim()
		.split_on("\n")
		.keep_oks(I64.from_str)

	count_increasing : List(I64) -> (I64, I64)
	count_increasing = |arg| {
		arg
			.fold(
				(0, I64.highest),
				|(counter, prev), next| {
					new_counter = if (next > prev) counter + 1 else counter

					(new_counter, next)
				},
			)
	}

	ans1 = count_increasing(input)
	Stdout.line!("[debug] ans1 = ${Str.inspect(ans1.0)}")?

	ans2_input = input.map_with_index(|_, i| input.sublist({ start: i, len: 3 }))
		.drop_last(2)
		.map(|l| l.fold(0, I64.plus))

	ans2 = count_increasing(ans2_input)
	Stdout.line!("[debug] ans2 = ${Str.inspect(ans2.0)}")?

	Ok({})
}

read! = |filename| {
	Cmd.new("cat")
		.arg(filename)
		.capture!()
}
