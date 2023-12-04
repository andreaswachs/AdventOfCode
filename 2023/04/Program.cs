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

Func<string, int[]> winningCards = (string input) =>
{
    return input
        .Split(Environment.NewLine)
        .ToList()
        .Where(s => !String.IsNullOrWhiteSpace(s))
        .Select(line =>
            {
                var split = line.Substring(line.IndexOf(":")+1)
                    .Split("|");

                // We make a nice list of the winning numbers
                var winningNumbers = split[0]
                    .Split(" ")
                    .Where(s => !String.IsNullOrWhiteSpace(s))
                    .Select(int.Parse)
                    .ToList();

                // Same as before, for the drawn numbers
                var drawnNumbers = split[1]
                    .Split(" ")
                    .Where(s => !String.IsNullOrWhiteSpace(s))
                    .Select(int.Parse)
                    .ToList();

                // Find the size of the intersection of winning and drawn numbers
                // This shows how many drawn cards were winning cards
                return winningNumbers.Intersect(drawnNumbers).Count();
            })
        .ToArray();
};

Action<string> part1 = (string input) =>
{
    var res = winningCards(input)
        .Where(i => i > 0)
        .Select(pointsCalc)
        .Sum();

    Console.WriteLine($"Part 1: {res}");
};

Action<string> part2 = (string input) =>
{

    var cardWins = winningCards(input);
    // We'll use a dictionary to keep track of how many copies we win of each
    // subsequent card, as we might have more cards than the input has,
    // so a fixed sized array doesnt work. Could make work with a list, but
    // this is easier to implement fast
    var cardCount = new Dictionary<int, int>();

    // Initialize the card count
    for(int i = 0; i < cardWins.Count(); i++)
    {
        cardCount.Add(1, 1);
    }

    // We go through each of the cards
    for(int i = 1; i < cardWins.Count(); i++)
    {
        // For as many times as we have copies of the given card
        for(int k = 1; k <= cardCount[i]; k++)
        {
            // Add copies of the subsequent cards, if we have any wins
            var wins = cardWins[i-1];
            for(int j = 1; j <= wins; j++)
            {
                int buf = -1;
                if (cardCount.TryGetValue(i+j, out buf))
                {
                    cardCount[i+j] = buf + 1;
                } else {
                    cardCount.Add(i+j, 1);
                }
            }
        }
    }

    Console.WriteLine($"Part 2: {cardCount.Values.Sum()}");
};

part1(input);
part2(input);
