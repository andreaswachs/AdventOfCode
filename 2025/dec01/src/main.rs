use clap::Parser;

use std::fs::File;
use std::io::{BufRead, BufReader};

#[derive(Parser)]
#[command(name = "dec01")]
#[command(about = "Solve December 1 2025 Advent of Code problem")]
struct Args {
    #[arg(value_name = "FILENAME")]
    filename: String,
}

fn main() {
    let args = Args::parse();

    // Medium 🧠 level: we use some basic math that could probably be done
    // smarter if I were smarter
    let file = File::open(args.filename.clone()).expect("should have been able to opepn");
    let reader = BufReader::new(file);
    let mut part1_dial: isize = 50;
    let mut part1 = 0;
    for line in reader.lines() {
        let line = line.unwrap(); // fuck it, bail if we have a problem
        let direction = line.chars().take(1).collect::<String>();
        assert!(direction == "L" || direction == "R");

        let mut number = line
            .chars()
            .skip(1)
            .collect::<String>()
            .parse::<isize>()
            .unwrap();

        if direction == "L" {
            number = number * -1;
        }

        let mut buffer_num = part1_dial + number;
        if buffer_num < 0 {
            buffer_num += 100;
        }

        part1_dial = buffer_num % 100;
        if part1_dial == 0 {
            part1 += 1;
        }
    }

    println!("Part 1: {}", part1);

    // part2: Zero 🧠 idea: we open and read the file again
    // instead of doing math that I might get wrong we just iterate one
    // step at the time to determine the number that the dial is pointing at
    let file = File::open(args.filename.clone()).expect("should have been able to opepn");
    let reader = BufReader::new(file);
    let mut part2_dial: isize = 50;
    let mut part2 = 0;
    for line in reader.lines() {
        let line = line.unwrap(); // fuck it, bail if we have a problem
        let direction = line.chars().take(1).collect::<String>();
        assert!(direction == "L" || direction == "R");

        let mut number = line
            .chars()
            .skip(1)
            .collect::<String>()
            .parse::<isize>()
            .unwrap();

        let delta = match direction.as_str() {
            "L" => -1,
            "R" => 1,
            &_ => todo!(), // We will never get here
        };

        while number > 0 {
            part2_dial += delta;

            if part2_dial == -1 {
                part2_dial = 99;
            }

            if part2_dial == 100 {
                part2_dial = 0;
            }

            if part2_dial == 0 {
                part2 += 1;
            }

            number -= 1;
        }
    }

    println!("Part 2: {}", part2);
}
