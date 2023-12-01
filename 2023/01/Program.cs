internal class Program
{
    private static string[] digits = 
    {
        "one",
        "two",
        "three",
        "four",
        "five",
        "six",
        "seven",
        "eight",
        "nine"
    };

    private static int stringDigitToInt(string s)
    {
        switch (s)
        {
            case "one":
                return 1;
            case "two":
                return 2;
            case "three":
                return 3;
            case "four":
                return 4;
            case "five":
                return 5;
            case "six":
                return 6;
            case "seven":
                return 7;
            case "eight":
                return 8;
            case "nine":
                return 9;

            default:
                return -100000;
        }
    }

    private static int getA(string s)
    {
        int limit = s.Length;
        for(int i = 0; i < limit; i++)
        {
            if (Char.IsDigit(s[i]))
            {
                return (int) s[i] - '0';
            }

            var sub = s.Substring(i);
            foreach(var digit in digits)
            {
                if (sub.StartsWith(digit))
                {
                    return stringDigitToInt(digit);
                }
            }
        }

        throw new Exception($"I could not compute A for {s}");
    }
        
    private static int getB(string s)
    {
        int limit = s.Length - 1;
        for(int i = limit; i >= 0; i--)
        {
            if (Char.IsDigit(s[i]))
            {
                return (int) s[i] - '0';
            }

            var sub = s.Substring(i);
            foreach(var digit in digits)
            {
                if (sub.StartsWith(digit))
                {
                    return stringDigitToInt(digit);
                }
            }
        }

        throw new Exception($"I could not compute B for {s}");
    }
    private static void Main(string[] args)
    {
        if (args.Length != 1)
        {
            Console.WriteLine("only expected path to input file as argument");
            return;
        }

        using (StreamReader sr = File.OpenText(args[0]))
        {
            int sum = 0;

            string? s = string.Empty;
            while ((s = sr.ReadLine()) != null)
            {

                sum += getA(s) * 10 + getB(s);
            }

            Console.WriteLine($"{sum}");
        }
    }
}
