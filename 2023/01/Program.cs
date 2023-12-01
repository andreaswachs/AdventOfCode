internal class Program
{
    private static string[] digits = 
    {
        "one", "two", "three",
        "four", "five", "six",
        "seven", "eight", "nine"
    };

    private static IEnumerable<string> GetForwardEnumerableOf(string s) {
        // We find a maximum window size of either 5 or the length of the string
        int right = s.Length > 5 ? 5 : s.Length;
        for(int left = 0; left < s.Length; left++) {

            // Yield the substring  
            yield return s.Substring(left, right - left);

            // If the window has not reached the right side, we move the window
            // to the right
            if (right < s.Length) {
                right++;
            }
        }
    }

    private static IEnumerable<string> GetBackwardEnumerableOf(string s) {
        // Move a window through the string from the right to the left
        // with variable length from 1 to 5
        for(int left = s.Length - 1, right = s.Length; right >= 1;) {

            // Yield the substring
            yield return s.Substring(left, right - left);

            // If the window has reached the left side, we tighten the window from the right
            if (left == 0) {
                right--;
            } 

            // If we have not reached the left side, we move the window
            if (left > 0) {
                left--;
            }

            // If the window is too large, we tighten it from the right
            if (right - left > 5) {
                right--;
            }
        }
    }

    private static int GetFirstDigitInDirection(IEnumerable<string> directionEnumerable)
    {
        foreach(string window in directionEnumerable)
        {
            // First we check if we have reached a "simple" digit, i.e. 0-9
            if (Char.IsDigit(window[0]))
            {
                return (int) window[0] - '0';
            }

            // We now check if we have reached a word digit, i.e. one, two, three, ...
            foreach(var digit in digits)
            {
                if (window.StartsWith(digit))
                {
                    return Array.IndexOf(digits, window) + 1;
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

                var a = GetFirstDigitInDirection(forward) * 10;
                var b = GetFirstDigitInDirection(backward);

                sum += a + b;
            }

            Console.WriteLine($"{sum}");
        }
    }
}
