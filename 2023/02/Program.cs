const int RED = 0;
const int GREEN = 1;
const int BLUE = 2;
int [] maxAmounts = new int[3] {12, 13, 14};

Func<string, int> toColourIdx = s => {
    switch (s) {
        case "blue": 
            return BLUE;
        case "red": 
            return RED;
        case "green": 
            return GREEN;
        default: 
            throw new InvalidDataException($"{s} is not a valid colour idx");
    };
};

Func<string, bool> drawPossible = s => {
    // Evaluates if "<amount> <colour>" is valid
    var draw = s.Split(" ");
    var idx = toColourIdx(draw[1]);
    var amount = int.Parse(draw[0]);

    return amount <= maxAmounts[idx];
};

Func<string, bool> gamePossible = s => 
    s
    .Substring(s.IndexOf(":") + 1)
    .Split(",;".ToCharArray())
    .Select(s => s.Trim())
    .Where(s => s.Length > 4)
    .All(drawPossible);


Func<string, (int, int)> toColourAmount = s => {
    var draw = s.Split(" ");
    var idx = toColourIdx(draw[1]);
    var amount = int.Parse(draw[0]);
    return (amount, idx);
};

Func<string, int[]> minimumPossibleCubeSet = s =>
    s.Substring(s.IndexOf(":") + 1)
    .Split(",;".ToCharArray())
    .Where(s => s.Length > 4)
    .Select(s => s.Trim())
    .Select(toColourAmount)
    .Aggregate(new int[3] {0, 0, 0}, (acc, x) => {
        // Count the minimum number of cubes of a given colour
        // for this to be a possible play
        (int amount, int idx) = x;
        if (acc[idx] < amount) { 
            acc[idx] = amount; 
        }
        return acc;
    });

Func<string, int> gameId = s =>
     int.Parse(s.Split(":")[0].Split(" ")[1]);

if (args.Length != 1) {
    Console.WriteLine("Usage: dotnet run <path to input file>");
    return;
}

using(StreamReader sr = File.OpenText(args[0])) {   
    string? input = String.Empty;

    var data = sr
        .ReadToEnd()
        .Split('\n')
        .Where(s => s.Length > 0);

    var part1 = data
        .Where(gamePossible)
        .Select(gameId)
        .Sum();

    Console.WriteLine($"Part 1: {part1}");

    var part2 = data
        .Select(minimumPossibleCubeSet)
        .Select(s => s.Aggregate(1, (acc, x) => acc * x))
        .Sum();

    Console.WriteLine($"Part 2: {part2}");
}
