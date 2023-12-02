const int RED = 0;
const int GREEN = 1;
const int BLUE = 2;
int[] sums = new int[3];
int [] maxAmounts = new int[3] {12, 13, 14};

Action<int[]> resetArray = a => Array.Fill(a, 0);

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

Func<string, bool> gamePossible = s => {
    return s
        .Substring(s.IndexOf(":") + 1)
        .Split(",;".ToCharArray())
        .Select(s => s.Trim())
        .Where(s => s.Length > 4)
        .All(drawPossible);
};

Func<string, int> gameId = s =>
     int.Parse(s.Split(":")[0].Split(" ")[1]);

if (args.Length != 1) {
    Console.WriteLine("Usage: dotnet run <path to input file>");
    return;
}

using(StreamReader sr = File.OpenText(args[0])) {   
    string? input = String.Empty;

    var res = sr.ReadToEnd()
        .Split("\n")
        .Where(s => s.Length > 0)
        .Where(gamePossible)
        .Select(gameId)
        .Sum();

    Console.WriteLine($"Part 1: {res}");
}
