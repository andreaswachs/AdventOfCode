if (args.Length != 1)
{
    Console.WriteLine("Usage: dotnet run <path-to-input>");
    return;
}

var input = String.Empty;

using(StreamReader sr = File.OpenText(args[0]))
{
    input = await sr.ReadToEndAsync();
}

Func<int, int> pointsCalc = (int winningCards) =>
{
    if (winningCards == 1)
    {
        return 1;
    }

    return (int)Math.Pow(2, winningCards - 1);
};

Action<string> part1 = (string input) =>
{
    var res = input
        .Split(Environment.NewLine)
        .ToList()
        .Where(s => !String.IsNullOrWhiteSpace(s))
        .Select(line =>
            {
                var split = line.Substring(line.IndexOf(":")+1)
                    .Split("|");

                var winningNumbers = split[0]
                    .Split(" ")
                    .Where(s => !String.IsNullOrWhiteSpace(s))
                    .Select(int.Parse)
                    .ToList();

                var drawnNumbers = split[1]
                    .Split(" ")
                    .Where(s => !String.IsNullOrWhiteSpace(s))
                    .Select(int.Parse)
                    .ToList();

                return winningNumbers.Intersect(drawnNumbers).Count();
            })
        .Where(i => i > 0)
        .Select(pointsCalc)
        .Sum();


    Console.WriteLine(res);
};

part1(input);

