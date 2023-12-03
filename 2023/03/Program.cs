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

        return positions.Where(_withinBounds(maxX, maxY));
    }

    private Func<(int, int), bool> _withinBounds(int maxx, int maxy)
    {
        Func<(int, int), bool> pred = pos => {
            return pos.Item1 >= 0 
                && pos.Item1 < maxx 
                && pos.Item2 >= 0 
                && pos.Item2 <= maxy;
        };

        return pred;
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

    private static void Main(string[] args)
    {
        var input = _readInput(args);

        Part1(input);
    }

}
