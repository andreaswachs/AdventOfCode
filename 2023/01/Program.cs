internal class Program
{
    private static string[] digits = 
    {
        "one", "two", "three",
        "four", "five", "six",
        "seven", "eight", "nine"
    };

    private static IEnumerable<string> GetForwardEnumerableOf(string s) {
        int right = s.Length > 5 ? 5 : s.Length;
        for(int left = 0; left < s.Length; left++) {
            yield return s.Substring(left, right - left);

            if (right < s.Length) {
                right++;
            }
        }
    }

    private static IEnumerable<string> GetBackwardEnumerableOf(string s) {
        for(int left = s.Length - 1, right = s.Length; right >= 1;) {
            yield return s.Substring(left, right - left);
            if (left == 0) {
                right--;
            } 
            if (left > 0) {
                left--;
            }
            if (right - left > 5) {
                right--;
            }
        }
    }
    private static int stringDigitToInt(string s)
        => Array.IndexOf(digits, s) + 1;

    private static int get(IEnumerable<string> e)
    {
        foreach(string window in e)
        {
            if (Char.IsDigit(window[0]))
            {
                return (int) window[0] - '0';
            }

            foreach(var digit in digits)
            {
                if (window.StartsWith(digit))
                {
                    return stringDigitToInt(digit);
                }
            }

        }

        throw new Exception("I should not have reached this place");
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
                var forward = GetForwardEnumerableOf(s);
                var backward = GetBackwardEnumerableOf(s);

                var a = get(forward) * 10;
                var b = get(backward);

                sum += a + b;
            }

            Console.WriteLine($"{sum}");
        }
    }
}
