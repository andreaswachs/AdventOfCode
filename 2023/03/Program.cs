internal class Utils
{
    public static Func<(int, int), bool> WithinBounds(int maxx, int maxy)
        => new Func<(int, int), bool>((pos) => {
            return pos.Item1 >= 0 
                && pos.Item1 < maxx 
                && pos.Item2 >= 0 
                && pos.Item2 <= maxy;
        });
}

internal class PartsNumber
{
    public int Value { get; set; }
    public int X { get; set; }
    public int Y { get; set; }

    public PartsNumber(int value, int x, int y)
    {
        this.Value = value;
        this.X = x;
        this.Y = y;
    }

    public IEnumerable<(int, int)> AdjacentPositions(int maxX, int maxY)
    {
        IList<(int, int)> positions = new List<(int, int)>();

        int length = this.Value.ToString().Length;

        // We want to iterate within the vertical bars, including the area
        // the vertical bars point at:
        // ...12345...
        //   |-----|
        for(int i = this.X - 1; i <= this.X + length; i++)
        {
            // Adding positions
            // .###.
            // .?A?. (? might be numbers)
            // .###.
            positions.Add((i, this.Y - 1));
            positions.Add((i, this.Y + 1));
        }

        // Add the following # positions
        // ..........
        // ..#123#...
        // ..........
        positions.Add((this.X - 1, this.Y));
        positions.Add((this.X + length, this.Y));

        return positions.Where(Utils.WithinBounds(maxX, maxY));
    }


}

internal class ProgramInput
{
    public int GridX { get; }
    public int GridY { get; }
    public IList<PartsNumber> PartsNumbers { get; }
    public IList<string> Grid { get; }

    public ProgramInput(int gridx, int gridy, IList<PartsNumber> numbers, IList<string> grid)
    {
        this.GridX = gridx;
        this.GridY = gridy;
        this.PartsNumbers = numbers;
        this.Grid = grid;
    }
}

internal class Program
{
    private static ProgramInput _readInput(string[] args)
    {
        int gridx = 0;
        int gridy = 0;

        IList<PartsNumber> partsNumbers = new List<PartsNumber>();
        IList<string> grid = new List<string>();                          

        using(StreamReader sr = File.OpenText(args[0]))
        {
            int y = 0;
            string? line = String.Empty;

            while((line = sr.ReadLine()) != null)
            {
                if (line.Length == 0) continue;
                gridx = line.Length;
                gridy = y;

                grid.Add(line);

                IList<string> buf = new List<string>();
                System.Text.StringBuilder sb = new System.Text.StringBuilder();

                // We find numbers on a line old-school style: scan left to right
                for(int i = 0; i < line.Length; i++)
                {
                    if (Char.IsDigit(line[i]))
                    {
                        sb.Append(line[i]);
                    } else {
                        if (sb.Length > 0) {
                            buf.Add(sb.ToString());
                            sb.Clear();
                        }
                    }
                }
                // Check if there were any hanging number in the string builder
                // (this is caused from a number ending at the very rightmost position
                // in the grid)
                if (sb.Length > 0) {
                    buf.Add(sb.ToString());
                }

                // Offset is used to prevent errors with duplicate numbers on the same line
                // and determining the same position for the two numbers
                int offset = 0;
                foreach(string pn in buf)
                {
                    int x = line.IndexOf(pn, offset);
                    int value = int.Parse(pn);
                    partsNumbers.Add(new PartsNumber(value, x, y));
                    offset = x + pn.Length;
                }

                y++;
            }
        }

        return new ProgramInput(gridx, gridy, partsNumbers, grid);
    }

    private static void Part1(ProgramInput input)
    {
        var part1 = 0;

        foreach(PartsNumber pn in input.PartsNumbers)
        {
            foreach((int, int) pos in pn.AdjacentPositions(input.GridX, input.GridY))
            {
                char buf = input.Grid[pos.Item2][pos.Item1];
                if (buf != '.' && !Char.IsDigit(buf))
                {
                    part1 += pn.Value;
                    break;
                }
            }
        }

        Console.WriteLine($"Part 1: {part1}");
    }

    private static void Part2(ProgramInput input)
    {
        int y = 0;

        long sum = 0;

        foreach(string line in input.Grid)
        {
            for(int i = 0; i < line.Length; i++)
            {
                // Check if there is a possible gear connection
                if (line[i] == '*')
                {
                    (string num1, string num2) = adjacentPartsNumbers(input, (i, y));
                    
                    // We might reach cases where there is not two adjacent part numbers
                    // then its not a gear
                    if (num1 == String.Empty || num2 == String.Empty)
                    {
                        continue;
                    }

                    sum += long.Parse(num1) * long.Parse(num2);
                }
            }

            y++;
        }

        Console.WriteLine($"Part 2: {sum}");
    }

    private static (string, string) adjacentPartsNumbers(ProgramInput input, (int, int) source)
    {
        (int x, int y) = source;
        IList<(int, int)> adj = new List<(int, int)>();
        var withinBounds = Utils.WithinBounds(input.GridX, input.GridY);

        // Find all adjacent numbers around source
        for(int i = x - 1; i <= x + 1; i++)
        {
            for(int j = y - 1; j <= y + 1; j++) {
                if (withinBounds((i, j)) && Char.IsDigit(input.Grid[j][i])) {
                    adj.Add((i, j));
                }
            }
        }

        IList<((int, int), string)> numbersInScope = new List<((int, int), string)>();
        foreach((int, int) a in adj)
        {
            numbersInScope.Add(seekPartsNumber(input, a));
        }

        // We risk having duplicate numbers, so we ask for the distrinct
        // positions and numbers
        var res = numbersInScope.Distinct().ToList();
        if (res.Count() != 2)
        {
            return (string.Empty, string.Empty);
        }

        return (res[0].Item2, res[1].Item2);
    }

    /// seekPartsNumber returns the beginnign position of a number and the
    /// number in string representation
    private static ((int, int), string) seekPartsNumber(ProgramInput input, (int, int) source)
    {
        System.Text.StringBuilder sb = new System.Text.StringBuilder();
        (int x, int y) = source;

        while(x > 0 && Char.IsDigit(input.Grid[y][x]))
        {
            // We move the cursor to the left
            x--;
        }
        // If we accidentally moved the curser beyond the number, we move it back one spot
        if (!Char.IsDigit(input.Grid[y][x])) 
        {
            x++;
        }

        (int, int) pos = (x, y);
        while(x < input.GridX && Char.IsDigit(input.Grid[y][x]))
        {
            // We eat up all the digits until we reach an end to tne number or 
            // we run out of space on the grid
            sb.Append(input.Grid[y][x++]);
        }

        return (pos, sb.ToString());
    }

    private static void Main(string[] args)
    {
        if (args.Length != 1) {
            Console.WriteLine("Usage: dotnet run <input_file>");
            System.Environment.Exit(1);
        }

        var input = _readInput(args);
        Part1(input);
        Part2(input);
    }
}
