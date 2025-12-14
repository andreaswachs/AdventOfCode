use clap::Parser;
use std::fs::File;
use std::io::{BufReader, Read};

#[derive(Parser)]
#[command(name = "dec02")]
struct Args {
    #[arg(value_name = "FILENAME")]
    filename: String,
}

#[derive(Debug)]
#[warn(dead_code)]
struct Range {
    low: usize,
    high: usize,
}

impl Range {
    fn from(input: &str) -> Range {
        let tokens = input.split("-").collect::<Vec<&str>>();

        let low = tokens[0]
            .parse::<usize>()
            .expect("low number should been a valid unsigned integer");
        let high = tokens[1]
            .parse::<usize>()
            .expect("high number should been a valid unsigned integer");

        Range {
            low,
            high,
        }
    }

    fn part1_sum_invalid_ids(self: &Self) -> usize {
        (self.low..(self.high + 1))
            .filter(|n| {
                let as_str = n.to_string();
                let (a, b) = as_str.split_at(as_str.len() / 2);
                a == b
            })
            .sum::<usize>()
    }

    fn part2_sum_invalid_ids(self: &Self) -> usize {
      (self.low..(self.high + 1))
        .filter(|n| *n >= 10)
        .map(|n| {
          let as_str = n.to_string();


          let mut result = 0;
          let max = as_str.len()/2;
          let mut i = 1;

          loop {
            let chunks =
              as_str.as_bytes()
                .chunks(i)
                .map(|chunk| std::str::from_utf8(chunk).unwrap()) // Unwrap hack
                .collect::<Vec<&str>>();

            if chunks.iter().all(|x| *x == chunks[0]) {
              result += n;
              break
            }
            i += 1;
            if i > max {
              break
            }
          }

          result
        })
        .sum::<usize>()
    }
}

fn main() {
    let args = Args::parse();

    let file =
        File::open(args.filename.clone()).expect("should be able to read file path from args");
    let mut reader = BufReader::new(file);

    let mut content = String::new();

    reader
        .read_to_string(&mut content)
        .expect("should be able to read file till EOF");

    let part1 = content
        .trim()
        .split(",")
        .map(Range::from)
        .map(|r| r.part1_sum_invalid_ids())
        .sum::<usize>();

    println!("{:?}", part1);

    let part2 = content
        .trim()
        .split(",")
        .map(Range::from)
        .map(|r| r.part2_sum_invalid_ids())
        .sum::<usize>();

    println!("{:?}", part2);



}

#[cfg(test)]
mod tests {

    use super::*;

    #[test]
    fn test_invalid_ids_11_to_22() {
        let r = Range::from("11-22");
        assert_eq!(r.part1_sum_invalid_ids(), 33);
    }

    #[test]
    fn test_invalid_ids_95_to_115() {
        let r = Range::from("95-115");
        assert_eq!(r.part1_sum_invalid_ids(), 99);
    }

    #[test]
    fn test_part_2_from_10_to_12() {
        let r = Range::from("10-12");
        assert_eq!(r.part2_sum_invalid_ids(), 11);
    }

    #[test]
    fn test_part_2_from_10_to_22() {
        let r = Range::from("10-22");
        assert_eq!(r.part2_sum_invalid_ids(), 33);
    }

    #[test]
    fn test_part_2_from_2121212118_to_2121212124() {
        let r = Range::from("2121212118-2121212124");
        assert_eq!(r.part2_sum_invalid_ids(), 2121212121);
    }


    #[test]
    fn test_part_2_from_824824821_to_824824827() {
        let r = Range::from("824824821-824824827");
        assert_eq!(r.part2_sum_invalid_ids(), 824824824);
    }

    #[test]
    fn test_part_2_from_565653_to_565659() {
        let r = Range::from("565653-565659");
        assert_eq!(r.part2_sum_invalid_ids(), 565656);
    }

    #[test]
    fn test_part_2_from_38593856_to_38593862() {
        let r = Range::from("38593856-38593862");
        assert_eq!(r.part2_sum_invalid_ids(), 38593859);
    }

    #[test]
    fn test_part_2_from_1698522_to_1698528() {
        let r = Range::from("1698522-1698528");
        assert_eq!(r.part2_sum_invalid_ids(), 0);
    }

    #[test]
    fn test_part_2_from_222220_to_222224() {
        let r = Range::from("222220-222224");
        assert_eq!(r.part2_sum_invalid_ids(), 222222);
    }

    #[test]
    fn test_part_2_from_11_to_22() {
        let r = Range::from("11-22");
        assert_eq!(r.part2_sum_invalid_ids(), 11 + 22);
    }

    #[test]
    fn test_part_2_from_95_to_115() {
        let r = Range::from("95-115");
        assert_eq!(r.part2_sum_invalid_ids(), 99 + 111);
    }

    #[test]
    fn test_part_2_from_998_to_1012() {
        let r = Range::from("998-1012");
        assert_eq!(r.part2_sum_invalid_ids(), 999 + 1010);
    }

    #[test]
    fn test_part_2_from_1188511880_to_1188511890() {
        let r = Range::from("1188511880-1188511890");
        assert_eq!(r.part2_sum_invalid_ids(), 1188511885);
    }

    #[test]
    fn test_part_2_from_1111110_to_1111112() {
        let r = Range::from("1111111-1111112");
        assert_eq!(r.part2_sum_invalid_ids(), 1111111);
    }

    // Test cases from problem description: sequences repeated at least twice

    #[test]
    fn test_part_2_1234_repeated_twice() {
        let r = Range::from("12341234-12341234");
        assert_eq!(r.part2_sum_invalid_ids(), 12341234);
    }

    #[test]
    fn test_part_2_123_repeated_three_times() {
        let r = Range::from("123123123-123123123");
        assert_eq!(r.part2_sum_invalid_ids(), 123123123);
    }

    #[test]
    fn test_part_2_12_repeated_five_times() {
        let r = Range::from("1212121212-1212121212");
        assert_eq!(r.part2_sum_invalid_ids(), 1212121212);
    }

    #[test]
    fn test_part_2_1_repeated_seven_times() {
        let r = Range::from("1111111-1111111");
        assert_eq!(r.part2_sum_invalid_ids(), 1111111);
    }

    // Additional edge cases for repeated patterns

    #[test]
    fn test_part_2_two_digit_repetitions() {
        // 11, 22, 33, etc. should all be invalid
        let r = Range::from("11-33");
        assert_eq!(r.part2_sum_invalid_ids(), 11 + 22 + 33);
    }

    #[test]
    fn test_part_2_three_digit_repetitions() {
        // 111, 222, etc. should be invalid
        let r = Range::from("111-222");
        assert_eq!(r.part2_sum_invalid_ids(), 111 + 222);
    }

    #[test]
    fn test_part_2_single_digits() {
        // Single digits like 1, 2, 3 could be considered as the digit repeated once
        // Testing what the current implementation does
        let r = Range::from("1-9");
        // Based on the algorithm, single digits won't match because max = 1/2 = 0
        assert_eq!(r.part2_sum_invalid_ids(), 0);
    }

    #[test]
    fn test_part_2_non_repeating_patterns() {
        // Numbers that are NOT repeating patterns
        let r = Range::from("12345-12345");
        assert_eq!(r.part2_sum_invalid_ids(), 0);
    }

    #[test]
    fn test_part_2_partial_repetitions_dont_count() {
        // 12312 is NOT a valid repetition (123 doesn't repeat fully)
        let r = Range::from("12312-12312");
        assert_eq!(r.part2_sum_invalid_ids(), 0);
    }

    #[test]
    fn test_part_2_mixed_range_with_repetitions() {
        // Range containing both valid and invalid IDs
        let r = Range::from("100-111");
        // Only 111 should be invalid
        assert_eq!(r.part2_sum_invalid_ids(), 111);
    }
}
